package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

var EntityNotFoundErr = errors.New("Entidade não encontrada")
var InternalServerErr = errors.New("Ocorreu um erro inesperado no nosso servidor. Por favor tente novamente daqui a pouco.")

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
