package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"

	"github.com/cayo-rodrigues/nff/web/internal/bg-workers"
	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/handlers"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/services"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
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

	PREFORK := false
	_, isThere = os.LookupEnv("PREFORK")
	if isThere {
		PREFORK = true
	}

	bgworkers.SS_API_BASE_URL, isThere = os.LookupEnv("SS_API_BASE_URL")
	if !isThere || bgworkers.SS_API_BASE_URL == "" {
		log.Fatal("SS_API_BASE_URL env not set or has an empty value")
	}

	dbpool := db.GetDBPool()
	defer dbpool.Close()

	rdb := db.GetRedisConn()
	defer rdb.Close()

	engine := html.New("internal/views", ".html")

	engine.Reload(DEBUG)
	engine.Debug(DEBUG)

	engine.AddFunc("GetInvoiceItemSelectFields", utils.GetInvoiceItemSelectFields)
	engine.AddFunc("GetReqCardErrSummary", utils.GetReqCardErrSummary)
	engine.AddFunc("GetReqCardData", models.NewRequestCard)
	engine.AddFunc("FormatDate", utils.FormatDate)
	engine.AddFunc("FormatDateAsBR", utils.FormatDateAsBR)

	entityService := services.NewEntityService()
	itemsService := services.NewItemsService()
	invoiceService := services.NewInvoiceService(entityService, itemsService)
	cancelingService := services.NewCancelingService(entityService)
	metricsService := services.NewMetricsService(entityService)

	siareBGWorker := bgworkers.NewSiareBGWorker(invoiceService, cancelingService, metricsService)

	entitiesPage := handlers.NewEntitiesPage(entityService)
	invoicesPage := handlers.NewInvoicesPage(invoiceService, entityService, siareBGWorker)
	cancelInvoicesPage := handlers.NewCancelInvoicesPage(cancelingService, entityService, siareBGWorker)
	metricsPage := handlers.NewMetricsPage(metricsService, entityService, siareBGWorker)

	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
		Prefork:           PREFORK,
		AppName:           "NFF",
	})

	app.Use(logger.New())

	app.Get("/static/styles/:stylesheet", handlers.ServeStyles)
	app.Get("/static/scripts/:script", handlers.ServeJS)
	app.Get("/static/icons/:icon", handlers.ServeIcons)

	app.Get("/", handlers.Index)

	app.Get("/entities", entitiesPage.Render)
	app.Get("/entities/:id/form", entitiesPage.GetEntityForm)
	app.Post("/entities", entitiesPage.CreateEntity)
	app.Put("/entities/:id", entitiesPage.UpdateEntity)
	app.Delete("/entities/:id", entitiesPage.DeleteEntity)

	app.Get("/invoices", invoicesPage.Render)
	app.Post("/invoices", invoicesPage.RequireInvoice)
	app.Get("/invoices/:id/form", invoicesPage.GetInvoiceForm)
	app.Get("/invoices/:id/request-card-details", invoicesPage.GetRequestCardDetails)
	app.Get("/invoices/:id/request-card-status", invoicesPage.GetRequestStatus)
	app.Get("/invoices/items/form-section", invoicesPage.GetItemFormSection)

	app.Get("/invoices/cancel", cancelInvoicesPage.Render)
	app.Post("/invoices/cancel", cancelInvoicesPage.CancelInvoice)
	app.Get("/invoices/cancel/:id/form", cancelInvoicesPage.GetInvoiceCancelForm)
	app.Get("/invoices/cancel/:id/request-card-details", cancelInvoicesPage.GetRequestCardDetails)
	app.Get("/invoices/cancel/:id/request-card-status", cancelInvoicesPage.GetRequestStatus)

	app.Get("/metrics", metricsPage.Render)
	app.Post("/metrics", metricsPage.GenerateMetrics)
	app.Get("/metrics/:id/form", metricsPage.GetMetricsForm)
	app.Get("/metrics/:id/request-card-details", metricsPage.GetRequestCardDetails)
	app.Get("/metrics/:id/request-card-status", metricsPage.GetRequestStatus)

	err := app.Listen(":" + PORT)
	if err != nil {
		log.Fatalln(">:(", err)
	}
}
