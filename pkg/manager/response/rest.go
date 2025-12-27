package response

import (
	"github.com/gofiber/fiber/v2"
)

func New(ctx *fiber.Ctx) *Response {
	return &Response{
		ctx: ctx,
	}
}

func (r *Response) SetMessage(message string) *Response {
	r.Message = message
	return r
}

func (r *Response) SetData(data any) *Response {
	r.Data = data
	return r
}

func (r *Response) SetErr(err error) *Response {
	r.Error = err
	return r
}

func (r *Response) SetHttpCode(httpCode int) *Response {
	r.httpCode = httpCode
	return r
}

func (r *Response) Send(arg ...string) (resp error) {
	//args := strings.Join(arg, "|")

	//valida http code
	if r.httpCode <= 0 {
		r.httpCode = fiber.StatusInternalServerError
	}

	if r.httpCode < fiber.StatusContinue {
		r.httpCode = fiber.StatusOK
	}

	//validate message for http code
	switch r.httpCode / 100 {
	case fiber.StatusOK / 100:
		r.Message = r.Message + " successfully"
	case fiber.StatusBadRequest / 100:
		////replace message from args
		//if strings.TrimSpace(args) != "" {
		//	r.Message = args
		//}
		//
		//if strings.TrimSpace(args) == "" && r.error != nil {
		//	r.Message = errors.Cause(r.error).Error()
		//}

	case fiber.StatusInternalServerError / 100:
		r.Message = "please try again"
	}

	resp = r.ctx.Status(r.httpCode).JSON(&r)
	return resp
}
