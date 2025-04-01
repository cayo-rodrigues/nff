package handlers

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/siare"
	"github.com/cayo-rodrigues/nff/web/ui/components"
	"github.com/cayo-rodrigues/nff/web/ui/forms"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/gofiber/fiber/v2"
)

func CancelInvoicePage(c *fiber.Ctx) error {
	filters := c.Queries()

	cancelingsList, err := services.ListCancelings(c.Context(), filters)
	if err != nil {
		return err
	}

	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}
	entitiesByType := models.NewEntitiesByType(entities)

	cancelingForForm := models.NewInvoiceCancelWithSamples(entitiesByType.Senders)
	cancelingsByDate := services.GroupListByDate(cancelingsList)
	page := pages.InvoicesCancelingsPage(cancelingsByDate, cancelingForForm, entitiesByType)

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter", "highlight-current-page", "notification-list-loaded")
	return Render(c, layouts.Base(page))
}

func CancelInvoice(c *fiber.Ctx) error {
	decryptionKey, err := services.GetEncryptionKeySession(c)
	if err != nil {
		return RetargetToReauth(c)
	}

	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}
	entitiesByType := models.NewEntitiesByType(entities)

	canceling := models.NewInvoiceCancelFromForm(c)

	entity, err := services.RetrieveEntity(c.Context(), canceling.Entity.ID)
	if err != nil {
		return err
	}
	canceling.Entity = entity

	if !canceling.IsValid() {
		return Render(c, forms.CancelInvoiceForm(canceling, entitiesByType.Senders))
	}

	err = services.CreateCanceling(c.Context(), canceling)
	if err != nil {
		return err
	}

	ssapi := siare.GetSSApiClient().WithDecryptionKey(decryptionKey)
	go ssapi.CancelInvoice(canceling)

	c.Append("HX-Trigger-After-Swap", "reload-canceling-list")
	return Render(c, forms.CancelInvoiceForm(canceling, entitiesByType.Senders))
}

func CancelInvoiceByID(c *fiber.Ctx) error {
	decryptionKey, err := services.GetEncryptionKeySession(c)
	if err != nil {
		return RetargetToReauth(c)
	}

	invoiceID, err := c.ParamsInt("invoice_id")
	if err != nil {
		return err
	}

	canceling, err := services.CreateCancelingFromInvoiceID(c.Context(), invoiceID)
	if err != nil {
		return err
	}

	ssapi := siare.GetSSApiClient().WithDecryptionKey(decryptionKey)
	go ssapi.CancelInvoice(canceling)

	return nil
}

func ListInvoiceCancelings(c *fiber.Ctx) error {
	filters := c.Queries()
	cancelings, err := services.ListCancelings(c.Context(), filters)
	if err != nil {
		return err
	}

	cancelingsByDate := services.GroupListByDate(cancelings)

	return Render(c, components.InvoicesCancelingsList(cancelingsByDate))
}

func GetCancelInvoiceForm(c *fiber.Ctx) error {
	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}
	entitiesByType := models.NewEntitiesByType(entities)

	baseCancelingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	baseCanceling, err := services.RetrieveCanceling(c.Context(), baseCancelingID)
	if err != nil {
		return err
	}

	c.Append("HX-Trigger-After-Swap", "scroll-to-top")
	return Render(c, forms.CancelInvoiceForm(baseCanceling, entitiesByType.Senders))
}

func RetrieveInvoiceCancelCard(c *fiber.Ctx) error {
	cancelingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	canceling, err := services.RetrieveCanceling(c.Context(), cancelingID)
	if err != nil {
		return err
	}

	return Render(c, components.InvoiceCancelCard(canceling))
}
