package services

import (
	"context"
	"errors"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/jackc/pgx/v5"
)

func CreateUser(ctx context.Context, user *models.User) error {
	_, err := storage.RetrieveUser(ctx, user.Email)
	userAlreadyExists := true
	if errors.Is(err, pgx.ErrNoRows) {
		userAlreadyExists = false
	} else if err != nil {
		return err
	}
	if userAlreadyExists {
		user.Errors["Email"] = utils.EmailNotAvailableMsg
		return err
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	err = storage.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
