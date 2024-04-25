package handlers

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, opt := range options {
		opt(componentHandler)
	}
	c.Set("HX-Trigger-After-Swap", "rebuild-icons")
	return adaptor.HTTPHandler(componentHandler)(c)
}

func RetargetResponse(c *fiber.Ctx, component templ.Component, hxTarget string, hxSwap string, options ...func(*templ.ComponentHandler)) error {
	c.Set("HX-Retarget", hxTarget)
	c.Set("HX-Reswap", hxSwap)
	return Render(c, component, options...)
}

func RetargetToForm(c *fiber.Ctx, resourceName string, form templ.Component, options ...func(*templ.ComponentHandler)) error {
	hxTarget := fmt.Sprintf("#%s-form", resourceName)
	return RetargetResponse(c, form, hxTarget, "outerHTML", options...)
}

func RetargetToPageHandler(c *fiber.Ctx, url string, pageHandler fiber.Handler) error {
	c.Set("HX-Location", url)
	return pageHandler(c)
}
