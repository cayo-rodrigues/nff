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
	i1 := models.NewInvoice()
	i2 := models.NewInvoice()
	i3 := models.NewInvoice()
	i4 := models.NewInvoice()
	i5 := models.NewInvoice()

	i1.Sender.Name = "Emerson"
	i1.Recipient.Name = "Lúcio da Silva"
	i1.ReqStatus = "success"
	i1.Number = "123.456.789"

	i2.Sender.Name = "Cayo Rodrigues"
	i2.Recipient.Name = "Ivy Rodrigues"
	i2.ReqStatus = "warning"

	i3.Sender.Name = "Joelson do nome desnecessauramente grande só pra enxer o saco"
	i3.Recipient.Name = "Oto cara com nome muito grande DISTRIBOI LTDA. MONSTROS S.A."
	i3.ReqStatus = "error"

	invoices := []*models.Invoice{
		i1, i2, i3, i4, i5,
	}
	return Render(c, layouts.Base(pages.InvoicesPage(invoices)))
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
