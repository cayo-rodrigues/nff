package utils

import (
	"errors"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/gofiber/fiber/v2"
)

var EntityNotFoundErr = errors.New(globals.EntityNotFoundMsg)
var InvoiceNotFoundErr = errors.New(globals.InvoiceNotFoundMsg)
var CancelingNotFoundErr = errors.New(globals.CancelingNotFoundMsg)
var MetricsNotFoundErr = errors.New(globals.MetricsNotFoundMsg)
var PrintingNotFoundErr = errors.New(globals.PrintingNotFoundMsg)
var UserNotFoundErr = errors.New(globals.UserNotFoundMsg)

var InternalServerErr = errors.New(globals.InternalServerErrMsg)

func ErrorResponse(c *fiber.Ctx, tmplName string, tmplData interface{}, event string) error {
	// w.WriteHeader(statusCode)
	c.Set("HX-Trigger-After-Settle", event)
	return c.Render(tmplName, tmplData)
}

func GeneralErrorResponse(c *fiber.Ctx, err error) error {
	c.Set("HX-Retarget", "#general-error-msg")
	c.Set("HX-Reswap", "outerHTML")
	return ErrorResponse(c, "partials/general-error", err.Error(), "general-error")
}
