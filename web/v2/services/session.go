package services

import (
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/gofiber/fiber/v2"
)

func SaveUserSession(c *fiber.Ctx, userID int) error {
	sessionOpts := []*storage.SessionOpts{
		{Key: "IsAuthenticated", Val: true},
		{Key: "UserID", Val: userID},
	}
	return storage.SetSessionKVs(c, sessionOpts...)
}

func DestroyUserSession(c *fiber.Ctx) error {
	sessionOpts := []*storage.SessionOpts{
		{Key: "IsAuthenticated"},
		{Key: "UserID"},
	}
	return storage.DeleteSessionKeys(c, sessionOpts...)
}
