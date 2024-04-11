package handlers

import (
	"github.com/cayo-rodrigues/nff/web/ui/forms"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/gofiber/fiber/v2"
)

func LoginPage(c *fiber.Ctx) error {
	user := models.NewUser()
	return Render(c, layouts.Base(pages.LoginPage(user)))
}

func LoginUser(c *fiber.Ctx) error {
	user := models.NewUserFromForm(c)
	if !user.IsValid() {
		return RetargetToForm(c, "login", forms.LoginForm(user))
	}

	if !services.IsLoginDataValid(c.Context(), user) {
		return RetargetToForm(c, "login", forms.LoginForm(user))
	}

	err := services.SaveUserSession(c, user.ID)
	if err != nil {
		return err
	}

	return RetargetToPageHandler(c, "/entities", EntitiesPage)
}

func LogoutUser(c *fiber.Ctx) error {
	err := services.DestroyUserSession(c)
	if err != nil {
		return err
	}

	return RetargetToPageHandler(c, "/", HomePage)
}
