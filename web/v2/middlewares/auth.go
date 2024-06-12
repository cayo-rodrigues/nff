package middlewares

import (
	"github.com/cayo-rodrigues/nff/web/handlers"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	isAuthenticated, userID, err := services.GetUserSession(c)
	if err != nil {
		return err
	}

	if !isAuthenticated || userID == 0 {
		err := services.DestroyUserSession(c)
		if err != nil {
			return err
		}

		switch c.Path() {
		case "/":
			return handlers.RetargetToPageHandler(c, "/", handlers.HomePage)
		default:
			return handlers.RetargetToPageHandler(c, "/login", handlers.LoginPage)

		}
	}

	// refresh user session
	services.SaveUserSession(c, userID)

	return c.Next()
}
