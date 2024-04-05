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
		user.SetError("Email", utils.EmailNotAvailableMsg)
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

func IsLoginDataValid(ctx context.Context, user *models.User) bool {
	userFromDB, err := storage.RetrieveUser(ctx, user.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			user.SetError("Email", utils.InvalidLoginDataMsg)
			user.SetError("Password", utils.InvalidLoginDataMsg)
		}
		return false
	}

	passwordsMatch := utils.IsPasswordCorrect(user.Password, userFromDB.Password)
	if !passwordsMatch {
		user.SetError("Email", utils.InvalidLoginDataMsg)
		user.SetError("Password", utils.InvalidLoginDataMsg)
		return false
	}

	return true
}
