package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var EntityNotFoundErr = errors.New(EntityNotFoundMsg)
var InvoiceNotFoundErr = errors.New(InvoiceNotFoundMsg)
var CancelingNotFoundErr = errors.New(CancelingNotFoundMsg)
var MetricsNotFoundErr = errors.New(MetricsNotFoundMsg)
var MetricsResultNotFoundErr = errors.New(MetricsResultNotFoundMsg)
var PrintingNotFoundErr = errors.New(PrintingNotFoundMsg)
var UserNotFoundErr = errors.New(UserNotFoundMsg)
var InvalidLoginDataErr = errors.New(InvalidLoginDataMsg)

var InternalServerErr = errors.New(InternalServerErrMsg)

func ErrorResponse(c *fiber.Ctx, tmplName string, tmplData interface{}, event string) error {
	// w.WriteHeader(statusCode)
	c.Append("HX-Trigger-After-Settle", event)
	return c.Render(tmplName, tmplData)
}

func GeneralErrorResponse(c *fiber.Ctx, err error) error {
	c.Append("HX-Retarget", "#general-error-msg")
	c.Append("HX-Reswap", "outerHTML")
	return ErrorResponse(c, "partials/general-error", err.Error(), "general-error")
}
