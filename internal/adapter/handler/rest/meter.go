package rest

import (
	"github.com/alhamsya/voltron/internal/core/domain/constant"
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

		date, err := time.ParseInLocation(constant.DateWithSlash, data.Date, time.UTC)
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
