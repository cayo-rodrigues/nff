package storage

import (
	"context"
	"log"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
)

func RetrieveUser(ctx context.Context, email string) (*models.User, error) {
	db := database.GetDB()

	row := db.SQLite.QueryRowContext(
		ctx,
		"SELECT * FROM users WHERE users.email = ?",
		email,
	)

	user := models.NewUser()
	err := Scan(row, user)
	if err != nil {
		log.Println("Error scaning user row: ", err)
		return nil, err
	}

	return user, nil
}

func CreateUser(ctx context.Context, user *models.User) error {
	db := database.GetDB()

	row := db.SQLite.QueryRowContext(
		ctx,
		`INSERT INTO users (email, password) VALUES (?, ?) RETURNING id`,
		user.Email, user.Password,
	)
	err := row.Scan(&user.ID)
	if err != nil {
		log.Println("Error when running insert user query: ", err)
		return err
	}

	return nil
}
