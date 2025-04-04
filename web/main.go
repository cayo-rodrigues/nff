package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/handlers"

	"github.com/cayo-rodrigues/nff/web/handlers/sse"
	"github.com/cayo-rodrigues/nff/web/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

//go:embed static/*
var staticFiles embed.FS

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
		// ErrorHandler: SERIA UTIL?
	})

	app.Use(logger.New(logger.Config{
		TimeZone: "America/Sao_Paulo",
	}))

	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(staticFiles),
		PathPrefix: "static",
		Browse:     true,
	}))

	app.Get("/register", handlers.RegisterPage)
	app.Post("/register", handlers.RegisterUser)

	app.Get("/login", handlers.LoginPage)
	app.Post("/login", handlers.LoginUser)

	app.Use(middlewares.AuthMiddleware)

	app.Get("/logout", handlers.LogoutUser)

	app.Get("/sse/notify-operations-results", sse.NotifyOperationsResults)

	app.Get("/notifications", handlers.ListNotifications)
	app.Get("/notifications/latest", handlers.GetLatestNotification)
	app.Delete("/notifications", handlers.ClearNotifications)

	app.Use(middlewares.CacheManagementMiddleware)

	app.Post("/reauthenticate", handlers.ReauthUser)

	app.Get("/", handlers.HomePage)

	app.Get("/entities", handlers.EntitiesPage)
	app.Get("/entities/search", handlers.SearchEntities)
	app.Get("/entities/create", handlers.CreateEntityPage)
	app.Post("/entities/create", handlers.CreateEntity)
	app.Get("/entities/update/:id", handlers.EditEntityPage)
	app.Put("/entities/update/:id", handlers.UpdateEntity)
	app.Delete("/entities/delete/:id", handlers.DeleteEntity)

	app.Get("/invoices", handlers.InvoicesPage)
	app.Post("/invoices", handlers.CreateInvoice)
	app.Get("/invoices/list", handlers.ListInvoices)
	app.Get("/invoices/form/get-sender-ie-input", handlers.GetSenderIeInput)
	app.Get("/invoices/form/get-recipient-ie-input", handlers.GetRecipientIeInput)
	app.Get("/invoices/form/get-cfops-input", handlers.GetCfopsInput)
	app.Get("/invoices/:id/form", handlers.GetInvoiceForm)
	app.Get("/invoices/:id/items-details", handlers.RetrieveInvoiceItemsDetails)
	app.Get("/invoices/:id/card", handlers.RetrieveInvoiceCard)

	app.Get("/invoices/cancel", handlers.CancelInvoicePage)
	app.Post("/invoices/cancel", handlers.CancelInvoice)
	app.Post("/invoices/cancel/:invoice_id", handlers.CancelInvoiceByID)
	app.Get("/invoices/cancel/list", handlers.ListInvoiceCancelings)
	app.Get("/invoices/cancel/:id/form", handlers.GetCancelInvoiceForm)
	app.Get("/invoices/cancel/:id/card", handlers.RetrieveInvoiceCancelCard)

	app.Get("/invoices/print", handlers.PrintInvoicePage)
	app.Post("/invoices/print", handlers.PrintInvoice)
	app.Post("/invoices/print/:record_id/:invoice_number/:entity_id", handlers.PrintInvoiceFromMetricsRecord)
	app.Get("/invoices/print/list", handlers.ListInvoicePrintings)
	app.Get("/invoices/print/:id/form", handlers.GetPrintInvoiceForm)
	app.Get("/invoices/print/:id/card", handlers.RetrieveInvoicePrintCard)

	app.Get("/metrics", handlers.MetricsPage)
	app.Post("/metrics", handlers.GenerateMetrics)
	app.Get("/metrics/list", handlers.ListMetrics)
	app.Get("/metrics/:id", handlers.MetricsDetailsPage)
	app.Get("/metrics/:id/form", handlers.GetMetricsForm)
	app.Get("/metrics/:id/results-details", handlers.RetrieveMetricsResultsDetails)
	app.Get("/metrics/:id/card", handlers.RetrieveMetricsCard)
	app.Get("/metrics/:id/download-from-record-status-icon", handlers.GetDownloadFromRecordStatusIcon)

	app.Use(handlers.NotFoundPage)

	err = app.Listen(":" + PORT)
	if err != nil {
		log.Fatalln(">:(", err)
	}
}
