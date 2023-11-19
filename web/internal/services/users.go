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

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) RetrieveUser(ctx context.Context, userId int) (*models.User, error) {
	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM users WHERE users.id = $1",
		userId,
	)

	user := models.NewEmptyUser()
	err := user.Scan(row)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("User with id %v not found: %v", userId, err)
		return nil, utils.UserNotFoundErr
	}
	if err != nil {
		log.Println("Error scaning user row, likely because it has not been found: ", err)
		return nil, utils.InternalServerErr
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
		user.Email, user.Password,
	)
	err := row.Scan(&user.ID)
	if err != nil {
		log.Println("Error when running insert user query: ", err)
		return utils.InternalServerErr
	}

	return nil
}
