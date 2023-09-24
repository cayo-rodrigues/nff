package handlers

import (
	"log"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/cayo-rodrigues/nff/web/internal/workers"
	"github.com/gofiber/fiber/v2"
)

type CancelInvoicesPage struct{}

type CancelInvoicesPageData struct {
	IsAuthenticated  bool
	InvoiceCancel    *models.InvoiceCancel
	GeneralError     string
	FormMsg          string
	FormSuccess      bool
	FormSelectFields *models.InvoiceCancelFormSelectFields
}

func (page *CancelInvoicesPage) NewEmptyData() *CancelInvoicesPageData {
	return &CancelInvoicesPageData{
		IsAuthenticated: true,
		FormSelectFields: &models.InvoiceCancelFormSelectFields{
			Entities: &[]models.Entity{},
		},
	}
}

func (page *CancelInvoicesPage) Render(c *fiber.Ctx) error {
	data := page.NewEmptyData()

	entities, err := workers.ListEntities(c.Context())
	if err != nil {
		data.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
	}

	data.FormSelectFields.Entities = entities
	data.InvoiceCancel = models.NewEmptyInvoiceCancel()

	return c.Render("invoices-cancel", data, "layouts/base")
}

func (page *CancelInvoicesPage) CancelInvoice(c *fiber.Ctx) error {
	data := page.NewEmptyData()
	data.FormSuccess = true
	data.FormMsg = "Requerimento efetuado com sucesso! Acompanhe o progresso na sessão abaixo."

	entities, err := workers.ListEntities(c.Context())
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	data.FormSelectFields.Entities = entities

	entityId, err := strconv.Atoi(c.FormValue("entity"))
	if err != nil {
		log.Println("Error converting entity id from string to int: ", err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	entity, err := workers.RetrieveEntity(c.Context(), entityId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	invoiceCancel, err := models.NewInvoiceCancelFromForm(c)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	invoiceCancel.Entity = entity

	if !invoiceCancel.IsValid() {
		data.FormMsg = "Corrija os campos abaixo."
		data.FormSuccess = false
	}

	data.InvoiceCancel = invoiceCancel

	return c.Render("partials/invoice-cancel-form", data)
}