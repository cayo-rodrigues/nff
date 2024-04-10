package services

import (
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func CreateEntity(c *fiber.Ctx, entity *models.Entity) error {
	userID := utils.GetCurrentUserID(c)
	entity.CreatedBy = userID
	return storage.CreateEntity(c.Context(), entity)
}

func ListEntities(c *fiber.Ctx) ([]*models.Entity, error) {
	userID := utils.GetCurrentUserID(c)
	entities, err := storage.ListEntities(c.Context(), userID)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
