package middlewares

import "github.com/gofiber/fiber/v2"

func AuthMiddleware(c *fiber.Ctx) error {
	sess, err := SessionStore.Get(c)
	if err != nil {
		return err
	}

	if auth, ok := sess.Get("IsAuthenticated").(bool); !ok || !auth {
		c.Locals("IsAuthenticated", false)
		if c.Path() == "/" {
			return c.Next()
		}
		return c.Redirect("/login")
	}

	c.Locals("IsAuthenticated", true)
	return c.Next()
}
