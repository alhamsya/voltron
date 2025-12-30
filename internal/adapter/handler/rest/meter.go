package rest

import (
	modelRequest "github.com/alhamsya/voltron/internal/core/domain/request"
	"github.com/alhamsya/voltron/pkg/manager/response"
	"github.com/alhamsya/voltron/pkg/util"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
	"time"
)

func (h *Handler) Reading(ctx *fiber.Ctx) error {
	req := new(modelRequest.ReqHandlerMeterReading)
	err := ctx.BodyParser(req)
	if err != nil {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetErr(err).SetMessage("please check request body").Send()
	}

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
	deviceID := strings.TrimSpace(ctx.Query("device_id", "iot"))
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
	deviceID := ctx.Query("device_id", "iot")
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

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return response.New(ctx).SetHttpCode(http.StatusBadRequest).
			SetMessage("invalid from format, expected YYYY-MM-DD").Send()
	}
	to, err := time.Parse("2006-01-02", toStr)
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
