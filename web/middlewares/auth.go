package middlewares

import (
	"context"
	"slices"

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

		allowedPaths := []string{
			"/", "/sse/notify-operations-results",
		}

		if slices.Contains(allowedPaths, c.Path()) {
			return c.Next()
		}

		return handlers.RetargetToPageHandler(c, "/login", handlers.LoginPage)
	}

	// save user data to request context
	userData := &utils.ReqUserData{
		ID:              userID,
		IsAuthenticated: isAuthenticated,
	}
	ctx := context.WithValue(c.Context(), "UserData", userData)
	adaptor.CopyContextToFiberContext(ctx, c.Context())

	err = c.Next()
	if err != nil {
		return err
	}

	// TODO
	// A cada 24h renovar o ID da sessão

	// refresh user session
	// sess.Save must be called last
	return services.SaveUserSession(c, userID)
}
