package services

import (
	"fmt"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func SaveUserSession(c *fiber.Ctx, userData *utils.UserCtxData) error {
	userData.IsAuthenticated = true
	sessionOpts := []*storage.SessionOpts{
		{Key: "UserData", Val: userData},
	}

	err := storage.SetSessionKVs(c, sessionOpts...)
	if err != nil {
		return err
	}

	return nil
}

func DestroyUserSession(c *fiber.Ctx) error {
	sessionOpts := []*storage.SessionOpts{
		{Key: "UserData"},
	}

	err := storage.DeleteSessionKeys(c, sessionOpts...)
	if err != nil {
		return err
	}

	return nil
}

func GetUserSession(c *fiber.Ctx) (*utils.UserCtxData, error) {
	sessionOpts := []*storage.SessionOpts{
		{Key: "UserData"},
	}

	vals, err := storage.GetSessionValsByKeys(c, sessionOpts...)
	if err != nil {
		return nil, err
	}

	userData, ok := vals["UserData"].(*utils.UserCtxData)
	if !ok {
		return &utils.UserCtxData{}, nil
	}

	return userData, nil
}

func SaveEncryptionKeySession(c *fiber.Ctx, key []byte) error {
	db := database.GetDB()

	sess, err := db.SessionStore.Get(c)
	if err != nil {
		return err
	}

	id := sess.ID()
	sessKey := fmt.Sprintf("key:%s", id)

	return db.Redis.Set(c.Context(), sessKey, key, 4*time.Hour).Err()
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
