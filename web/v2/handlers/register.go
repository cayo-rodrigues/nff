package handlers

import (
	"github.com/cayo-rodrigues/nff/web/ui/forms"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/gofiber/fiber/v2"
)

func RegisterPage(c *fiber.Ctx) error {
	user := models.NewUser()
	isAuthenticated := false
	return Render(c, layouts.Base(pages.RegisterPage(user), isAuthenticated))
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

	return RetargetToPageHandler(c, "/entities", EntitiesPage)
}
