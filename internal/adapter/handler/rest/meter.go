package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alhamsya/voltron/internal/core/domain/constant"
	"github.com/alhamsya/voltron/pkg/manager/logging"
	"github.com/alhamsya/voltron/pkg/manager/response"
	"github.com/alhamsya/voltron/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"

	modelPower "github.com/alhamsya/voltron/internal/core/domain/power"
	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"
)

func (h *Handler) Reading(ctx *fiber.Ctx) error {
	req := new(modelRequest.ReqHandlerMeterReading)
	err := ctx.BodyParser(req)
	if err != nil {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetErr(err).SetMessage("please check request body").Send()
	}

	logging.FromContextInfo(ctx.Context()).Interface("request", req).Send()

	var powerMeter []modelRequest.PowerMater
	for _, data := range req.PowerMeter {
		dataMeter, err := util.ParseStrToFloat(data.Data)
		if err != nil {
			return response.New(ctx).SetHttpCode(http.StatusBadRequest).
				SetErr(err).SetMessage("please check request data").Send()
		}

		date, err := util.ParseDeviceTime(data.Date, time.UTC)
		if err != nil {
			return response.New(ctx).SetHttpCode(http.StatusBadRequest).
				SetErr(err).SetMessage("please check request date").Send()
		}

		powerMeter = append(powerMeter, modelRequest.PowerMater{
			Date: date,
			Data: dataMeter,
			Name: strings.ToLower(data.Name),
		})
	}

	resp, err := h.Interactor.MeterService.Reading(ctx.Context(), powerMeter)
	if err != nil {
		return response.New(ctx).SetHttpCode(resp.HttpCode).
			SetErr(err).SetMessage("failed meter reading").Send()
	}

	return response.New(ctx).SetHttpCode(resp.HttpCode).
		SetData(resp.Data).SetMessage("success meter reading").Send()
}

func (h *Handler) Latest(ctx *fiber.Ctx) error {
	deviceID := strings.TrimSpace(ctx.Query("device_id", constant.DeviceIoT))
	if deviceID == "" {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetMessage("device_id is required").Send()
	}

	resp, err := h.Interactor.MeterService.Latest(ctx.Context(), deviceID)
	if err != nil {
		return response.New(ctx).
			SetHttpCode(fiber.StatusInternalServerError).SetErr(err).
			SetMessage("failed to get latest metrics").Send()
	}

	return response.New(ctx).SetHttpCode(http.StatusOK).
		SetData(resp.Data).SetMessage("success power latest").Send()
}

func (h *Handler) TimeSeries(ctx *fiber.Ctx) error {
	queryParams := ctx.Queries()

	metric := queryParams["metric"]
	from := queryParams["from"]
	to := queryParams["to"]
	deviceID := queryParams["device_id"]

	resp, err := h.Interactor.MeterService.TimeSeries(ctx.Context(), &modelRequest.TimeSeries{
		Metric:   metric,
		From:     from,
		To:       to,
		DeviceID: deviceID,
	})
	if err != nil {
		return response.New(ctx).SetHttpCode(resp.HttpCode).
			SetErr(err).SetMessage("failed meter reading").Send()
	}

	return response.New(ctx).SetHttpCode(resp.HttpCode).
		SetData(resp.Data).SetMessage("success power time series").Send()
}

func (h *Handler) DailyUsage(ctx *fiber.Ctx) error {
	deviceID := ctx.Query("device_id", constant.DeviceIoT)
	fromStr := ctx.Query("from", "")
	toStr := ctx.Query("to", "")

	if deviceID == "" {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetMessage("device_id is required").Send()
	}
	if fromStr == "" || toStr == "" {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetMessage("from and to are required (YYYY-MM-DD)").Send()
	}

	from, err := time.Parse(constant.DateOnly, fromStr)
	if err != nil {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetMessage("invalid from format, expected YYYY-MM-DD").Send()
	}
	to, err := time.Parse(constant.DateOnly, toStr)
	if err != nil {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetMessage("invalid to format, expected YYYY-MM-DD").Send()
	}
	if !from.Before(to) {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetMessage("from must be before to").Send()
	}

	resp, err := h.Interactor.MeterService.DailyUsage(ctx.Context(), deviceID, from, to)
	if err != nil {
		return response.New(ctx).SetHttpCode(http.StatusInternalServerError).SetErr(err).
			SetMessage("failed to get daily usage").Send()
	}

	return response.New(ctx).SetHttpCode(http.StatusOK).
		SetData(resp.Data).SetMessage("success power latest").Send()
}

