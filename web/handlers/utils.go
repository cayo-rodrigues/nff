package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/a-h/templ"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/ui/forms"
	"github.com/cayo-rodrigues/nff/web/utils"
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
	isSearchRequest := strings.HasSuffix(c.Path(), "/search")

	if isPageRequest || isListRequest || isSearchRequest {
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

func RetargetToReauth(c *fiber.Ctx) error {
	userID := utils.GetUserID(c.Context())
	log.Printf("EncryptionKey not found for user %d, asking for reauthentication...\n", userID)
	return RetargetResponse(c, forms.ReauthenticateForm(models.NewUser()), "#reauth-form-dialog-content", "innerHTML")
}

func handleBrowserQueryParams(c *fiber.Ctx) {
	filters := c.Queries()
	if len(filters) > 0 {
		jsonFilters, err := json.Marshal(filters)
		if err == nil {
			c.Append("HX-Trigger-After-Settle", fmt.Sprintf(`{"append-query-params": { "queries": %s }}`, jsonFilters))
		}
	}
}
