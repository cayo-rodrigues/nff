package services

import (
	"context"
	"log"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
)

func ListEntities(ctx context.Context, userID int) (entities []*models.Entity, err error) {
	namespace := "entities"

	db := database.GetDB()

	if db.Redis.GetDecodedCache(ctx, userID, namespace, &entities); entities != nil {
		return entities, nil
	}

	rows, _ := db.PG.Query(ctx, "SELECT * FROM entities WHERE created_by = $1 ORDER BY name", userID)
	defer rows.Close()

	for rows.Next() {
		entity := models.NewEntity()
		err := scan(rows, entity)
		if err != nil {
			log.Println("Error scaning entity rows: ", err)
			return nil, err
		}
		entities = append(entities, entity)
	}

	db.Redis.SetEncodedCache(ctx, userID, namespace, entities, time.Hour)

	return entities, nil
}
