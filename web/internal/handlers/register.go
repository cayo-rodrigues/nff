package handlers

import (
	"fmt"

	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
	"github.com/cayo-rodrigues/nff/web/internal/middlewares"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type RegisterPage struct {
	userService interfaces.UserService
}

func NewRegisterPage(userService interfaces.UserService) *RegisterPage {
	return &RegisterPage{
		userService: userService,
	}
}

func (p *RegisterPage) Render(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{}, "layouts/base")
}

func (p *RegisterPage) CreateUser(c *fiber.Ctx) error {
	user := models.NewUserFromForm(c)
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		fmt.Println("Error hashing new user password:", err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}
	user.Password = hash

	err = p.userService.CreateUser(c.Context(), user)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	sess, err := middlewares.SessionStore.Get(c)
	if err != nil {
		return err
	}
	sess.Set("IsAuthenticated", true)
	err = sess.Save()
	if err != nil {
		return err
	}

	return c.Redirect("/entities")
}
