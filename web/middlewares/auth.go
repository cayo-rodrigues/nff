package middlewares

import (
	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	sess, err := db.SessionStore.Get(c)
	if err != nil {
		return err
	}

	isAuthenticated, authOk := sess.Get("IsAuthenticated").(bool)
	userID, idOk := sess.Get("UserID").(int)

	if !authOk || !idOk || !isAuthenticated || userID == 0 {
		c.Locals("IsAuthenticated", false)
		c.Locals("UserID", 0)
		if c.Path() != "/" {
			sess.Destroy()
			return c.Redirect("/login")
		}
	}

	// maybe we could refresh the user session here
	c.Locals("IsAuthenticated", isAuthenticated)
	c.Locals("UserID", userID)
	return c.Next()
}
