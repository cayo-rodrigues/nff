package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type PrintInvoicesPage struct {
	service       interfaces.PrintingService
	entityService interfaces.EntityService
	siareBGWorker interfaces.SiareBGWorker
}

func NewPrintInvoicesPage(service interfaces.PrintingService, entityService interfaces.EntityService, siareBGWorker interfaces.SiareBGWorker) *PrintInvoicesPage {
	return &PrintInvoicesPage{
		service:       service,
		entityService: entityService,
		siareBGWorker: siareBGWorker,
	}
}

type PrintInvoicesPageData struct {
	IsAuthenticated  bool
	Filters          *models.ReqCardFilters
	InvoicePrint     *models.InvoicePrint
	InvoicePrintings []*models.InvoicePrint
	GeneralError     string
	FormMsg          string
	FormSuccess      bool
	FormSelectFields *models.InvoicePrintFormSelectFields
	ResourceName     string
}

func (p *PrintInvoicesPage) NewEmptyData() *PrintInvoicesPageData {
	return &PrintInvoicesPageData{
		IsAuthenticated:  true,
		Filters:          models.NewRequestCardFilters(),
		FormSelectFields: models.NewInvoicePrintFromSelectFields(),
		ResourceName:     "invoices/print",
	}
}

func (p *PrintInvoicesPage) Render(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()

	userID := c.Locals("UserID").(int)

	entities, err := p.entityService.ListEntities(c.Context(), userID)
	if err != nil {
		pageData.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
	}

	pageData.FormSelectFields.Entities = entities
	pageData.InvoicePrint = models.NewEmptyInvoicePrint()

	// get the latest 10 printings
	printings, err := p.service.ListInvoicePrintings(c.Context(), userID, nil)
	if err != nil {
		pageData.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
		return c.Render("invoices-print", pageData, "layouts/base")
	}

	pageData.InvoicePrintings = printings

	return c.Render("invoices-print", pageData, "layouts/base")
}

func (p *PrintInvoicesPage) PrintInvoice(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()
	userID := c.Locals("UserID").(int)

	entities, err := p.entityService.ListEntities(c.Context(), userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	pageData.FormSelectFields.Entities = entities

	entityID, err := strconv.Atoi(c.FormValue("entity"))
	if err != nil {
		log.Println("Error converting entity id from string to int: ", err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	entity, err := p.entityService.RetrieveEntity(c.Context(), entityID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	invoicePrint := models.NewInvoicePrintFromForm(c)
	invoicePrint.Entity = entity
	invoicePrint.CreatedBy = userID

	if !invoicePrint.IsValid() {
		pageData.InvoicePrint = invoicePrint
		pageData.FormMsg = "Corrija os campos abaixo."
		return utils.RetargetToForm(c, "invoice-print", pageData)
	}

	err = p.service.CreateInvoicePrinting(c.Context(), invoicePrint)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	go p.siareBGWorker.RequestInvoicePrinting(invoicePrint)

	filters := models.NewRawFiltersFromForm(c)

	printings, err := p.service.ListInvoicePrintings(c.Context(), userID, filters)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Trigger-After-Settle", "invoice-print-required")

	shouldWarnUser := utils.FiltersExcludeToday(filters)
	if shouldWarnUser {
		return utils.GeneralInfoResponse(c, globals.ReqCardNotVisibleMsg)
	}
	return c.Render("partials/requests-overview", printings)
}

func (p *PrintInvoicesPage) GetRequestCardDetails(c *fiber.Ctx) error {
	printingId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.PrintingNotFoundErr)
	}
	userID := c.Locals("UserID").(int)
	printing, err := p.service.RetrieveInvoicePrinting(c.Context(), printingId, userID)

	c.Set("HX-Trigger-After-Settle", "open-request-card-details")
	return c.Render("partials/request-card-details", printing)
}

func (p *PrintInvoicesPage) GetInvoicePrintForm(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()
	userID := c.Locals("UserID").(int)

	entities, err := p.entityService.ListEntities(c.Context(), userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	pageData.FormSelectFields.Entities = entities

	printingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.PrintingNotFoundErr)
	}
	printing, err := p.service.RetrieveInvoicePrinting(c.Context(), printingID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	pageData.InvoicePrint = printing

	c.Set("HX-Trigger-After-Settle", "scroll-to-top")
	return c.Render("partials/invoice-print-form", pageData)
}

func (p *PrintInvoicesPage) GetRequestStatus(c *fiber.Ctx) error {
	printingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.PrintingNotFoundErr)
	}

	key := fmt.Sprintf("reqstatus:printing:%v", printingID)
	err = db.Redis.GetDel(c.Context(), key).Err()
	if err == redis.Nil {
		return c.Render("partials/request-card-status", "pending")
	}
	if err != nil {
		log.Printf("Error reading redis key %v: %v\n", key, err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	userID := c.Locals("UserID").(int)

	printing, err := p.service.RetrieveInvoicePrinting(c.Context(), printingID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	targetID := fmt.Sprintf("#request-card-%v", c.Params("id"))
	c.Set("HX-Retarget", targetID)
	c.Set("HX-Reswap", "outerHTML")
	return c.Render("partials/request-card", printing)
}

func (p *PrintInvoicesPage) FilterRequests(c *fiber.Ctx) error {
	userID := c.Locals("UserID").(int)
	filters := c.Queries()

	printings, err := p.service.ListInvoicePrintings(c.Context(), userID, filters)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	return c.Render("partials/requests-overview", printings)
}
