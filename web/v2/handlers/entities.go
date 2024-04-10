package handlers

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/components/forms"
	"github.com/cayo-rodrigues/nff/web/components/layouts"
	"github.com/cayo-rodrigues/nff/web/components/pages"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/gofiber/fiber/v2"
)

func EntitiesPage(c *fiber.Ctx) error {
	entities, err := services.ListEntities(c)
	if err != nil {
		return err
	}
	return Render(c, layouts.Base(pages.EntitiesPage(entities)))
}

func CreateEntityPage(c *fiber.Ctx) error {
	entity := models.NewEntity()
	return Render(c, layouts.Base(pages.EntityFormPage(entity)))
}

func EditEntityPage(c *fiber.Ctx) error {
	entityID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	entity, err := services.RetrieveEntity(c, entityID)
	if err != nil {
		return err
	}
	return Render(c, layouts.Base(pages.EntityFormPage(entity)))
}

func CreateEntity(c *fiber.Ctx) error {
	entity := models.NewEntityFromForm(c)
	if !entity.IsValid() {
		return RetargetToForm(c, "entity", forms.EntityForm(entity))
	}

	err := services.CreateEntity(c, entity)
	if err != nil {
		return err
	}

	return c.Redirect("/entities")
}

func UpdateEntity(c *fiber.Ctx) error {
	entity := models.NewEntityFromForm(c)
	if !entity.IsValid() {
		return RetargetToForm(c, "entity", forms.EntityForm(entity))
	}

	err := services.UpdateEntity(c, entity)
	if err != nil {
		return err
	}

	return c.Redirect("/entities")
}
