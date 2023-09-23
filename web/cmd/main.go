package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"

	"github.com/cayo-rodrigues/nff/web/internal/handlers"
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

	// AddFunc adds a function to the template's global function map.
	engine.AddFunc("greet", func(name string) string {
		return "Hello, " + name + "!"
	})

	app := fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	app.Use(cors.New())

	// app.Static("/static", "internal/static")
	app.Get("/static/styles/:stylesheet", handlers.ServeStyles)
	app.Get("/static/scripts/:script", handlers.ServeJS)

	app.Get("/", handlers.Index)
	app.Get("/p", func(c *fiber.Ctx) error {
		return c.Render("partials/p", fiber.Map{
			"Words": []string{"aaaa", "bbb", "c"},
		})
	})

	fmt.Println("Server running on port", PORT)
	log.Fatal(app.Listen(":" + PORT))

	r := chi.NewRouter()

	entitiesPage := handlers.NewEntitiesPage()

	r.Get("/entities", entitiesPage.Render)
	r.Post("/entities", entitiesPage.CreateEntity)
	r.Put("/entities/{id}", entitiesPage.UpdateEntity)
	r.Delete("/entities/{id}", entitiesPage.DeleteEntity)
	r.Get("/entities/{id}/form", entitiesPage.GetEntityForm)

	invoicesPage := handlers.NewInvoicesPage()

	r.Get("/invoices", invoicesPage.Render)
	r.Post("/invoices", invoicesPage.RequireInvoice)
	r.Get("/invoices/items/form-section", invoicesPage.GetItemFormSection)
}
