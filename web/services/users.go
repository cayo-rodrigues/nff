package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/cayo-rodrigues/nff/web/utils/cryptoutils"
)

func CreateUser(ctx context.Context, user *models.User) error {
	_, err := storage.RetrieveUser(&storage.RetrieveUserParams{
		Ctx:   ctx,
		Email: user.Email,
	})
	userAlreadyExists := true
	if errors.Is(err, sql.ErrNoRows) {
		userAlreadyExists = false
	} else if err != nil {
		return err
	}
	if userAlreadyExists {
		user.SetError("Email", utils.EmailNotAvailableMsg)
		return &user.Errors
	}

	user.Salt, err = cryptoutils.GenerateSalt()
	if err != nil {
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
	userFromDB, err := storage.RetrieveUser(&storage.RetrieveUserParams{
		Ctx:   ctx,
		Email: user.Email,
	})
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

func IsReauthDataValid(ctx context.Context, user *models.User) bool {
	if user.Password == "" {
		user.SetError("Password", utils.MandatoryFieldMsg)
		return false
	}

	userFromDB, err := storage.RetrieveUser(&storage.RetrieveUserParams{
		Ctx: ctx,
		ID:  user.ID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			user.SetError("Password", utils.WrongPassword)
		}
		return false
	}

	passwordsMatch := utils.IsPasswordCorrect(user.Password, userFromDB.Password)
	if !passwordsMatch {
		user.SetError("Password", utils.WrongPassword)
		return false
	}

	return true
}
