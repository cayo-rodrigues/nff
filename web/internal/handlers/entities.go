package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/services"
	"github.com/go-chi/chi/v5"
)

type EntitiesPage struct {
	tmpl *template.Template
}

type EntitiesPageData struct {
	IsAuthenticated bool
	Entities        *[]models.Entity
	Entity          *models.Entity
}

func (page *EntitiesPage) ParseTemplates() {
	page.tmpl = template.Must(template.ParseFiles(
		"internal/templates/layout.html",
		"internal/templates/entities.html",
		"internal/templates/entity-form.html",
	))
}

func (page *EntitiesPage) Render(w http.ResponseWriter, r *http.Request) {
	entities, err := services.ListEntities(r.Context())
	if err != nil {
		return
	}

	w.Header().Add("HX-Trigger-After-Settle", "entity-page-loaded")

	page.tmpl.ExecuteTemplate(w, "layout", &EntitiesPageData{
		IsAuthenticated: true,
		Entities:        entities,
		Entity: &models.Entity{
			Address: &models.Address{},
		},
	})
}

func (page *EntitiesPage) GetEntityForm(w http.ResponseWriter, r *http.Request) {
	entityId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return
	}
	entity, err := services.RetrieveEntity(r.Context(), entityId)
	if err != nil {
		return
	}

	entity.IsSelected = true
	page.tmpl.ExecuteTemplate(w, "entity-form", &EntitiesPageData{
		Entity: entity,
	})
}

func (page *EntitiesPage) CreateEntity(w http.ResponseWriter, r *http.Request) {
	// TODO validate data
	entity := models.NewEntityFromForm(r)

	err := services.RegisterEntity(r.Context(), entity)
	if err != nil {
		return
	}

	page.tmpl.ExecuteTemplate(w, "entity-card", entity)
}

func (page *EntitiesPage) UpdateEntity(w http.ResponseWriter, r *http.Request) {
	entity := models.NewEntityFromForm(r)
	entityId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return
	}
	entity.Id = entityId

	err = services.UpdateEntity(r.Context(), entity)
	if err != nil {
		return
	}
	entity.IsSelected = true

	page.tmpl.ExecuteTemplate(w, "entity-card", entity)
}

func (page *EntitiesPage) DeleteEntity(w http.ResponseWriter, r *http.Request) {
	entityId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return
	}
	err = services.DeleteEntity(r.Context(), entityId)
	if err != nil {
		return
	}

	eventMsg := fmt.Sprintf("{\"entity-deleted\": %v}", entityId)
	w.Header().Add("HX-Trigger-After-Settle", eventMsg)

	page.tmpl.ExecuteTemplate(w, "entity-form", nil)
}
