package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func CreateUser(ctx context.Context, user *models.User) error {
	_, err := storage.RetrieveUser(ctx, user.Email)
	userAlreadyExists := true
	if errors.Is(err, sql.ErrNoRows) {
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
		if errors.Is(err, sql.ErrNoRows) {
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

	user.ID = userFromDB.ID

	return true
}
