package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/gofiber/fiber/v2"
)

func NotFoundPage(c *fiber.Ctx) error {
	return Render(c, layouts.Base(pages.NotFoundPage()), templ.WithStatus(http.StatusNotFound))
}
