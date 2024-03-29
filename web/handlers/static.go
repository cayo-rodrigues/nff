package handlers

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

func ServeStyles(ctx *fiber.Ctx) error {
	stylesheet := ctx.Params("stylesheet")

	filepath := fmt.Sprintf("static/styles/%s", stylesheet)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		msg := fmt.Sprintf("File '%s' does not exist.\n", filepath)
		return ctx.Status(fiber.StatusNotFound).Send([]byte(msg))
	}

	return ctx.SendFile(filepath)
}

func ServeJS(ctx *fiber.Ctx) error {
	script := ctx.Params("script")

	filepath := fmt.Sprintf("static/scripts/%s", script)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		msg := fmt.Sprintf("File '%s' does not exist.\n", filepath)
		return ctx.Status(fiber.StatusNotFound).Send([]byte(msg))
	}

	return ctx.SendFile(filepath)
}

func ServeIcons(ctx *fiber.Ctx) error {
	icon := ctx.Params("icon")

	filepath := fmt.Sprintf("static/icons/%s", icon)
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		msg := fmt.Sprintf("File '%s' does not exist.\n", filepath)
		return ctx.Status(fiber.StatusNotFound).Send([]byte(msg))
	}

	return ctx.SendFile(filepath)
}
