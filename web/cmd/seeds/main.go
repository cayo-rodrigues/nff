package main

import (
	"context"
	"fmt"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/globals"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func main() {
	dbpool := db.GetDBPool()
	defer dbpool.Close()

	fmt.Println("Planting seeds...")

	userID, err := UserSeeds()
	if err != nil {
		return
	}

	err = EntitySeeds(userID)
	if err != nil {
		return
	}

	fmt.Println("A brand new garden OK!")
}

func UserSeeds() (int, error) {
	tx, err := db.PG.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())

	hashedPassword, err := utils.HashPassword("asdf")
	if err != nil {
		fmt.Printf("Something went wrong trying to hash user sample password: %v\n", err)
		return 0, err
	}

	user := &models.User{
		Email:    "foo@bar.baz",
		Password: hashedPassword,
	}

	row := tx.QueryRow(
		context.Background(),
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		user.Email, user.Password,
	)
	err = row.Scan(&user.ID)
	if err != nil {
		fmt.Printf("Something went wrong trying to insert user sample in db: %v\n", err)
		return 0, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return 0, err
	}

	fmt.Println("User seeds OK")

	return user.ID, nil
}

func EntitySeeds(userID int) error {
	entities := []*models.Entity{
		{
			Name:      "Foo",
			UserType:  globals.EntityUserTypes[0],
			CpfCnpj:   "13213486060",
			Ie:        "1134711168910",
			Email:     "bar@baz.com",
			Password:  "123456",
			Address:   &models.Address{},
			CreatedBy: userID,
		},
		{
			Name:      "Bar",
			UserType:  globals.EntityUserTypes[2],
			Ie:        "9970020436725",
			Address:   &models.Address{},
			CreatedBy: userID,
		},
	}

	tx, err := db.PG.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	for i, entity := range entities {
		_, err := tx.Exec(
			context.Background(),
			`INSERT INTO entities (
                    name, user_type, cpf_cnpj, ie, email, password,
                    postal_code, neighborhood, street_type, street_name, number,
                    created_by
                )
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
            `,
			entity.Name, entity.UserType, entity.CpfCnpj, entity.Ie, entity.Email, entity.Password,
			entity.PostalCode, entity.Neighborhood, entity.StreetType, entity.StreetName, entity.Number,
			entity.CreatedBy,
		)

		if err != nil {
			fmt.Printf("Something went wrong trying to insert entity sample #%d in db: %v\n", i, err)
			return err
		}
	}

	fmt.Println("Entity seeds OK")
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil

}
