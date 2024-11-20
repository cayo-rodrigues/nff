package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
)

func ListEntities(ctx context.Context, userID int, filters *models.Filters) (entities []*models.Entity, err error) {
	namespace := "entities"
	filtersKey := filters.String() + filters.StringValues()

	db := database.GetDB()

	if db.Redis.GetDecodedCache(ctx, userID, namespace, filtersKey, &entities); entities != nil {
		return entities, nil
	}

	query := new(strings.Builder)
	query.WriteString("SELECT * FROM entities")
	query.WriteString(filters.String())

	rows, err := db.SQLite.QueryContext(ctx, query.String(), filters.Values()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		entity := models.NewEntity()
		err := Scan(rows, entity)
		if err != nil {
			log.Println("Error scaning entity rows: ", err)
			return nil, err
		}
		err = json.Unmarshal(entity.OtherIesJSON, &entity.OtherIes)
		if err != nil {
			log.Println("Error unmarshaling entity other ies: ", err)
			return nil, err
		}
		entities = append(entities, entity)
	}

	db.Redis.SetEncodedCache(ctx, userID, namespace, filtersKey, entities, time.Hour)

	return entities, nil
}

func CreateEntity(ctx context.Context, entity *models.Entity) (err error) {
	db := database.GetDB()

	entity.OtherIesJSON, err = json.Marshal(entity.OtherIes)
	if err != nil {
		return err
	}

	row := db.SQLite.QueryRowContext(
		ctx,
		`INSERT INTO entities (
			name, user_type, cpf_cnpj, ie, email, password,
			postal_code, neighborhood, street_type, street_name, number,
			created_by, other_ies
		)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, created_at, updated_at
		`,
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
		entity.PostalCode, entity.Neighborhood, entity.StreetType, entity.StreetName, entity.Number,
		entity.CreatedBy, entity.OtherIesJSON,
	)
	err = row.Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)
	if err != nil {
		log.Println("Error when running insert entity query: ", err)
		return err
	}

	return nil
}

func RetrieveEntity(ctx context.Context, entityID int, userID int) (*models.Entity, error) {
	db := database.GetDB()

	row := db.SQLite.QueryRowContext(
		ctx,
		"SELECT * FROM entities WHERE entities.id = ? AND created_by = ?",
		entityID, userID,
	)

	entity := models.NewEntity()
	err := Scan(row, entity)
	if errors.Is(err, sql.ErrNoRows) {
		log.Printf("Entity with id %v not found: %v", entityID, err)
		return nil, err
	}
	if err != nil {
		log.Println("Error scaning entity row, likely because it has not been found: ", err)
		return nil, err
	}

	return entity, nil
}

func UpdateEntity(ctx context.Context, entity *models.Entity) (err error) {
	db := database.GetDB()

	entity.OtherIesJSON, err = json.Marshal(entity.OtherIes)
	if err != nil {
		return err
	}

	result, err := db.SQLite.ExecContext(
		ctx,
		`UPDATE entities
			SET name = ?, user_type = ?, cpf_cnpj = ?, ie = ?, email = ?,
				password = ?, postal_code = ?, neighborhood = ?, street_type = ?,
				street_name = ?, number = ?, updated_at = ?, other_ies = ?
		WHERE id = ? AND created_by = ?`,
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
		entity.PostalCode, entity.Neighborhood, entity.StreetType, entity.StreetName, entity.Number, time.Now(), entity.OtherIesJSON,
		entity.ID, entity.CreatedBy,
	)
	if err != nil {
		log.Println("Error when running update entity query: ", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Entity with id %v not found when running update query", entity.ID)
		return err
	}

	return nil
}

func DeleteEntity(ctx context.Context, entityID int, userID int) error {
	db := database.GetDB()

	result, err := db.SQLite.ExecContext(
		ctx,
		"DELETE FROM entities WHERE id = ? AND created_by = ?",
		entityID, userID,
	)
	if err != nil {
		log.Println("Error when running delete entity query: ", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error when getting rows affected by delete entity query: ", err)
		return err
	}
	if rowsAffected == 0 {
		log.Printf("Entity with id %v not found when running delete query", entityID)
		return err
	}

	return nil
}
