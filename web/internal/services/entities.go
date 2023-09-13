package services

import (
	"context"
	"errors"
	"log"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/jackc/pgx/v5"
)

// TODO accept filters
func ListEntities(ctx context.Context) (*[]models.Entity, error) {
	dbpool := sql.GetDatabasePool()
	rows, _ := dbpool.Query(ctx, "SELECT * FROM entities")
	defer rows.Close()

	entities := []models.Entity{}

	for rows.Next() {
		entity := models.Entity{
			Address: &models.Address{},
		}
		err := entity.Scan(rows)
		if err != nil {
			log.Println("Error scaning entity rows: ", err)
			return nil, &utils.InternalServerError{}
		}
		entities = append(entities, entity)
	}

	return &entities, nil
}

func RetrieveEntity(ctx context.Context, entityId int) (*models.Entity, error) {
	dbpool := sql.GetDatabasePool()
	row := dbpool.QueryRow(
		ctx,
		"SELECT * FROM entities WHERE entities.id = $1",
		entityId,
	)

	entity := models.Entity{
		Address: &models.Address{},
	}
	err := entity.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Entity with id %v not found: %v", entityId, err)
		return nil, &utils.EntityNotFoundError{}
	}
	if err != nil {
		log.Println("Error scaning entity row, likely because it has not been found: ", err)
		return nil, &utils.InternalServerError{}
	}

	return &entity, nil
}

func RegisterEntity(ctx context.Context, entity *models.Entity) error {
	passwordHash, err := utils.HashPassword(entity.Password)
	if err != nil {
		log.Println("Error hashing entity password: ", err)
		return &utils.InternalServerError{}
	}

	dbpool := sql.GetDatabasePool()
	row := dbpool.QueryRow(
		ctx,
		"INSERT INTO entities (name, user_type, cpf_cnpj, ie, email, password, postal_code, neighborhood, street_type, street_name, number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id",
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, passwordHash,
		entity.Address.PostalCode, entity.Address.Neighborhood, entity.Address.StreetType, entity.Address.StreetName, entity.Address.Number,
	)
	err = row.Scan(&entity.Id)
	if err != nil {
		log.Println("Error when running insert entity query: ", err)
		return &utils.InternalServerError{}
	}

	return nil
}

func UpdateEntity(ctx context.Context, entity *models.Entity) error {
	dbpool := sql.GetDatabasePool()
	result, err := dbpool.Exec(
		ctx,
		"UPDATE entities SET name = $1, user_type = $2, cpf_cnpj = $3, ie = $4, email = $5, password = $6, postal_code = $7, neighborhood = $8, street_type = $9, street_name = $10, number = $11 WHERE id = $12",
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
		entity.Address.PostalCode, entity.Address.Neighborhood, entity.Address.StreetType, entity.Address.StreetName, entity.Address.Number,
		entity.Id,
	)
	if err != nil {
		log.Println("Error when running update entity query: ", err)
		return &utils.InternalServerError{}
	}
	if result.RowsAffected() == 0 {
		log.Printf("Entity with id %v not found when running update query", entity.Id)
		return &utils.EntityNotFoundError{}
	}

	return nil
}

func DeleteEntity(ctx context.Context, entityId int) error {
	dbpool := sql.GetDatabasePool()
	result, err := dbpool.Exec(
		ctx,
		"DELETE FROM entities WHERE id = $1",
		entityId,
	)
	if err != nil {
		log.Println("Error when running delete entity query: ", err)
		return &utils.InternalServerError{}
	}
	if result.RowsAffected() == 0 {
		log.Printf("Entity with id %v not found when running delete query", entityId)
		return &utils.EntityNotFoundError{}
	}

	return nil
}
