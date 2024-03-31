package middlewares

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/cayo-rodrigues/nff/web/components/layouts"
	"github.com/cayo-rodrigues/nff/web/components/pages"
	"github.com/cayo-rodrigues/nff/web/handlers"
	"github.com/gofiber/fiber/v2"
)

func NotFoundMiddleware(c *fiber.Ctx) error {
	return handlers.Render(c, layouts.Base(pages.NotFoundPage()), templ.WithStatus(http.StatusNotFound))
}
