package handlers

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/components"
	"github.com/cayo-rodrigues/nff/web/ui/forms"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func PrintInvoicePage(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	printingsList, err := services.ListPrintings(c.Context(), userID)
	if err != nil {
		return err
	}

	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}
	printingForForm := models.NewInvoicePrintWithSamples(entities)

	return Render(c, layouts.Base(pages.InvoicesPrintPage(printingsList, printingForForm, entities)))
}

func PrintInvoice(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)

	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	printing := models.NewInvoicePrintFromForm(c)

	entity, err := services.RetrieveEntity(c.Context(), printing.Entity.ID, userID)
	if err != nil {
		return err
	}
	printing.Entity = entity

	if printing.IsValid() {
		err := services.CreatePrinting(c.Context(), printing, userID)
		if err != nil {
			return err
		}
		c.Set("HX-Trigger-After-Swap", "reload-printing-list")
	}

	return Render(c, forms.PrintInvoiceForm(printing, entities))
}

func ListInvoicePrintings(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)
	filters := c.Queries()
	printings, err := services.ListPrintings(c.Context(), userID, filters)
	if err != nil {
		return err
	}

	return Render(c, components.InvoicesPrintingsList(printings))
}

func GetPrintInvoiceForm(c *fiber.Ctx) error {
	userID := utils.GetCurrentUserID(c)

	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}

	basePrintingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}

	basePrinting, err := services.RetrievePrinting(c.Context(), basePrintingID, userID)
	if err != nil {
		return err
	}

	c.Set("HX-Trigger-After-Swap", "scroll-to-top")
	return Render(c, forms.PrintInvoiceForm(basePrinting, entities))
}
