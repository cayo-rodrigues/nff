package models

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/cayo-rodrigues/safe"
)

type User struct {
	ID        int
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Errors    safe.ErrorMessages
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
	
	fields := safe.Fields{
		{
			Name:  "Email",
			Value: u.Email,
			Rules: safe.Rules{safe.Required(), safe.Email(), safe.Max(128)},
		},
		{
			Name:  "Password",
			Value: u.Password,
			Rules: safe.Rules{safe.Required()},
		},
	}
	errors, ok := safe.Validate(fields)
	u.Errors = errors
	return ok
}

func (u *User) Values() []any {
	return []any{&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt}
}

func (u *User) SetError(key, val string) {
	if u.Errors == nil {
		u.Errors = make(safe.ErrorMessages)
	}
	u.Errors[key] = val
}
