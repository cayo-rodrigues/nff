package utils

import "github.com/gofiber/fiber/v2"

func GeneralInfoResponse(c *fiber.Ctx, msg string) error {
	c.Set("HX-Retarget", "#general-info-msg")
	c.Set("HX-Reswap", "outerHTML")
	c.Set("HX-Trigger-After-Settle", "general-info")
	return c.Render("partials/general-info", msg)
}
