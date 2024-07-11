package handlers

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/components"
	"github.com/cayo-rodrigues/nff/web/ui/forms"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func EntitiesPage(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	entities, err := services.ListEntities(c.Context(), userID)
	if err != nil {
		return err
	}
	c.Append("HX-Trigger-After-Settle", "highlight-current-page")
	return Render(c, layouts.Base(pages.EntitiesPage(entities)))
}

func CreateEntityPage(c *fiber.Ctx) error {
	entity := models.NewEntity()
	return Render(c, layouts.Base(pages.EntityFormPage(entity)))
}

func EditEntityPage(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	entity, err := services.RetrieveEntity(c.Context(), entityID, userID)
	if err != nil {
		return err
	}
	return Render(c, layouts.Base(pages.EntityFormPage(entity)))
}

func CreateEntity(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	entity := models.NewEntityFromForm(c)
	if !entity.IsValid() {
		return RetargetToForm(c, "entity", forms.EntityForm(entity))
	}

	err := services.CreateEntity(c.Context(), entity, userID)
	if err != nil {
		return err
	}

	return RetargetToPageHandler(c, "/entities", EntitiesPage)
}

func UpdateEntity(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	entity := models.NewEntityFromForm(c)
	if !entity.IsValid() {
		return RetargetToForm(c, "entity", forms.EntityForm(entity))
	}

	err := services.UpdateEntity(c.Context(), entity, userID)
	if err != nil {
		return err
	}

	return RetargetToPageHandler(c, "/entities", EntitiesPage)
}

func DeleteEntity(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	err = services.DeleteEntity(c.Context(), entityID, userID)
	if err != nil {
		return err
	}
	return Render(c, components.Nothing())
}

func SearchEntities(c *fiber.Ctx) error {
	userID := utils.GetUserData(c.Context()).ID
	filters := c.Queries()
	entities, err := services.ListEntities(c.Context(), userID, filters)
	if err != nil {
		return err
	}
	return Render(c, components.EntityList(entities))
}
