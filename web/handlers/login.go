package handlers

import (
	"errors"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/interfaces"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type LoginPage struct {
	userService interfaces.UserService
}

func NewLoginPage(userService interfaces.UserService) *LoginPage {
	return &LoginPage{
		userService: userService,
	}
}

func (p *LoginPage) Render(c *fiber.Ctx) error {
	return c.Render("login", nil, "layouts/base")
}

func (p *LoginPage) Login(c *fiber.Ctx) error {
	loginData := models.NewUserFromForm(c)
	formData := fiber.Map{
		"Email":    loginData.Email,
		"Password": loginData.Password,
	}

	if !loginData.IsValid() {
		formData["Errors"] = loginData.Errors
		return utils.RetargetToForm(c, "login", formData)
	}

	user, err := p.userService.RetrieveUser(c.Context(), loginData.Email)
	if errors.Is(err, utils.UserNotFoundErr) {
		loginData.Errors.Email = err.Error()
		formData["Errors"] = loginData.Errors
		return utils.RetargetToForm(c, "login", formData)
	}
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}

	passwordsMatch := utils.IsPasswordCorrect(loginData.Password, user.Password)
	if !passwordsMatch {
		loginData.Errors.Password = "Senha incorreta"
		formData["Errors"] = loginData.Errors
		return utils.RetargetToForm(c, "login", formData)
	}

	sess, err := db.SessionStore.Get(c)
	if err != nil {
		return err
	}
	sess.Set("IsAuthenticated", true)
	sess.Set("UserID", user.ID)
	err = sess.Save()
	if err != nil {
		return err
	}

	return c.Redirect("/entities", fiber.StatusFound)
}

func Logout(c *fiber.Ctx) error {
	sess, err := db.SessionStore.Get(c)
	if err != nil {
		return err
	}
	sess.Set("IsAuthenticated", false)
	sess.Set("UserID", 0)
	err = sess.Save()
	if err != nil {
		return err
	}

	return c.Redirect("/")
}
