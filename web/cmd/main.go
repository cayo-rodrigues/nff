package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

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

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/static/styles/{stylesheet}", handlers.ServeStyles)
	r.Get("/static/scripts/{script}", handlers.ServeJS)

	r.Get("/", handlers.Index)

	entitiesPage := handlers.EntitiesPage{}
	entitiesPage.ParseTemplates()

	r.Get("/entities", entitiesPage.Render)
	r.Post("/entities", entitiesPage.CreateEntity)
	r.Put("/entities/{id}", entitiesPage.UpdateEntity)
	r.Delete("/entities/{id}", entitiesPage.DeleteEntity)
	r.Get("/entities/{id}/form", entitiesPage.GetEntityForm)

	fmt.Println("Server running on port ", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, r))
}
