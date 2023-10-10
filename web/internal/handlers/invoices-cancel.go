package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/cayo-rodrigues/nff/web/internal/workers"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
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
		data.InvoiceCancel = invoiceCancel
		data.FormMsg = "Corrija os campos abaixo."
		c.Set("HX-Retarget", "#invoice-cancel-form")
		c.Set("HX-Reswap", "outerHTML")
		return c.Render("partials/invoice-cancel-form", data)
	}

	err = workers.CreateInvoiceCanceling(c.Context(), invoiceCancel)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	go func(invoiceCancel *models.InvoiceCancel) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel()

		// do the thing
		time.Sleep(time.Second * 15)

		invoiceCancel.ReqStatus = "success"
		invoiceCancel.ReqMsg = "Cancelamento efetuado com sucesso!"
		err := workers.UpdateInvoiceCanceling(ctx, invoiceCancel)
		if err != nil {
			log.Println("ops")
		}

		key := fmt.Sprintf("reqstatus:canceling:%v", invoiceCancel.Id)
		db.Redis.Set(ctx, key, true, time.Minute)
	}(invoiceCancel)

	c.Set("HX-Trigger-After-Settle", "invoice-cancel-required")
	return c.Render("partials/request-card", invoiceCancel)
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

func (page *CancelInvoicesPage) GetRequestStatus(c *fiber.Ctx) error {
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

	canceling, err := workers.RetrieveInvoiceCanceling(c.Context(), cancelingId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Retarget", "#request-card-" + c.Params("id"))
	c.Set("HX-Reswap", "outerHTML")
	return c.Status(286).Render("partials/request-card", canceling)
}
