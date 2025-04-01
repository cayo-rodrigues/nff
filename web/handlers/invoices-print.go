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
	filters := c.Queries()

	printingsList, err := services.ListPrintings(c.Context(), filters)
	if err != nil {
		return err
	}

	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}

	entitiesByType := models.NewEntitiesByType(entities)

	printingForForm := models.NewInvoicePrintWithSamples(entitiesByType.Senders)
	printingsByDate := services.GroupListByDate(printingsList)
	page := pages.InvoicesPrintPage(printingsByDate, printingForForm, entitiesByType)

	c.Append("HX-Trigger-After-Settle", "highlight-current-filter", "highlight-current-page", "notification-list-loaded")
	return Render(c, layouts.Base(page))
}

func PrintInvoice(c *fiber.Ctx) error {
	decryptionKey, err := services.GetEncryptionKeySession(c)
	if err != nil {
		return RetargetToReauth(c)
	}

	entities, err := services.ListEntities(c.Context())
	if err != nil {
		return err
	}
	entitiesByType := models.NewEntitiesByType(entities)

	printing := models.NewInvoicePrintFromForm(c)

	entity, err := services.RetrieveEntity(c.Context(), printing.Entity.ID)
	if err != nil {
		return err
	}
	printing.Entity = entity


	if !printing.IsValid() {
		return Render(c, forms.PrintInvoiceForm(printing, entitiesByType.Senders))
	}

	err = services.CreatePrinting(c.Context(), printing)
	if err != nil {
		return err
	}

	ssapi := siare.GetSSApiClient().WithDecryptionKey(decryptionKey)
	go ssapi.PrintInvoice(printing)

	c.Append("HX-Trigger-After-Swap", "reload-printing-list")
	return Render(c, forms.PrintInvoiceForm(printing, entitiesByType.Senders))
}

func PrintInvoiceFromMetricsRecord(c *fiber.Ctx) error {
	decryptionKey, err := services.GetEncryptionKeySession(c)
	if err != nil {
		return RetargetToReauth(c)
	}

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

	ssapi := siare.GetSSApiClient().WithDecryptionKey(decryptionKey)
	go ssapi.PrintInvoiceFromMetricsRecord(printing, recordID, userID)

	record := models.NewMetricsResult()
	record.ID = recordID

	return Render(c, components.DownloadInvoiceFromRecordLoadingIcon(record))
}

func ListInvoicePrintings(c *fiber.Ctx) error {
	filters := c.Queries()
	printings, err := services.ListPrintings(c.Context(), filters)
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
	entitiesByType := models.NewEntitiesByType(entities)

	basePrintingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	basePrinting, err := services.RetrievePrinting(c.Context(), basePrintingID)
	if err != nil {
		return err
	}

	c.Append("HX-Trigger-After-Swap", "scroll-to-top")
	return Render(c, forms.PrintInvoiceForm(basePrinting, entitiesByType.Senders))
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
