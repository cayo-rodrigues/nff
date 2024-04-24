package handlers

import (
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/forms"
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
	if err != nil {
		return err
	}
	cancelingForForm := models.NewInvoiceCancelWithSamples(entities)

	return Render(c, layouts.Base(pages.InvoicesCancelingsPage(cancelingsList, cancelingForForm, entities)))
}

func CancelInvoice(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)

	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	canceling := models.NewInvoiceCancelFromForm(c)

	entity, err := services.RetrieveEntity(c.Context(), canceling.Entity.ID, userID)
	if err != nil {
		return err
	}
	canceling.Entity = entity

	if canceling.IsValid() {
		c.Set("HX-Trigger-After-Swap", "reload-cancelings-list")
	}

	return Render(c, forms.CancelInvoiceForm(canceling, entities))
}
