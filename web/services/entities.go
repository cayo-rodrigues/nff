package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/jackc/pgx/v5"
)

type EntityService struct{}

func NewEntityService() *EntityService {
	return &EntityService{}
}


func (s *EntityService) ListEntities(ctx context.Context, userID int) (entities []*models.Entity, err error) {
    namespace := "entities"

	if utils.GetDecodedCache(ctx, userID, namespace, &entities); entities != nil {
        return entities, nil
    }

	rows, _ := db.PG.Query(ctx, "SELECT * FROM entities WHERE created_by = $1 ORDER BY name", userID)
	defer rows.Close()

	for rows.Next() {
		entity := models.NewEmptyEntity()
		err := entity.Scan(rows)
		if err != nil {
			log.Println("Error scaning entity rows: ", err)
			return nil, utils.InternalServerErr
		}
		entities = append(entities, entity)
	}

	utils.SetEncodedCache(ctx, userID, namespace, entities, time.Hour)

	return entities, nil
}

func (s *EntityService) CreateEntity(ctx context.Context, entity *models.Entity) error {
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO entities (name, user_type, cpf_cnpj, ie, email, password, postal_code, neighborhood, street_type, street_name, number, created_by, other_ies)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
		`,
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
		entity.Address.PostalCode, entity.Address.Neighborhood, entity.Address.StreetType, entity.Address.StreetName, entity.Address.Number,
		entity.CreatedBy, entity.OtherIes,
	)
	err := row.Scan(&entity.ID)
	if err != nil {
		log.Println("Error when running insert entity query: ", err)
		return utils.InternalServerErr
	}

	return nil
}

func (s *EntityService) RetrieveEntity(ctx context.Context, entityID int, userID int) (*models.Entity, error) {
	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM entities WHERE entities.id = $1 AND created_by = $2",
		entityID, userID,
	)

	entity := models.NewEmptyEntity()
	err := entity.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("Entity with id %v not found: %v", entityID, err)
		return nil, utils.EntityNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning entity row, likely because it has not been found: ", err)
		return nil, utils.InternalServerErr
	}

	return entity, nil
}

func (s *EntityService) UpdateEntity(ctx context.Context, entity *models.Entity) error {
	result, err := db.PG.Exec(
		ctx,
		`UPDATE entities
			SET name = $1, user_type = $2, cpf_cnpj = $3, ie = $4, email = $5,
				password = $6, postal_code = $7, neighborhood = $8, street_type = $9,
				street_name = $10, number = $11, other_ies = $12
		WHERE id = $13 AND created_by = $14`,
		entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
		entity.Address.PostalCode, entity.Address.Neighborhood, entity.Address.StreetType, entity.Address.StreetName, entity.Address.Number, entity.OtherIes,
		entity.ID, entity.CreatedBy,
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

func (s *EntityService) DeleteEntity(ctx context.Context, entityID int, userID int) error {
	result, err := db.PG.Exec(
		ctx,
		"DELETE FROM entities WHERE id = $1 AND created_by = $2",
		entityID, userID,
	)
	if err != nil {
		log.Println("Error when running delete entity query: ", err)
		return utils.InternalServerErr
	}
	if result.RowsAffected() == 0 {
		log.Printf("Entity with id %v not found when running delete query", entityID)
		return utils.EntityNotFoundErr
	}

	return nil
}
