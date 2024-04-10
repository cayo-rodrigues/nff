package storage

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/jackc/pgx/v5"
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
		err := Scan(rows, entity)
		if err != nil {
			log.Println("Error scaning entity rows: ", err)
			return nil, err
		}
		entities = append(entities, entity)
	}

	db.Redis.SetEncodedCache(ctx, userID, namespace, entities, time.Hour)

	return entities, nil
}

func CreateEntity(ctx context.Context, entity *models.Entity) error {
	db := database.GetDB()

	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO entities (
			name, user_type, cpf_cnpj, ie, email, password,
			postal_code, neighborhood, street_type, street_name, number,
			created_by
		)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
		`,
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
		entity.PostalCode, entity.Neighborhood, entity.StreetType, entity.StreetName, entity.Number,
		entity.CreatedBy,
	)
	err := row.Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	if err != nil {
		log.Println("Error when running insert entity query: ", err)
		return err
	}

	return nil
}

func RetrieveEntity(ctx context.Context, entityID int, userID int) (*models.Entity, error) {
	db := database.GetDB()

	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM entities WHERE entities.id = $1 AND created_by = $2",
		entityID, userID,
	)

	entity := models.NewEntity()
	err := Scan(row, entity)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Entity with id %v not found: %v", entityID, err)
		return nil, err
	}
	if err != nil {
		log.Println("Error scaning entity row, likely because it has not been found: ", err)
		return nil, err
	}

	return entity, nil
}

func UpdateEntity(ctx context.Context, entity *models.Entity) error {
	db := database.GetDB()

	result, err := db.PG.Exec(
		ctx,
		`UPDATE entities
			SET name = $1, user_type = $2, cpf_cnpj = $3, ie = $4, email = $5,
				password = $6, postal_code = $7, neighborhood = $8, street_type = $9,
				street_name = $10, number = $11, updated_at = $12
		WHERE id = $13 AND created_by = $14`,
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
		entity.PostalCode, entity.Neighborhood, entity.StreetType, entity.StreetName, entity.Number, time.Now(),
		entity.ID, entity.CreatedBy,
	)
	if err != nil {
		log.Println("Error when running update entity query: ", err)
		return err
	}
	if result.RowsAffected() == 0 {
		log.Printf("Entity with id %v not found when running update query", entity.ID)
		return err
	}

	return nil
}

func DeleteEntity(ctx context.Context, entityID int, userID int) error {
	db := database.GetDB()

	result, err := db.PG.Exec(
		ctx,
		"DELETE FROM entities WHERE id = $1 AND created_by = $2",
		entityID, userID,
	)
	if err != nil {
		log.Println("Error when running delete entity query: ", err)
		return err
	}
	if result.RowsAffected() == 0 {
		log.Printf("Entity with id %v not found when running delete query", entityID)
		return err
	}

	return nil
}