func (h *Handler) Invoice(c *fiber.Ctx) error {
	deviceID := c.Query("device_id", constant.DeviceIoT)
	fromStr := c.Query("from", "")
	toStr := c.Query("to", "")
	if fromStr == "" || toStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "from and to are required (YYYY-MM-DD)")
	}

	from, err := time.Parse(constant.DateOnly, fromStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid from")
	}
	to, err := time.Parse(constant.DateOnly, toStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid to")
	}
	if !from.Before(to) {
		return fiber.NewError(fiber.StatusBadRequest, "from must be before to")
	}

	summary, lines, err := h.Interactor.MeterService.BuildBilling(c.Context(), deviceID, from, to)
	if err != nil {
		return err
	}

	pdfBytes, err := BuildInvoicePDF(summary, lines)
	if err != nil {
		return err
	}

	filename := fmt.Sprintf("invoice-%s-%s.pdf", deviceID, from.Format("2006-01"))
	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, filename))
	c.Status(http.StatusOK)
	return c.Send(pdfBytes)
}

func BuildInvoicePDF(summary modelPower.BillingSummary, lines []modelPower.DailyLine) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", gofpdf.PageSizeA4, "")
	pdf.SetTitle("Invoice", false)
	pdf.AddPage()

	// Header
	pdf.SetFont(constant.PDFFamilyStr, gofpdf.AlignBaseline, 18)
	pdf.Cell(0, 12, "Power Meter Invoice")
	pdf.Ln(10)

	pdf.SetFont(constant.PDFFamilyStr, "", 11)
	pdf.Cell(0, 6, fmt.Sprintf("Device: %s", summary.DeviceID))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Period: %s to %s", summary.From, summary.To))
	pdf.Ln(10)

	// Summary
	pdf.SetFont(constant.PDFFamilyStr, gofpdf.AlignBaseline, 12)
	pdf.Cell(0, 7, "Summary")
	pdf.Ln(8)

	pdf.SetFont(constant.PDFFamilyStr, "", 11)
	pdf.Cell(0, 6, fmt.Sprintf("Total Usage: %.6f kWh", summary.TotalKwh))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Rate: %.2f / kWh", summary.RatePerKwh))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Subtotal: %.2f", summary.Subtotal))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Tax: %.2f", summary.Tax))
	pdf.Ln(6)

	pdf.SetFont(constant.PDFFamilyStr, gofpdf.AlignBaseline, 11)
	pdf.Cell(0, 8, fmt.Sprintf("Total: %.2f", summary.Total))
	pdf.Ln(12)

	// ===== FULL WIDTH TABLE =====
	pageW, _ := pdf.GetPageSize()
	left, _, right, _ := pdf.GetMargins()
	tableW := pageW - left - right

	dayColW := tableW * 0.4
	usageColW := tableW * 0.6

	pdf.SetFont(constant.PDFFamilyStr, gofpdf.AlignBaseline, 12)
	pdf.Cell(0, 4, "Daily Usage")
	pdf.Ln(6)

	pdf.SetFont(constant.PDFFamilyStr, gofpdf.AlignBaseline, 10)
	pdf.CellFormat(dayColW, 8, "Day", "1", 0, "C", false, 0, "")
	pdf.CellFormat(usageColW, 8, "Usage (kWh)", "1", 1, "C", false, 0, "")

	pdf.SetFont(constant.PDFFamilyStr, "", 10)
	for _, l := range lines {
		pdf.CellFormat(dayColW, 8, l.Day, "1", 0, "C", false, 0, "")
		pdf.CellFormat(
			usageColW,
			8,
			fmt.Sprintf("%.6f", l.UsageKwh),
			"1",
			1,
			"R",
			false,
			0,
			"",
		)
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
