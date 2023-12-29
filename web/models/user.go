package models

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/globals"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type UserFormErrors struct {
	Email    string
	Password string
}

type User struct {
	ID        int
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Errors    *UserFormErrors
}

func NewEmptyUser() *User {
	return &User{
		Email:    "",
		Password: "",
		Errors:   &UserFormErrors{},
	}
}

func NewUserFromForm(c *fiber.Ctx) *User {
	return &User{
		Email:    strings.TrimSpace(c.FormValue("email")),
		Password: strings.TrimSpace(c.FormValue("password")),
		Errors:   &UserFormErrors{},
	}
}

func (u *User) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := globals.MandatoryFieldMsg
	valueTooLongMsg := globals.ValueTooLongMsg
	invalidFormatMsg := globals.InvalidFormatMsg

	hasEmail := u.Email != ""
	hasPassword := u.Password != ""

	hasValidEmailFormat := globals.ReEmail.MatchString(u.Email)

	emailTooLong := utf8.RuneCount([]byte(u.Email)) > 128

	fields := [4]*utils.Field{
		{ErrCondition: !hasEmail, ErrField: &u.Errors.Email, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasPassword, ErrField: &u.Errors.Password, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: hasEmail && !hasValidEmailFormat, ErrField: &u.Errors.Email, ErrMsg: &invalidFormatMsg},
		{ErrCondition: emailTooLong, ErrField: &u.Errors.Email, ErrMsg: &valueTooLongMsg},
	}

	for _, field := range fields {
		utils.ValidateField(field, &isValid)
	}

	return isValid
}

func (u *User) Scan(rows db.Scanner) error {
	return rows.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}
