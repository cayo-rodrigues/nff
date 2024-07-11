package handlers

import (
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/shared"
	"github.com/gofiber/fiber/v2"
)

func ListNotifications(c *fiber.Ctx) error {
	notificationItems := services.GetNotifications(c.Context())
	return Render(c, shared.NotificationList(notificationItems))
}

func ClearNotifications(c *fiber.Ctx) error {
	return nil
}
