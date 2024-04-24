package handlers

import (
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func InvoicesCancelingsPage(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	cancelingsList, err := services.ListCancelings(c.Context(), userID)
	if err != nil {
		return err
	}

	entities, err := services.ListEntities(c.Context(), userID)
	cancelingForForm := models.NewInvoiceCancelWithSamples(entities)

	return Render(c, layouts.Base(pages.InvoicesCancelingsPage(cancelingsList, cancelingForForm, entities)))
}
