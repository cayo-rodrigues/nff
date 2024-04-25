package handlers

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/components"
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
		err := services.CreateCanceling(c.Context(), canceling, userID)
		if err != nil {
			return err
		}
		c.Set("HX-Trigger-After-Settle", "reload-canceling-list")
	}

	return Render(c, forms.CancelInvoiceForm(canceling, entities))
}

func ListInvoiceCancelings(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	cancelings, err := services.ListCancelings(c.Context(), userID)
	if err != nil {
		return err
	}

	return Render(c, components.InvoicesCancelingsList(cancelings))
}

func GetCancelInvoiceForm(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)

	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	baseCancelingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	baseCanceling, err := services.RetrieveCanceling(c.Context(), baseCancelingID, userID)
	if err != nil {
		return err
	}

	return Render(c, forms.CancelInvoiceForm(baseCanceling, entities))
}
