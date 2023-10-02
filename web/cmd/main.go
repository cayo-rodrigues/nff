package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/handlers"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
)

func main() {
	PORT, isThere := os.LookupEnv("PORT")
	if !isThere || PORT == "" {
		log.Fatal("PORT env not set or has an empty value")
	}

	DEBUG := false
	_, isThere = os.LookupEnv("DEBUG")
	if isThere {
		DEBUG = true
	}

	dbpool := sql.GetDatabasePool()
	defer dbpool.Close()

	err := dbpool.Ping(context.Background())
	if err != nil {
		log.Fatal("Database connection is not OK, ping failed: ", err)
	}

	engine := html.New("internal/views", ".html")

	engine.Reload(DEBUG)
	engine.Debug(DEBUG)

	engine.AddFunc("GetInvoiceItemSelectFields", func() *models.InvoiceItemFormSelectFields {
		return &models.InvoiceItemFormSelectFields{
			Groups:               &globals.InvoiceItemGroups,
			Origins:              &globals.InvoiceItemOrigins,
			UnitiesOfMeasurement: &globals.InvoiceItemUnitiesOfMeaasurement,
		}
	})

	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
		Prefork:           true,
		AppName:           "NFF",
	})

	app.Use(logger.New())

	app.Get("/static/styles/:stylesheet", handlers.ServeStyles)
	app.Get("/static/scripts/:script", handlers.ServeJS)
	app.Get("/static/icons/:icon", handlers.ServeIcons)

	app.Get("/", handlers.Index)

	entitiesPage := new(handlers.EntitiesPage)
	app.Get("/entities", entitiesPage.Render)
	app.Get("/entities/:id/form", entitiesPage.GetEntityForm)
	app.Post("/entities", entitiesPage.CreateEntity)
	app.Put("/entities/:id", entitiesPage.UpdateEntity)
	app.Delete("/entities/:id", entitiesPage.DeleteEntity)

	invoicesPage := new(handlers.InvoicesPage)
	app.Get("/invoices", invoicesPage.Render)
	app.Post("/invoices", invoicesPage.RequireInvoice)
	app.Get("/invoices/:id/form", invoicesPage.GetInvoiceForm)
	app.Get("/invoices/:id/request-card-details", invoicesPage.GetRequestCardDetails)
	app.Get("/invoices/items/form-section", invoicesPage.GetItemFormSection)

	cancelInvoicesPage := new(handlers.CancelInvoicesPage)
	app.Get("/invoices/cancel", cancelInvoicesPage.Render)
	app.Post("/invoices/cancel", cancelInvoicesPage.CancelInvoice)
	app.Get("/invoices/cancel/:id/form", cancelInvoicesPage.GetInvoiceCancelForm)
	app.Get("/invoices/cancel/:id/request-card-details", cancelInvoicesPage.GetRequestCardDetails)

	log.Fatal(app.Listen(":" + PORT))
}
