package services

import (
	"context"
	"encoding/hex"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/cayo-rodrigues/nff/web/utils/cryptoutils"
)

func CreateEntity(ctx context.Context, encryptionKey []byte, entity *models.Entity) error {
	userID := utils.GetUserID(ctx)
	entity.CreatedBy = userID

	cyphertext, err := cryptoutils.Encrypt(encryptionKey, []byte(entity.Password))
	if err != nil {
		return err
	}
	entity.Password = hex.EncodeToString(cyphertext)

	return storage.CreateEntity(ctx, entity)
}

func ListEntities(ctx context.Context, filters ...map[string]string) ([]*models.Entity, error) {
	userID := utils.GetUserID(ctx)

	f := models.NewFilters().Where("created_by = ").Placeholder(userID)

	for _, filter := range filters {
		if q, ok := filter["q"]; ok {
			f.Raw("AND ( ")
			f.Raw("name").ILike().WildPlaceholder(q)
			f.Or("cpf_cnpj").ILike().WildPlaceholder(q)
			f.Or("ie").ILike().WildPlaceholder(q)
			f.Or("user_type").ILike().WildPlaceholder(q)
			f.Or("email").ILike().WildPlaceholder(q)
			f.Raw(" ) ")
		}
	}

	f.OrderBy("name")

	return storage.ListEntities(ctx, userID, f)
}

func RetrieveEntity(ctx context.Context, entityID int) (*models.Entity, error) {
	userID := utils.GetUserID(ctx)
	return storage.RetrieveEntity(ctx, entityID, userID)
}

func UpdateEntity(ctx context.Context, encryptionKey []byte, entity *models.Entity) error {
	userID := utils.GetUserID(ctx)
	entity.CreatedBy = userID

	cyphertext, err := cryptoutils.Encrypt(encryptionKey, []byte(entity.Password))
	if err != nil {
		return err
	}
	entity.Password = hex.EncodeToString(cyphertext)

	return storage.UpdateEntity(ctx, entity)
}

func DeleteEntity(ctx context.Context, entityID int) error {
	userID := utils.GetUserID(ctx)
	return storage.DeleteEntity(ctx, entityID, userID)
}
