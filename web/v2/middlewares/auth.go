package middlewares

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/handlers"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
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
			return c.Next()
		case "/sse/notify-operations-results":
			return c.Next()
		default:
			return handlers.RetargetToPageHandler(c, "/login", handlers.LoginPage)

		}
	}

	// refresh user session
	services.SaveUserSession(c, userID)

	// save user data to request context
	userData := &utils.UserData{
		ID:              userID,
		IsAuthenticated: isAuthenticated,
	}
	ctx := context.WithValue(c.Context(), "UserData", userData)
	adaptor.CopyContextToFiberContext(ctx, c.Context())

	return c.Next()
}
