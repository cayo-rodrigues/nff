package middlewares

import (
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

		if c.Path() != "/" {
			c.Set("HX-Target", "body")
			return c.Redirect("/login")
		}
	}

	// refresh user session
	services.SaveUserSession(c, userID)

	return c.Next()
}
