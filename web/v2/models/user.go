package models

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        int
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Errors    ErrorMessages
}

func NewUser() *User {
	return &User{}
}

func NewUserFromForm(c *fiber.Ctx) *User {
	return &User{
		Email:    strings.TrimSpace(c.FormValue("email")),
		Password: strings.TrimSpace(c.FormValue("password")),
	}
}

func (u *User) IsValid() bool {
	fields := Fields{
		{
			Name:  "Email",
			Value: u.Email,
			Rules: Rules(Required, Email, Max(128)),
		},
		{
			Name:  "Password",
			Value: u.Password,
			Rules: Rules(Required),
		},
	}
	errors, ok := Validate(fields)
	u.Errors = errors
	return ok
}

func (u *User) Values() []any {
	return []any{&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt}
}
