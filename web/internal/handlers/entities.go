package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/cayo-rodrigues/nff/web/internal/workers"
)

type EntitiesPage struct{}

type EntitiesPageData struct {
	IsAuthenticated  bool
	Entities         []*models.Entity
	Entity           *models.Entity
	GeneralError     string
	FormSelectFields *models.EntityFormSelectFields
}

func (page *EntitiesPage) NewEmptyData() *EntitiesPageData {
	return &EntitiesPageData{
		IsAuthenticated:  true,
		Entity:           models.NewEmptyEntity(),
		FormSelectFields: models.NewEntityFormSelectFields(),
	}
}

func (page *EntitiesPage) Render(c *fiber.Ctx) error {
	data := page.NewEmptyData()
	time.Sleep(time.Second * 3)
	entities, err := workers.ListEntities(c.Context())
	if err != nil {
		data.GeneralError = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
	}

	data.Entities = entities

	return c.Render("entities", data, "layouts/base")
}

func (page *EntitiesPage) GetEntityForm(c *fiber.Ctx) error {
	data := page.NewEmptyData()

	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.EntityNotFoundErr)
	}

	entity, err := workers.RetrieveEntity(c.Context(), entityId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	data.Entity = entity
	data.Entity.IsSelected = true

	return c.Render("partials/entity-form", data)
}

func (page *EntitiesPage) CreateEntity(c *fiber.Ctx) error {
	entity := models.NewEntityFromForm(c)

	if !entity.IsValid() {
		data := page.NewEmptyData()
		data.Entity = entity
		c.Set("HX-Retarget", "#entity-form")
		c.Set("HX-Reswap", "outerHTML")
		return c.Render("partials/entity-form", data)
	}

	err := workers.CreateEntity(c.Context(), entity)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Trigger-After-Settle", "entity-created")
	return c.Render("partials/entity-card", entity)
}

func (page *EntitiesPage) UpdateEntity(c *fiber.Ctx) error {
	entity := models.NewEntityFromForm(c)
	entity.IsSelected = true

	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.EntityNotFoundErr)
	}
	entity.Id = entityId

	if !entity.IsValid() {
		data := page.NewEmptyData()
		data.Entity = entity
		c.Set("HX-Retarget", "#entity-form")
		c.Set("HX-Reswap", "outerHTML")
		return c.Render("partials/entity-form", data)
	}

	err = workers.UpdateEntity(c.Context(), entity)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Trigger-After-Settle", "entity-updated")
	return c.Render("partials/entity-card", entity)
}

func (page *EntitiesPage) DeleteEntity(c *fiber.Ctx) error {
	entityId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.EntityNotFoundErr)
	}
	err = workers.DeleteEntity(c.Context(), entityId)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	eventMsg := fmt.Sprintf("{\"entity-deleted\": %v}", entityId)
	c.Set("HX-Trigger-After-Settle", eventMsg)

	return c.Render("partials/entity-form", page.NewEmptyData())
}
