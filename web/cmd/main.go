package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	dbpool := sql.GetDatabasePool()
	defer dbpool.Close()

	err := dbpool.Ping(context.Background())
	if err != nil {
		log.Fatal("Database connection is not OK, ping failed: ", err)
	}

	engine := html.New("internal/views", ".html")

	// Reload the templates on each render, good for development
	engine.Reload(true) // Optional. Default: false

	// Debug will print each template that is parsed, good for debugging
	engine.Debug(true) // Optional. Default: false

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
	})

	app.Use(cors.New())

	app.Get("/static/styles/:stylesheet", handlers.ServeStyles)
	app.Get("/static/scripts/:script", handlers.ServeJS)

	app.Get("/", handlers.Index)

	entitiesPage := &handlers.EntitiesPage{}
	app.Get("/entities", entitiesPage.Render)
	app.Get("/entities/:id/form", entitiesPage.GetEntityForm)
	app.Post("/entities", entitiesPage.CreateEntity)
	app.Put("/entities/:id", entitiesPage.UpdateEntity)
	app.Delete("/entities/:id", entitiesPage.DeleteEntity)

	invoicesPage := &handlers.InvoicesPage{}
	app.Get("/invoices", invoicesPage.Render)
	app.Post("/invoices", invoicesPage.RequireInvoice)
	app.Get("/invoices/items/form-section", invoicesPage.GetItemFormSection)

	cancelInvoicesPage := &handlers.CancelInvoicesPage{}
	app.Get("/invoices/cancel", cancelInvoicesPage.Render)
	app.Post("/invoices/cancel", cancelInvoicesPage.CancelInvoice)

	fmt.Println("Server running on port", PORT)
	log.Fatal(app.Listen(":" + PORT))

}