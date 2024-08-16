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
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func PrintInvoicePage(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID

	filters := c.Queries()

	printingsList, err := services.ListPrintings(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}

	printingForForm := models.NewInvoicePrintWithSamples(entities)
	printingsByDate := services.GroupListByDate(printingsList)
	page := pages.InvoicesPrintPage(printingsByDate, printingForForm, entities)

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter", "highlight-current-page", "notification-list-loaded")
	return Render(c, layouts.Base(page))
}

func PrintInvoice(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID

	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}

	printing := models.NewInvoicePrintFromForm(c)

	entity, err := services.RetrieveEntity(c.Context(), printing.Entity.ID)
	if err != nil {
		return err
	}
	printing.Entity = entity

	if !printing.IsValid() {
		return Render(c, forms.PrintInvoiceForm(printing, entities))
	}

	err = services.CreatePrinting(c.Context(), printing, userID)
	if err != nil {
		return err
	}

	ssapi := siare.GetSSApiClient()
	go ssapi.PrintInvoice(printing)

	c.Append("HX-Trigger-After-Swap", "reload-printing-list")
	return Render(c, forms.PrintInvoiceForm(printing, entities))
}

func PrintInvoiceFromMetricsRecord(c *fiber.Ctx) error {
	userID := utils.GetUserID(c.Context())

	recordID, err := c.ParamsInt("record_id")
	if err != nil {
		return err
	}
	invoiceNumber := c.Params("invoice_number")
	entityID, err := c.ParamsInt("entity_id")
	if err != nil {
		return err
	}

	printing, err := services.CreatePrintingFromMetricsRecord(c.Context(), invoiceNumber, entityID)

	ssapi := siare.GetSSApiClient()
	go ssapi.PrintInvoiceFromMetricsRecord(printing, recordID, userID)

	record := models.NewMetricsResult()
	record.ID = recordID

	return Render(c, components.DownloadInvoiceFromRecordLoadingIcon(record))
}

func ListInvoicePrintings(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	filters := c.Queries()
	printings, err := services.ListPrintings(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	printingsByDate := services.GroupListByDate(printings)

	return Render(c, components.InvoicesPrintingsList(printingsByDate))
}

func GetPrintInvoiceForm(c *fiber.Ctx) error {
	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}

	basePrintingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	basePrinting, err := services.RetrievePrinting(c.Context(), basePrintingID)
	if err != nil {
		return err
	}

	c.Append("HX-Trigger-After-Swap", "scroll-to-top")
	return Render(c, forms.PrintInvoiceForm(basePrinting, entities))
}

func RetrieveInvoicePrintCard(c *fiber.Ctx) error {
	printingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	printing, err := services.RetrievePrinting(c.Context(), printingID)
	if err != nil {
		return err
	}

	return Render(c, components.InvoicePrintCard(printing))
}
