package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func RetargetResponse(c *fiber.Ctx, tmplName string, tmplData any, hxTarget string, hxSwap string) error {
	c.Set("HX-Retarget", hxTarget)
	c.Set("HX-Reswap", hxSwap)
	return c.Render(tmplName, tmplData)
}

func RetargetToForm(c *fiber.Ctx, resourceName string, tmplData any) error {
	tmplName := fmt.Sprintf("partials/forms/%s-form", resourceName)
	hxTarget := fmt.Sprintf("#%s-form", resourceName)
	return RetargetResponse(c, tmplName, tmplData, hxTarget, "outerHTML")
}
