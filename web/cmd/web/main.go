package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"

	"github.com/cayo-rodrigues/nff/web/bg-workers"
	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/handlers"
	"github.com/cayo-rodrigues/nff/web/middlewares"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/utils"
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

	DEBUG := os.Getenv("DEBUG") == "true"
	PREFORK := os.Getenv("PREFORK") == "true"
	START_FRESH := os.Getenv("START_FRESH") == "true"

	dbpool := db.GetDBPool()
	defer dbpool.Close()

	rdb := db.GetRedisConn()
	defer rdb.Close()

	store := db.GetSessionStore()
	defer store.Storage.Close()

    if START_FRESH {
        utils.PurgeAllCachedData(context.Background())
    }

	engine := html.New("views", ".html")

	engine.Reload(DEBUG)
	engine.Debug(DEBUG)

	engine.AddFunc("GetInvoiceItemSelectFields", models.NewInvoiceItemSelectFields)
	engine.AddFunc("GetReqCardErrSummary", utils.GetReqCardErrSummary)
	engine.AddFunc("GetReqCardData", models.NewRequestCard)
	engine.AddFunc("FormatDate", utils.FormatDate)
	engine.AddFunc("FormatDateAsBR", utils.FormatDateAsBR)

	userService := services.NewUserService()
	entityService := services.NewEntityService()
	itemsService := services.NewItemsService()
	resultsService := services.NewMetricsResultService()
	filtersService := services.NewFiltersService()
	invoiceService := services.NewInvoiceService(itemsService, filtersService)
	cancelingService := services.NewCancelingService(filtersService)
	printingService := services.NewPrintingService(filtersService)
	metricsService := services.NewMetricsService(resultsService, filtersService)

	siareBGWorker := bgworkers.NewSiareBGWorker(invoiceService, cancelingService, printingService, metricsService, resultsService, SS_API_BASE_URL)

	registerPage := handlers.NewRegisterPage(userService)
	loginPage := handlers.NewLoginPage(userService)
	entitiesPage := handlers.NewEntitiesPage(entityService)
	invoicesPage := handlers.NewInvoicesPage(invoiceService, entityService, siareBGWorker)
	cancelInvoicesPage := handlers.NewCancelInvoicesPage(cancelingService, entityService, siareBGWorker)
	metricsPage := handlers.NewMetricsPage(metricsService, entityService, siareBGWorker)
	printInvoicesPage := handlers.NewPrintInvoicesPage(printingService, entityService, siareBGWorker)

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

	app.Get("/register", registerPage.Render)
	app.Post("/register", registerPage.CreateUser)

	app.Get("/login", loginPage.Render)
	app.Post("/login", loginPage.Login)
	app.Get("/logout", handlers.Logout)

	app.Use(middlewares.AuthMiddleware)
	app.Use(middlewares.CacheMiddleware)

	app.Get("/", handlers.Home)

	app.Get("/entities", entitiesPage.Render)
	app.Post("/entities", entitiesPage.CreateEntity)
	app.Get("/entities/:id/form", entitiesPage.GetEntityForm)
	app.Put("/entities/:id", entitiesPage.UpdateEntity)
	app.Delete("/entities/:id", entitiesPage.DeleteEntity)

	app.Get("/invoices", invoicesPage.Render)
	app.Post("/invoices", invoicesPage.RequireInvoice)
	app.Get("/invoices/:id/form", invoicesPage.GetInvoiceForm)
	app.Get("/invoices/:id/request-card-details", invoicesPage.GetRequestCardDetails)
	app.Get("/invoices/:id/request-card-status", invoicesPage.GetRequestStatus)
	app.Get("/invoices/items/form-section", invoicesPage.GetItemFormSection)
	app.Get("/invoices/request-card-filter", invoicesPage.FilterRequests)

	app.Get("/invoices/cancel", cancelInvoicesPage.Render)
	app.Post("/invoices/cancel", cancelInvoicesPage.CancelInvoice)
	app.Get("/invoices/cancel/:id/form", cancelInvoicesPage.GetInvoiceCancelForm)
	app.Get("/invoices/cancel/:id/request-card-details", cancelInvoicesPage.GetRequestCardDetails)
	app.Get("/invoices/cancel/:id/request-card-status", cancelInvoicesPage.GetRequestStatus)
	app.Get("/invoices/cancel/request-card-filter", cancelInvoicesPage.FilterRequests)

	app.Get("/invoices/print", printInvoicesPage.Render)
	app.Post("/invoices/print", printInvoicesPage.PrintInvoice)
	app.Get("/invoices/print/:id/form", printInvoicesPage.GetInvoicePrintForm)
	app.Get("/invoices/print/:id/request-card-details", printInvoicesPage.GetRequestCardDetails)
	app.Get("/invoices/print/:id/request-card-status", printInvoicesPage.GetRequestStatus)
	app.Get("/invoices/print/request-card-filter", printInvoicesPage.FilterRequests)

	app.Get("/metrics", metricsPage.Render)
	app.Post("/metrics", metricsPage.GenerateMetrics)
	app.Get("/metrics/:id/form", metricsPage.GetMetricsForm)
	app.Get("/metrics/:id/request-card-details", metricsPage.GetRequestCardDetails)
	app.Get("/metrics/:id/request-card-status", metricsPage.GetRequestStatus)
	app.Get("/metrics/request-card-filter", metricsPage.FilterRequests)

	err := app.Listen(":" + PORT)
	if err != nil {
		log.Fatalln(">:(", err)
	}
}
