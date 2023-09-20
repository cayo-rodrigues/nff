package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/cayo-rodrigues/nff/web/internal/workers"
)

type EntitiesPage struct {
	tmpl *template.Template
}

type EntitiesPageData struct {
	IsAuthenticated  bool
	Entities         *[]models.Entity
	Entity           *models.Entity
	GeneralError     string
	FormSelectFields *models.EntityFormSelectFields
}

func NewEntitiesPage() *EntitiesPage {
	entitiesPage := &EntitiesPage{}
	entitiesPage.ParseTemplates()
	return entitiesPage
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
			Errors:  &models.EntityFormError{},
		},
		FormSelectFields: &models.EntityFormSelectFields{
			UserTypes:   &globals.EntityUserTypes,
			StreetTypes: &globals.EntityAddressStreetTypes,
		},
	}
	entities, err := workers.ListEntities(r.Context())
	if err != nil {
		data.GeneralError = err.Error()
		utils.ErrorResponse(w, "general-error", page.tmpl, "layout", data)
		return
	}

	data.Entities = entities
	page.tmpl.ExecuteTemplate(w, "layout", data)
}

func (page *EntitiesPage) GetEntityForm(w http.ResponseWriter, r *http.Request) {
	data := &EntitiesPageData{
		Entity: &models.Entity{
			Address: &models.Address{},
			Errors:  &models.EntityFormError{},
		},
		FormSelectFields: &models.EntityFormSelectFields{
			UserTypes:   &globals.EntityUserTypes,
			StreetTypes: &globals.EntityAddressStreetTypes,
		},
	}

	entityId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.GeneralErrorResponse(w, utils.EntityNotFoundErr, page.tmpl)
		return
	}

	entity, err := workers.RetrieveEntity(r.Context(), entityId)
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}

	data.Entity = entity
	data.Entity.IsSelected = true
	page.tmpl.ExecuteTemplate(w, "entity-form", data)
}

func (page *EntitiesPage) CreateEntity(w http.ResponseWriter, r *http.Request) {
	entity := models.NewEntityFromForm(r)

	if !entity.IsValid() {
		w.Header().Add("HX-Retarget", "#entity-form")
		w.Header().Add("HX-Reswap", "outerHTML")
		page.tmpl.ExecuteTemplate(w, "entity-form", &EntitiesPageData{
			Entity: entity,
			FormSelectFields: &models.EntityFormSelectFields{
				UserTypes: &globals.EntityUserTypes,
				StreetTypes: &globals.EntityAddressStreetTypes,
			},
		})
		return
	}

	err := workers.RegisterEntity(r.Context(), entity)
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}

	w.Header().Add("HX-Trigger-After-Settle", "entity-created")
	page.tmpl.ExecuteTemplate(w, "entity-card", entity)
}

func (page *EntitiesPage) UpdateEntity(w http.ResponseWriter, r *http.Request) {
	entity := models.NewEntityFromForm(r)
	entity.IsSelected = true

	entityId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.GeneralErrorResponse(w, utils.EntityNotFoundErr, page.tmpl)
		return
	}
	entity.Id = entityId

	if !entity.IsValid() {
		w.Header().Add("HX-Retarget", "#entity-form")
		w.Header().Add("HX-Reswap", "outerHTML")
		page.tmpl.ExecuteTemplate(w, "entity-form", &EntitiesPageData{
			Entity: entity,
			FormSelectFields: &models.EntityFormSelectFields{
				UserTypes: &globals.EntityUserTypes,
				StreetTypes: &globals.EntityAddressStreetTypes,
			},
		})
		return
	}

	err = workers.UpdateEntity(r.Context(), entity)
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}

	w.Header().Add("HX-Trigger-After-Settle", "entity-updated")
	page.tmpl.ExecuteTemplate(w, "entity-card", entity)
}

func (page *EntitiesPage) DeleteEntity(w http.ResponseWriter, r *http.Request) {
	entityId, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}
	err = workers.DeleteEntity(r.Context(), entityId)
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}

	eventMsg := fmt.Sprintf("{\"entity-deleted\": %v}", entityId)
	w.Header().Add("HX-Trigger-After-Settle", eventMsg)

	page.tmpl.ExecuteTemplate(w, "entity-form", nil)
}
