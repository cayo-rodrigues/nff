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
	IsAuthenticated   bool
	InvoiceCancel     *models.InvoiceCancel
	InvoiceCancelings []*models.InvoiceCancel
	GeneralError      string
	FormMsg           string
	FormSuccess       bool
	FormSelectFields  *models.InvoiceCancelFormSelectFields
}

func (page *CancelInvoicesPage) NewEmptyData() *CancelInvoicesPageData {
	return &CancelInvoicesPageData{
		IsAuthenticated: true,
		FormSelectFields: &models.InvoiceCancelFormSelectFields{
			Entities: []*models.Entity{},
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

	// get the latest 10 cancelings
	cancelings, err := workers.ListInvoiceCancelings(c.Context())
	if err != nil {
		data.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
		return c.Render("invoices-cancel", data, "layouts/base")
	}

	data.InvoiceCancelings = cancelings

	return c.Render("invoices-cancel", data, "layouts/base")
}

func (page *CancelInvoicesPage) CancelInvoice(c *fiber.Ctx) error {
	data := page.NewEmptyData()
	// formSuccess := true

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

	data.InvoiceCancel = invoiceCancel
	data.InvoiceCancel.Entity = entity

	if !invoiceCancel.IsValid() {
		data.FormMsg = "Corrija os campos abaixo."
		// formSuccess = false
		c.Set("HX-Retarget", "#invoice-cancel-form")
		c.Set("HX-Reswap", "outerHTML")
		return c.Render("partials/invoice-cancel-form", data)
	}

	err = workers.CreateInvoiceCanceling(c.Context(), data.InvoiceCancel)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Trigger-After-Settle", "invoice-cancel-required")
	return c.Render("partials/request-card", data.InvoiceCancel)
}

func (page *CancelInvoicesPage) GetRequestCardDetails(c *fiber.Ctx) error {
	cancelingId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.CancelingNotFoundErr)
	}
	canceling, err := workers.RetrieveInvoiceCanceling(c.Context(), cancelingId)

	c.Set("HX-Trigger-After-Settle", "open-request-card-details")
	return c.Render("partials/request-card-details", canceling)
}

func (page *CancelInvoicesPage) GetInvoiceCancelForm(c *fiber.Ctx) error {
	data := page.NewEmptyData()

	entities, err := workers.ListEntities(c.Context())
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	data.FormSelectFields.Entities = entities

	cancelingId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.CancelingNotFoundErr)
	}
	canceling, err := workers.RetrieveInvoiceCanceling(c.Context(), cancelingId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	data.InvoiceCancel = canceling

	c.Set("HX-Trigger-After-Settle", "scroll-to-top")
	return c.Render("partials/invoice-cancel-form", data)
}
