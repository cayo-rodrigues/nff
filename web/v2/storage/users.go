package storage

import (
	"context"
	"errors"
	"log"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/jackc/pgx/v5"
)

func RetrieveUser(ctx context.Context, email string) (*models.User, error) {
	db := database.GetDB()

	row := db.PG.QueryRow(
		ctx,
		"SELECT * FROM users WHERE users.email = $1",
		email,
	)

	user := models.NewUser()
	err := Scan(row, user)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Printf("User with email %v not found: %v", email, err)
		return nil, err
	}
	if err != nil {
		log.Println("Error scaning user row: ", err)
		return nil, err
	}

	return user, nil
}

func CreateUser(ctx context.Context, user *models.User) error {
	db := database.GetDB()

	row := db.PG.QueryRow(
		ctx,
		`INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`,
		user.Email, user.Password,
	)
	err := row.Scan(&user.ID)
	if err != nil {
		log.Println("Error when running insert user query: ", err)
		return err
	}

	return nil
}
