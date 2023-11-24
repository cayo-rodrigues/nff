package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
)

type InvoicesPage struct {
	service       interfaces.InvoiceService
	entityService interfaces.EntityService
	siareBGWorker interfaces.SiareBGWorker
}

func NewInvoicesPage(service interfaces.InvoiceService, entityService interfaces.EntityService, siareBGWorker interfaces.SiareBGWorker) *InvoicesPage {
	return &InvoicesPage{
		service:       service,
		entityService: entityService,
		siareBGWorker: siareBGWorker,
	}
}

type InvoicesPageData struct {
	IsAuthenticated  bool
	Invoices         []*models.Invoice
	Invoice          *models.Invoice
	GeneralError     string
	FormMsg          string
	FormSelectFields *models.InvoiceFormSelectFields
}

func (p *InvoicesPage) NewEmptyData() *InvoicesPageData {
	return &InvoicesPageData{
		IsAuthenticated: true,
		FormSelectFields: &models.InvoiceFormSelectFields{
			Operations:   &globals.InvoiceOperations,
			Cfops:        &globals.InvoiceCfops,
			BooleanField: &globals.InvoiceBooleanField,
			IcmsOptions:  &globals.InvoiceIcmsOptions,
		},
	}
}

func (p *InvoicesPage) Render(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()
	userID := c.Locals("UserID").(int)

	// TODO async data aggregation with go routines
	entities, err := p.entityService.ListEntities(c.Context(), userID)
	if err != nil {
		pageData.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
		return c.Render("invoices", pageData, "layouts/base")
	}

	pageData.FormSelectFields.Entities = entities
	pageData.Invoice = models.NewEmptyInvoice()

	// get the latest 10 invoices
	invoices, err := p.service.ListInvoices(c.Context(), userID)
	if err != nil {
		pageData.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
		return c.Render("invoices", pageData, "layouts/base")
	}

	pageData.Invoices = invoices

	return c.Render("invoices", pageData, "layouts/base")
}

func (p *InvoicesPage) RequireInvoice(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()
	userID := c.Locals("UserID").(int)

	// TODO async data aggregation with go routines

	entities, err := p.entityService.ListEntities(c.Context(), userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	pageData.FormSelectFields.Entities = entities

	senderID, err := strconv.Atoi(c.FormValue("sender"))
	if err != nil {
		log.Println("Error converting sender id from string to int: ", err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}
	recipientID, err := strconv.Atoi(c.FormValue("recipient"))
	if err != nil {
		log.Println("Error converting recipient id from string to int: ", err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	sender, err := p.entityService.RetrieveEntity(c.Context(), senderID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	recipient, err := p.entityService.RetrieveEntity(c.Context(), recipientID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	invoice, err := models.NewInvoiceFromForm(c)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	invoice.Sender = sender
	invoice.Recipient = recipient
	invoice.CreatedBy = userID

	if !invoice.IsValid() {
		pageData.FormMsg = "Corrija os campos abaixo."
		pageData.Invoice = invoice
		return utils.RetargetToForm(c, "invoice", pageData)
	}

	err = p.service.CreateInvoice(c.Context(), invoice)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	go p.siareBGWorker.RequestInvoice(invoice)

	c.Set("HX-Trigger-After-Settle", "invoice-required")
	return c.Render("partials/request-card", invoice)
}

func (p *InvoicesPage) GetItemFormSection(c *fiber.Ctx) error {
	item := models.NewEmptyInvoiceItem()
	c.Set("HX-Trigger-After-Settle", "enumerate-item-sections")
	return c.Render("partials/invoice-form-item-section", item)
}

func (p *InvoicesPage) GetRequestCardDetails(c *fiber.Ctx) error {
	invoiceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.InvoiceNotFoundErr)
	}

	userID := c.Locals("UserID").(int)

	invoice, err := p.service.RetrieveInvoice(c.Context(), invoiceID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Trigger-After-Settle", "open-request-card-details")
	return c.Render("partials/request-card-details", invoice)
}

func (p *InvoicesPage) GetInvoiceForm(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()
	userID := c.Locals("UserID").(int)

	// TODO async data aggregation with go routines

	entities, err := p.entityService.ListEntities(c.Context(), userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	pageData.FormSelectFields.Entities = entities

	invoiceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.InvoiceNotFoundErr)
	}
	invoice, err := p.service.RetrieveInvoice(c.Context(), invoiceID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	pageData.Invoice = invoice

	c.Set("HX-Trigger-After-Settle", "scroll-to-top")
	return c.Render("partials/invoice-form", pageData)
}

func (p *InvoicesPage) GetRequestStatus(c *fiber.Ctx) error {
	invoiceID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.InvoiceNotFoundErr)
	}

	key := fmt.Sprintf("reqstatus:invoice:%v", invoiceID)
	err = db.Redis.GetDel(c.Context(), key).Err()
	if err == redis.Nil {
		return c.Render("partials/request-card-status", "pending")
	}
	if err != nil {
		log.Printf("Error reading redis key %v: %v\n", key, err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	userID := c.Locals("UserID").(int)

	invoice, err := p.service.RetrieveInvoice(c.Context(), invoiceID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	targetId := fmt.Sprintf("#request-card-%v", c.Params("id"))
	c.Set("HX-Retarget", targetId)
	c.Set("HX-Reswap", "outerHTML")
	return c.Render("partials/request-card", invoice)
}
