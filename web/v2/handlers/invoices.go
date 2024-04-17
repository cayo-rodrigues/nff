package handlers

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/ui/shared"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func InvoicesPage(c *fiber.Ctx) error {
	invoice := models.NewInvoice()
	return Render(c, layouts.Base(pages.InvoicesPage(invoice)))
}

func CreateInvoicePage(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	invoice := models.NewInvoice()
	if len(entities) > 0 {
		invoice.Sender = entities[0]
	}
	invoice.Items = append(invoice.Items, models.NewInvoiceItem())
	return Render(c, layouts.Base(pages.InvoiceFormPage(invoice, entities)))
}

func GetSenderIeInput(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	entityID, err := strconv.Atoi(c.Query("sender"))
	if err != nil {
		return err
	}
	entity, err := services.RetrieveEntity(c.Context(), entityID, userID)

	return Render(c, shared.SelectInput(&shared.InputData{
		ID:      "available-ies",
		Label:   "IE do Remetente",
		Value:   entity.Ie,
		Options: &shared.InputOptions{StringOptions: entity.AllIes()},
	}))
}
