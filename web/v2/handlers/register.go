package handlers

import (
	"github.com/cayo-rodrigues/nff/web/components/forms"
	"github.com/cayo-rodrigues/nff/web/components/layouts"
	"github.com/cayo-rodrigues/nff/web/components/pages"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/gofiber/fiber/v2"
)

func RegisterPage(c *fiber.Ctx) error {
	user := models.NewUser()
	return Render(c, layouts.Base(pages.RegisterPage(user)))
}

func RegisterUser(c *fiber.Ctx) error {
	user := models.NewUserFromForm(c)
	if !user.IsValid() {
		return RetargetToForm(c, "register", forms.RegisterForm(user))
	}

	err := services.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}

	err = services.SaveUserSession(c, user.ID)
	if err != nil {
		return err
	}

	return c.Redirect("/entities")
}
