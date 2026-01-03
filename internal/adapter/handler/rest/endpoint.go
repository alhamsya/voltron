package rest

import (
	"net/http"

	"github.com/alhamsya/voltron/internal/core/domain/constant"
	"github.com/benbjohnson/clock"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetRoot(ctx *fiber.Ctx) error {
	resp := map[string]any{
		"message": "success",
		"time": map[string]any{
			"utc":     clock.New().Now(),
			"jakarta": clock.New().Now().Format(constant.DateTime),
		},
	}
	return ctx.Status(http.StatusOK).JSON(resp)
}

func (h *Handler) Register() {
	h.App.Get("/", h.GetRoot)

	api := h.App.Group("/v1").Group("/api")
	powerMeter := api.Group("power")
	billing := api.Group("billing")

	powerMeter.Post("/meter", h.Reading)
	powerMeter.Get("/latest", h.Latest)
	powerMeter.Get("/time-series", h.TimeSeries)
	powerMeter.Get("/daily-usage", h.DailyUsage)

	billing.Get("/invoice", h.Invoice)
}
