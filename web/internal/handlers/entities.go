package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/services"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
)

type EntitiesPage struct{}

type EntitiesPageData struct {
	IsAuthenticated  bool
	Entities         []*models.Entity
	Entity           *models.Entity
	GeneralError     string
	FormSelectFields *models.EntityFormSelectFields
}

func (p *EntitiesPage) NewEmptyData() *EntitiesPageData {
	return &EntitiesPageData{
		IsAuthenticated:  true,
		Entity:           models.NewEmptyEntity(),
		FormSelectFields: models.NewEntityFormSelectFields(),
	}
}

func (p *EntitiesPage) Render(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()

	entities, err := services.ListEntities(c.Context())
	if err != nil {
		pageData.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
	}

	pageData.Entities = entities

	return c.Render("entities", pageData, "layouts/base")
}

func (p *EntitiesPage) GetEntityForm(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()

	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.EntityNotFoundErr)
	}

	entity, err := services.RetrieveEntity(c.Context(), entityId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	pageData.Entity = entity
	pageData.Entity.IsSelected = true

	return c.Render("partials/entity-form", pageData)
}

func (p *EntitiesPage) CreateEntity(c *fiber.Ctx) error {
	entity := models.NewEntityFromForm(c)

	if !entity.IsValid() {
		pageData := p.NewEmptyData()
		pageData.Entity = entity
		c.Set("HX-Retarget", "#entity-form")
		c.Set("HX-Reswap", "outerHTML")
		return c.Render("partials/entity-form", pageData)
	}

	err := services.CreateEntity(c.Context(), entity)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Trigger-After-Settle", "entity-created")
	return c.Render("partials/entity-card", entity)
}

func (p *EntitiesPage) UpdateEntity(c *fiber.Ctx) error {
	entity := models.NewEntityFromForm(c)
	entity.IsSelected = true

	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.EntityNotFoundErr)
	}
	entity.Id = entityId

	if !entity.IsValid() {
		pageData := p.NewEmptyData()
		pageData.Entity = entity
		c.Set("HX-Retarget", "#entity-form")
		c.Set("HX-Reswap", "outerHTML")
		return c.Render("partials/entity-form", pageData)
	}

	err = services.UpdateEntity(c.Context(), entity)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Trigger-After-Settle", "entity-updated")
	return c.Render("partials/entity-card", entity)
}

func (p *EntitiesPage) DeleteEntity(c *fiber.Ctx) error {
	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.EntityNotFoundErr)
	}
	err = services.DeleteEntity(c.Context(), entityId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	eventMsg := fmt.Sprintf("{\"entity-deleted\": %v}", entityId)
	c.Set("HX-Trigger-After-Settle", eventMsg)

	return c.Render("partials/entity-form", p.NewEmptyData())
}
