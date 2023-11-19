package handlers

import (
	"github.com/cayo-rodrigues/nff/web/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	// email := c.FormValue("email")
	// password := c.FormValue("password")

	// get user by email

	// compare passwords
	passwordsMatch := true

	if passwordsMatch {
		sess, err := middlewares.SessionStore.Get(c)
		if err != nil {
			return err
		}
		sess.Set("IsAuthenticated", true)
		err = sess.Save()
		if err != nil {
			return err
		}
		return c.Redirect("/entities")
	}

	return c.Render("/login", fiber.Map{}, "layouts/base")
}

func Logout(c *fiber.Ctx) error {
	sess, err := middlewares.SessionStore.Get(c)
	if err != nil {
		return err
	}
	sess.Set("IsAuthenticated", false)
	err = sess.Save()
	if err != nil {
		return err
	}
	return c.Redirect("/")
}
