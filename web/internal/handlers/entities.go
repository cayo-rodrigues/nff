package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/cayo-rodrigues/nff/web/internal/workers"
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
	data := &EntitiesPageData{
		IsAuthenticated: true,
		Entities:        nil,
		Entity: &models.Entity{
			Address: &models.Address{},
		},
	}
	entities, err := workers.ListEntities(r.Context())
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl, "layout", data)
		return
	}

	data.Entities = entities
	page.tmpl.ExecuteTemplate(w, "layout", data)
}

func (page *EntitiesPage) GetEntityForm(w http.ResponseWriter, r *http.Request) {
	data := &EntitiesPageData{
		Entity: &models.Entity{
			Address: &models.Address{},
		},
	}

	entityId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl, "entity-form", data)
		return
	}

	entity, err := workers.RetrieveEntity(r.Context(), entityId)
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl, "entity-form", data)
		return
	}

	data.Entity = entity
	data.Entity.IsSelected = true
	page.tmpl.ExecuteTemplate(w, "entity-form", data)
}

func (page *EntitiesPage) CreateEntity(w http.ResponseWriter, r *http.Request) {
	// TODO validate data
	entity := models.NewEntityFromForm(r)

	err := workers.RegisterEntity(r.Context(), entity)
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl, "entity-card", nil)
		return
	}

	page.tmpl.ExecuteTemplate(w, "entity-card", entity)
}

func (page *EntitiesPage) UpdateEntity(w http.ResponseWriter, r *http.Request) {
	entity := models.NewEntityFromForm(r)
	entityId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl, "entity-card", nil)
		return
	}
	entity.Id = entityId

	err = workers.UpdateEntity(r.Context(), entity)
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl, "entity-card", nil)
		return
	}

	entity.IsSelected = true

	page.tmpl.ExecuteTemplate(w, "entity-card", entity)
}

func (page *EntitiesPage) DeleteEntity(w http.ResponseWriter, r *http.Request) {
	entityId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl, "entity-form", nil)
		return
	}
	err = workers.DeleteEntity(r.Context(), entityId)
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl, "entity-form", nil)
	}

	eventMsg := fmt.Sprintf("{\"entity-deleted\": %v}", entityId)
	w.Header().Add("HX-Trigger-After-Settle", eventMsg)

	page.tmpl.ExecuteTemplate(w, "entity-form", nil)
}
