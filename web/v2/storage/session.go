package storage

import (
	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/gofiber/fiber/v2"
)

type SessionOpts struct {
	Key string
	Val any
}

func SaveSession(c *fiber.Ctx, opts ...*SessionOpts) error {
	db := database.GetDB()

	sess, err := db.SessionStore.Get(c)
	if err != nil {
		return err
	}
	for _, opt := range opts {
		sess.Set(opt.Key, opt.Val)
	}
	err = sess.Save()
	if err != nil {
		return err
	}

	return nil
}
