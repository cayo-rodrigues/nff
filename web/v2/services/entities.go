package services

import (
	"context"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
)

func CreateEntity(ctx context.Context, entity *models.Entity, userID int) error {
	entity.CreatedBy = userID
	return storage.CreateEntity(ctx, entity)
}

func ListEntities(ctx context.Context, userID int, filters ...map[string]string) ([]*models.Entity, error) {
	f := models.NewFilters().Where("created_by = ").Placeholder(userID)

	for _, filter := range filters {
		if name, ok := filter["name"]; ok {
			f.And("name").ILike().WildPlaceholder(name)
		}
	}

	f.OrderBy("name")

	return storage.ListEntities(ctx, userID, f)
}

func RetrieveEntity(ctx context.Context, entityID int, userID int) (*models.Entity, error) {
	entity, err := storage.RetrieveEntity(ctx, entityID, userID)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func UpdateEntity(ctx context.Context, entity *models.Entity, userID int) error {
	entity.CreatedBy = userID
	return storage.UpdateEntity(ctx, entity)
}

func DeleteEntity(ctx context.Context, entityID int, userID int) error {
	return storage.DeleteEntity(ctx, entityID, userID)
}
