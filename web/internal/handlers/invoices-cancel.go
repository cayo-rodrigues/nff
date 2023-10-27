package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
)

type CancelInvoicesPage struct {
	service       interfaces.CancelingService
	entityService interfaces.EntityService
	siareBGWorker interfaces.SiareBGWorker
}

func NewCancelInvoicesPage(service interfaces.CancelingService, entityService interfaces.EntityService, siareBGWorker interfaces.SiareBGWorker) *CancelInvoicesPage {
	return &CancelInvoicesPage{
		service:       service,
		entityService: entityService,
		siareBGWorker: siareBGWorker,
	}
}

type CancelInvoicesPageData struct {
	IsAuthenticated   bool
	InvoiceCancel     *models.InvoiceCancel
	InvoiceCancelings []*models.InvoiceCancel
	GeneralError      string
	FormMsg           string
	FormSuccess       bool
	FormSelectFields  *models.InvoiceCancelFormSelectFields
}

func (p *CancelInvoicesPage) NewEmptyData() *CancelInvoicesPageData {
	return &CancelInvoicesPageData{
		IsAuthenticated: true,
		FormSelectFields: &models.InvoiceCancelFormSelectFields{
			Entities: []*models.Entity{},
		},
	}
}

func (p *CancelInvoicesPage) Render(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()

	entities, err := p.entityService.ListEntities(c.Context())
	if err != nil {
		pageData.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
	}

	pageData.FormSelectFields.Entities = entities
	pageData.InvoiceCancel = models.NewEmptyInvoiceCancel()

	// get the latest 10 cancelings
	cancelings, err := p.service.ListInvoiceCancelings(c.Context())
	if err != nil {
		pageData.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
		return c.Render("invoices-cancel", pageData, "layouts/base")
	}

	pageData.InvoiceCancelings = cancelings

	return c.Render("invoices-cancel", pageData, "layouts/base")
}

func (p *CancelInvoicesPage) CancelInvoice(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()

	entities, err := p.entityService.ListEntities(c.Context())
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	pageData.FormSelectFields.Entities = entities

	entityId, err := strconv.Atoi(c.FormValue("entity"))
	if err != nil {
		log.Println("Error converting entity id from string to int: ", err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	entity, err := p.entityService.RetrieveEntity(c.Context(), entityId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	invoiceCancel, err := models.NewInvoiceCancelFromForm(c)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	invoiceCancel.Entity = entity

	if !invoiceCancel.IsValid() {
		pageData.InvoiceCancel = invoiceCancel
		pageData.FormMsg = "Corrija os campos abaixo."
		c.Set("HX-Retarget", "#invoice-cancel-form")
		c.Set("HX-Reswap", "outerHTML")
		return c.Render("partials/invoice-cancel-form", pageData)
	}

	err = p.service.CreateInvoiceCanceling(c.Context(), invoiceCancel)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	go p.siareBGWorker.RequestInvoiceCanceling(invoiceCancel)

	c.Set("HX-Trigger-After-Settle", "invoice-cancel-required")
	return c.Render("partials/request-card", invoiceCancel)
}

func (p *CancelInvoicesPage) GetRequestCardDetails(c *fiber.Ctx) error {
	cancelingId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.CancelingNotFoundErr)
	}
	canceling, err := p.service.RetrieveInvoiceCanceling(c.Context(), cancelingId)

	c.Set("HX-Trigger-After-Settle", "open-request-card-details")
	return c.Render("partials/request-card-details", canceling)
}

func (p *CancelInvoicesPage) GetInvoiceCancelForm(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()

	entities, err := p.entityService.ListEntities(c.Context())
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	pageData.FormSelectFields.Entities = entities

	cancelingId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.CancelingNotFoundErr)
	}
	canceling, err := p.service.RetrieveInvoiceCanceling(c.Context(), cancelingId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	pageData.InvoiceCancel = canceling

	c.Set("HX-Trigger-After-Settle", "scroll-to-top")
	return c.Render("partials/invoice-cancel-form", pageData)
}

func (p *CancelInvoicesPage) GetRequestStatus(c *fiber.Ctx) error {
	cancelingId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.CancelingNotFoundErr)
	}

	key := fmt.Sprintf("reqstatus:canceling:%v", cancelingId)
	err = db.Redis.GetDel(c.Context(), key).Err()
	if err == redis.Nil {
		return c.Render("partials/request-card-status", "pending")
	}
	if err != nil {
		log.Printf("Error reading redis key %v: %v\n", key, err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	canceling, err := p.service.RetrieveInvoiceCanceling(c.Context(), cancelingId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	targetId := fmt.Sprintf("#request-card-%v", c.Params("id"))
	c.Set("HX-Retarget", targetId)
	c.Set("HX-Reswap", "outerHTML")
	return c.Status(286).Render("partials/request-card", canceling)
}
