package services

import (
	"context"
	"errors"
	"log"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/jackc/pgx/v5"
)

type EntityService struct{}

func NewEntityService() *EntityService {
	return &EntityService{}
}

// TODO accept filters
func (s *EntityService) ListEntities(ctx context.Context) ([]*models.Entity, error) {
	rows, _ := db.PG.Query(ctx, "SELECT * FROM entities ORDER BY id")
	defer rows.Close()

	entities := []*models.Entity{}

	for rows.Next() {
		entity := models.NewEmptyEntity()
		err := entity.Scan(rows)
		if err != nil {
			log.Println("Error scaning entity rows: ", err)
			return nil, utils.InternalServerErr
		}
		entities = append(entities, entity)
	}

	return entities, nil
}

func (s *EntityService) RetrieveEntity(ctx context.Context, entityId int) (*models.Entity, error) {
	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM entities WHERE entities.id = $1",
		entityId,
	)

	entity := models.NewEmptyEntity()
	err := entity.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Entity with id %v not found: %v", entityId, err)
		return nil, utils.EntityNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning entity row, likely because it has not been found: ", err)
		return nil, utils.InternalServerErr
	}

	return entity, nil
}

func (s *EntityService) CreateEntity(ctx context.Context, entity *models.Entity) error {
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO entities (name, user_type, cpf_cnpj, ie, email, password, postal_code, neighborhood, street_type, street_name, number)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
		`,
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
		entity.Address.PostalCode, entity.Address.Neighborhood, entity.Address.StreetType, entity.Address.StreetName, entity.Address.Number,
		entity.CreatedBy,
	)
	err := row.Scan(&entity.ID)
	if err != nil {
		log.Println("Error when running insert entity query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func (s *EntityService) UpdateEntity(ctx context.Context, entity *models.Entity) error {
	result, err := db.PG.Exec(
		ctx,
		`UPDATE entities
			SET name = $1, user_type = $2, cpf_cnpj = $3, ie = $4, email = $5,
				password = $6, postal_code = $7, neighborhood = $8, street_type = $9,
				street_name = $10, number = $11
		WHERE id = $12`,
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
		entity.Address.PostalCode, entity.Address.Neighborhood, entity.Address.StreetType, entity.Address.StreetName, entity.Address.Number,
		entity.ID,
	)
	if err != nil {
		log.Println("Error when running update entity query: ", err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Entity with id %v not found when running update query", entity.ID)
		return utils.EntityNotFoundErr
	}

	return nil
}

func (s *EntityService) DeleteEntity(ctx context.Context, entityId int) error {
	result, err := db.PG.Exec(
		ctx,
		"DELETE FROM entities WHERE id = $1",
		entityId,
	)
	if err != nil {
		log.Println("Error when running delete entity query: ", err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Entity with id %v not found when running delete query", entityId)
		return utils.EntityNotFoundErr
	}

	return nil
}
