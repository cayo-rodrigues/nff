package services

import (
	"fmt"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/gofiber/fiber/v2"
)

func SaveUserSession(c *fiber.Ctx, userID int) error {
	sessionOpts := []*storage.SessionOpts{
		{Key: "IsAuthenticated", Val: true},
		{Key: "UserID", Val: userID},
	}

	err := storage.SetSessionKVs(c, sessionOpts...)
	if err != nil {
		return err
	}

	c.Locals("IsAuthenticated", true)
	c.Locals("UserID", userID)

	return nil
}

func DestroyUserSession(c *fiber.Ctx) error {
	sessionOpts := []*storage.SessionOpts{
		{Key: "IsAuthenticated"},
		{Key: "UserID"},
	}

	err := storage.DeleteSessionKeys(c, sessionOpts...)
	if err != nil {
		return err
	}

	c.Locals("IsAuthenticated", false)
	c.Locals("UserID", 0)

	return nil
}

func GetUserSession(c *fiber.Ctx) (isAuthenticated bool, userID int, err error) {
	sessionOpts := []*storage.SessionOpts{
		{Key: "IsAuthenticated"},
		{Key: "UserID"},
	}

	vals, err := storage.GetSessionValsByKeys(c, sessionOpts...)
	if err != nil {
		return false, 0, err
	}

	isAuthenticated, authOk := vals["IsAuthenticated"].(bool)
	userID, idOK := vals["UserID"].(int)

	if !authOk || !idOK {
		return false, 0, nil
	}

	return isAuthenticated, userID, nil
}

func SaveEncryptionKeySession(c *fiber.Ctx, key []byte) error {
	db := database.GetDB()

	sess, err := db.SessionStore.Get(c)
	if err != nil {
		return err
	}

	id := sess.ID()
	sessKey := fmt.Sprintf("key:%s", id)

	return db.Redis.Set(c.Context(), sessKey, key, 4 * time.Hour).Err()
}

func GetEncryptionKeySession(c *fiber.Ctx) ([]byte, error) {
	db := database.GetDB()

	sess, err := db.SessionStore.Get(c)
	if err != nil {
		return nil, err
	}

	id := sess.ID()
	sessKey := fmt.Sprintf("key:%s", id)

	key, err := db.Redis.Get(c.Context(), sessKey).Bytes()
	if err != nil {
		return nil, err
	}

	return key, nil
}
