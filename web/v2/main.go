package main

import (
	"context"
	"log"
	"os"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/handlers"
	"github.com/cayo-rodrigues/nff/web/middlewares"
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

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("Error initializing database connection: ", err)
	}
	defer db.Close()

	if START_FRESH {
		db.Redis.DestroyAllCachedData(context.Background())
	}

	app := fiber.New(fiber.Config{
		Prefork: PREFORK,
		AppName: "NFF",
	})

	app.Use(logger.New())

	app.Static("/static", "./static")

	app.Get("/register", handlers.RegisterPage)
	app.Post("/register", handlers.RegisterUser)

	app.Get("/login", handlers.LoginPage)
	app.Post("/login", handlers.LoginUser)
	app.Get("/logout", handlers.LogoutUser)

	app.Use(middlewares.AuthMiddleware)
	app.Use(middlewares.CacheMiddleware)

	app.Get("/", handlers.HomePage)

	app.Get("/entities", handlers.EntitiesPage)
	app.Get("/entities/search", handlers.SearchEntities)
	app.Get("/entities/create", handlers.CreateEntityPage)
	app.Post("/entities/create", handlers.CreateEntity)
	app.Get("/entities/update/:id", handlers.EditEntityPage)
	app.Put("/entities/update/:id", handlers.UpdateEntity)
	app.Delete("/entities/delete/:id", handlers.DeleteEntity)

	app.Get("/invoices", handlers.InvoicesPage)
	app.Get("/invoices/create", handlers.CreateInvoicePage)
	app.Get("/invoices/form/get-sender-ie-input", handlers.GetSenderIeInput)
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

	app.Get("/metrics", handlers.MetricsPage)
	// app.Post("/metrics", metricsPage.GenerateMetrics)
	// app.Get("/metrics/:id/form", metricsPage.GetMetricsForm)
	// app.Get("/metrics/:id/request-card-details", metricsPage.GetRequestCardDetails)
	// app.Get("/metrics/:id/request-card-status", metricsPage.GetRequestStatus)
	// app.Get("/metrics/request-card-filter", metricsPage.FilterRequests)
	// app.Get("/metrics/results/records/print", metricsPage.PrintInvoice)

	app.Use(handlers.NotFoundPage)

	err = app.Listen(":" + PORT)
	if err != nil {
		log.Fatalln(">:(", err)
	}
}
