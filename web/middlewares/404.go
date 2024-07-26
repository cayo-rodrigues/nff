package middlewares

import (
	// "net/http"

	// "github.com/a-h/templ"
	// "github.com/cayo-rodrigues/nff/web/components"
	// "github.com/cayo-rodrigues/nff/web/handlers"
	"github.com/gofiber/fiber/v2"
)

func NotFoundMiddleware(c *fiber.Ctx) error {
	return nil
	// return handlers.Render(c, components.NotFound(), templ.WithStatus(http.StatusNotFound))
}
