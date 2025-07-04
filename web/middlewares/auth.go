package middlewares

import (
	"context"
	"slices"

	"github.com/cayo-rodrigues/nff/web/handlers"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func AuthMiddleware(c *fiber.Ctx) error {
	userData, err := services.GetUserSession(c)
	if err != nil {
		return err
	}

	// save user data to request context
	ctx := context.WithValue(c.Context(), "UserData", userData)
	adaptor.CopyContextToFiberContext(ctx, c.Context())

	if !userData.IsAuthenticated || userData.ID == 0 {
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

	// if userData.IsBlocked || !userData.HasChosenPaymentPlan {
	// 	pathsToSkip := []string{
	// 		"/", "/prices", "/sse/notify-operations-results", "/logout",
	// 	}

	// 	if !slices.Contains(pathsToSkip, c.Path()) {
	// 		return handlers.RetargetToPageHandler(c, "/prices", handlers.PricesPage)
	// 	}
	// }

	err = c.Next()
	if err != nil {
		return err
	}

	// TODO
	// A cada 24h renovar o ID da sessão

	// refresh user session
	// sess.Save must be called last
	return services.SaveUserSession(c, userData)
}
