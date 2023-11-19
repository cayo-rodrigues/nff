package models

import (
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        int
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewEmptyUser() *User {
	return &User{
		Email:    "",
		Password: "",
	}
}

func NewUserFromForm(c *fiber.Ctx) *User {
	return &User{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}
}

func (u *User) Scan(rows db.Scanner) error {
	return rows.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}
