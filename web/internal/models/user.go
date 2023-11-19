package models

import (
	"sync"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
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
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		Errors:   &UserFormErrors{},
	}
}

func (u *User) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := "Campo obrigatório"
	invalidFormatMsg := "Formato inválido"
	validationsCount := 3

	var wg sync.WaitGroup
	wg.Add(validationsCount)
	ch := make(chan bool, validationsCount)

	go utils.ValidateField(u.Email == "", &u.Errors.Email, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(!globals.ReEmail.MatchString(u.Email), &u.Errors.Email, &invalidFormatMsg, ch, &wg)
	go utils.ValidateField(u.Password == "", &u.Errors.Password, &mandatoryFieldMsg, ch, &wg)

	wg.Wait()
	close(ch)

	for i := 0; i < validationsCount; i++ {
		if validationPassed := <-ch; !validationPassed {
			isValid = false
			break
		}
	}

	return isValid
}

func (u *User) Scan(rows db.Scanner) error {
	return rows.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}
