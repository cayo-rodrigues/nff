package storage

import (
	"context"
	"log"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
)

type RetrieveUserParams struct {
	Ctx context.Context
	Email string
	ID int
}

func RetrieveUser(params *RetrieveUserParams) (*models.User, error) {
	db := database.GetDB()

	row := db.SQLite.QueryRowContext(
		params.Ctx,
		"SELECT * FROM users WHERE users.email = ? OR users.id = ?",
		params.Email,
		params.ID,
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
		`INSERT INTO users (email, password, salt) VALUES (?, ?, ?) RETURNING id`,
		user.Email, user.Password, user.Salt,
	)
	err := row.Scan(&user.ID)
	if err != nil {
		log.Println("Error when running insert user query: ", err)
		return err
	}

	return nil
}
