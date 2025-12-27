package response

import "github.com/gofiber/fiber/v2"

type Response struct {
	Data     interface{} `json:"data"`
	Message  string      `json:"message"`
	Error    error       `json:"error,omitempty"`
	httpCode int
	ctx      *fiber.Ctx
}
