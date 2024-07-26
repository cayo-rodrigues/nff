package handlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func Render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, opt := range options {
		opt(componentHandler)
	}

	c.Append("HX-Trigger-After-Swap", "rebuild-icons")

	isPageRequest := c.Get("HX-Boosted") == "true"
	isListRequest := strings.HasSuffix(c.Path(), "/list")

	if isPageRequest || isListRequest {
		handleBrowserQueryParams(c)
	}

	return adaptor.HTTPHandler(componentHandler)(c)
}

func RetargetResponse(c *fiber.Ctx, component templ.Component, hxTarget string, hxSwap string, options ...func(*templ.ComponentHandler)) error {
	c.Append("HX-Retarget", hxTarget)
	c.Append("HX-Reswap", hxSwap)
	return Render(c, component, options...)
}

func RetargetToForm(c *fiber.Ctx, resourceName string, form templ.Component, options ...func(*templ.ComponentHandler)) error {
	hxTarget := fmt.Sprintf("#%s-form", resourceName)
	return RetargetResponse(c, form, hxTarget, "outerHTML", options...)
}

func RetargetToPageHandler(c *fiber.Ctx, url string, pageHandler fiber.Handler) error {
	c.Append("HX-Location", url)
	return pageHandler(c)
}

func handleBrowserQueryParams(c *fiber.Ctx) {
	filters := c.Queries()
	if len(filters) > 0 {
		jsonFilters, err := json.Marshal(filters)
		if err == nil {
			c.Append("HX-Trigger-After-Settle", fmt.Sprintf(`{"append-query-params": %s}`, jsonFilters))
		}
	}
}
