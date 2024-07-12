package handlers

import (
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/components"
	"github.com/cayo-rodrigues/nff/web/ui/shared"
	"github.com/gofiber/fiber/v2"
)

func ListNotifications(c *fiber.Ctx) error {
	notificationItems := services.GetNotifications(c.Context())
	c.Append("HX-Trigger-After-Settle", "notification-list-loaded")
	return Render(c, shared.NotificationList(notificationItems))
}

func ClearNotifications(c *fiber.Ctx) error {
	err := services.ClearNotifications(c.Context())
	if err != nil {
		return err
	}
	c.Append("HX-Trigger-After-Settle", "notification-list-cleared")
	return Render(c, components.Nothing())
}
