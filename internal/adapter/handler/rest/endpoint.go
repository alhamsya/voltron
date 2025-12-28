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

	flight := h.App.Group("/v1").Group("/api")
	flight.Post("/meter", h.Reading)
}
