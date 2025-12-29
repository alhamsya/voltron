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
		SetData(resp).SetMessage("success meter reading").Send()
}

func (h *Handler) Latest(ctx *fiber.Ctx) error {
	return response.New(ctx).SetHttpCode(http.StatusOK).
		SetData(nil).SetMessage("success power latest").Send()
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
			SetErr(err).SetMessage("failed power time series").Send()
	}

	return response.New(ctx).SetHttpCode(resp.HttpCode).
		SetData(resp).SetMessage("success power time series").Send()
}

func (h *Handler) DailyUsage(ctx *fiber.Ctx) error {
	return response.New(ctx).SetHttpCode(http.StatusOK).
		SetData(nil).SetMessage("success power latest").Send()
}
