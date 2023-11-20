package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
)

type EntitiesPage struct {
	service interfaces.EntityService
}

func NewEntitiesPage(entityService interfaces.EntityService) *EntitiesPage {
	return &EntitiesPage{
		service: entityService,
	}
}

func (p *EntitiesPage) NewEmptyData() fiber.Map {
	return fiber.Map{
		"Entity":           models.NewEmptyEntity(),
		"FormSelectFields": models.NewEntityFormSelectFields(),
	}
}

func (p *EntitiesPage) Render(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()
	userID := c.Locals("UserID").(int)

	entities, err := p.service.ListEntities(c.Context(), userID)
	if err != nil {
		pageData["GeneralError"] = err.Error()
		c.Set("HX-Trigger-After-Settle", "general-error")
	}

	pageData["Entities"] = entities

	return c.Render("entities", pageData, "layouts/base")
}

func (p *EntitiesPage) GetEntityForm(c *fiber.Ctx) error {
	pageData := p.NewEmptyData()
	userID := c.Locals("UserID").(int)

	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.EntityNotFoundErr)
	}

	entity, err := p.service.RetrieveEntity(c.Context(), entityID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	entity.IsSelected = true
	pageData["Entity"] = entity

	return c.Render("partials/entity-form", pageData)
}

func (p *EntitiesPage) CreateEntity(c *fiber.Ctx) error {
	userID := c.Locals("UserID").(int)
	entity := models.NewEntityFromForm(c)
	entity.CreatedBy = userID

	if !entity.IsValid() {
		pageData := p.NewEmptyData()
		pageData["Entity"] = entity
		return utils.RetargetToForm(c, "entity", pageData)
	}

	err := p.service.CreateEntity(c.Context(), entity)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Trigger-After-Settle", "entity-created")
	return c.Render("partials/entity-card", entity)
}

func (p *EntitiesPage) UpdateEntity(c *fiber.Ctx) error {
	entity := models.NewEntityFromForm(c)
	entity.IsSelected = true

	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.EntityNotFoundErr)
	}
	entity.ID = entityID

	userID := c.Locals("UserID").(int)
	entity.CreatedBy = userID

	if !entity.IsValid() {
		pageData := p.NewEmptyData()
		pageData["Entity"] = entity
		return utils.RetargetToForm(c, "entity", pageData)
	}

	err = p.service.UpdateEntity(c.Context(), entity)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	c.Set("HX-Trigger-After-Settle", "entity-updated")
	return c.Render("partials/entity-card", entity)
}

func (p *EntitiesPage) DeleteEntity(c *fiber.Ctx) error {
	userID := c.Locals("UserID").(int)
	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.GeneralErrorResponse(c, utils.EntityNotFoundErr)
	}
	err = p.service.DeleteEntity(c.Context(), entityID, userID)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	eventMsg := fmt.Sprintf("{\"entity-deleted\": %v}", entityID)
	c.Set("HX-Trigger-After-Settle", eventMsg)

	return c.Render("partials/entity-form", p.NewEmptyData())
}
