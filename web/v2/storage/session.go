package storage

import (
	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/gofiber/fiber/v2"
)

type SessionOpts struct {
	Key string
	Val any
}

func SetSessionKVs(c *fiber.Ctx, opts ...*SessionOpts) error {
	db := database.GetDB()

	sess, err := db.SessionStore.Get(c)
	if err != nil {
		return err
	}
	for _, opt := range opts {
		sess.Set(opt.Key, opt.Val)
	}
	return sess.Save()
}

func DeleteSessionKeys(c *fiber.Ctx, opts ...*SessionOpts) error {
	db := database.GetDB()

	sess, err := db.SessionStore.Get(c)
	if err != nil {
		return err
	}
	for _, opt := range opts {
		sess.Delete(opt.Key)
	}
	return sess.Save()
}
