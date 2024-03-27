package main

import (
	"context"
	"log"
	"os"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("PORT env not set or has an empty value")
	}

	SS_API_BASE_URL := os.Getenv("SS_API_BASE_URL")
	if SS_API_BASE_URL == "" {
		log.Fatal("SS_API_BASE_URL env not set or has an empty value")
	}

	PREFORK := os.Getenv("PREFORK") == "true"
	START_FRESH := os.Getenv("START_FRESH") == "true"

	db := database.GetDB()
	defer db.Close()

	if START_FRESH {
		utils.PurgeAllCachedData(context.Background())
	}

	app := fiber.New(fiber.Config{
		Prefork: PREFORK,
		AppName: "NFF",
	})

	app.Use(logger.New())

	app.Static("/static", "./static")

	// app.Get("/register", registerPage.Render)
	// app.Post("/register", registerPage.CreateUser)

	// app.Get("/login", loginPage.Render)
	// app.Post("/login", loginPage.Login)
	// app.Get("/logout", handlers.Logout)

	// app.Use(middlewares.AuthMiddleware)
	// app.Use(middlewares.CacheMiddleware)

	// app.Get("/", handlers.Home)

	// app.Get("/entities", entitiesPage.Render)
	// app.Post("/entities", entitiesPage.CreateEntity)
	// app.Get("/entities/:id/form", entitiesPage.GetEntityForm)
	// app.Put("/entities/:id", entitiesPage.UpdateEntity)
	// app.Delete("/entities/:id", entitiesPage.DeleteEntity)

	// app.Get("/invoices", invoicesPage.Render)
	// app.Post("/invoices", invoicesPage.RequireInvoice)
	// app.Get("/invoices/:id/form", invoicesPage.GetInvoiceForm)
	// app.Get("/invoices/:id/request-card-details", invoicesPage.GetRequestCardDetails)
	// app.Get("/invoices/:id/request-card-status", invoicesPage.GetRequestStatus)
	// app.Get("/invoices/items/form-section", invoicesPage.GetItemFormSection)
	// app.Get("/invoices/request-card-filter", invoicesPage.FilterRequests)

	// app.Get("/invoices/cancel", cancelInvoicesPage.Render)
	// app.Post("/invoices/cancel", cancelInvoicesPage.CancelInvoice)
	// app.Get("/invoices/cancel/:id/form", cancelInvoicesPage.GetInvoiceCancelForm)
	// app.Get("/invoices/cancel/:id/request-card-details", cancelInvoicesPage.GetRequestCardDetails)
	// app.Get("/invoices/cancel/:id/request-card-status", cancelInvoicesPage.GetRequestStatus)
	// app.Get("/invoices/cancel/request-card-filter", cancelInvoicesPage.FilterRequests)

	// app.Get("/invoices/print", printInvoicesPage.Render)
	// app.Post("/invoices/print", printInvoicesPage.PrintInvoice)
	// app.Get("/invoices/print/:id/form", printInvoicesPage.GetInvoicePrintForm)
	// app.Get("/invoices/print/:id/request-card-details", printInvoicesPage.GetRequestCardDetails)
	// app.Get("/invoices/print/:id/request-card-status", printInvoicesPage.GetRequestStatus)
	// app.Get("/invoices/print/request-card-filter", printInvoicesPage.FilterRequests)

	// app.Get("/metrics", metricsPage.Render)
	// app.Post("/metrics", metricsPage.GenerateMetrics)
	// app.Get("/metrics/:id/form", metricsPage.GetMetricsForm)
	// app.Get("/metrics/:id/request-card-details", metricsPage.GetRequestCardDetails)
	// app.Get("/metrics/:id/request-card-status", metricsPage.GetRequestStatus)
	// app.Get("/metrics/request-card-filter", metricsPage.FilterRequests)
	// app.Get("/metrics/results/records/print", metricsPage.PrintInvoice)

	// app.Use(middlewares.NotFoundMiddleware)

	err := app.Listen(":" + PORT)
	if err != nil {
		log.Fatalln(">:(", err)
	}
}
