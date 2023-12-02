package handlers

import (
	"errors"
	"fmt"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
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
	formData := fiber.Map{
		"Email":    user.Email,
		"Password": user.Password,
	}

	if !user.IsValid() {
		formData["Errors"] = user.Errors
		return utils.RetargetToForm(c, "register", formData)
	}

	_, err := p.userService.RetrieveUser(c.Context(), user.Email)
	userAlreadyExists := true
	if errors.Is(err, utils.UserNotFoundErr) {
		userAlreadyExists = false
	} else if err != nil {
		return utils.GeneralErrorResponse(c, err)
	}
	if userAlreadyExists {
		user.Errors.Email = "Email indisponível"
		formData["Errors"] = user.Errors
		return utils.RetargetToForm(c, "register", formData)
	}

	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		fmt.Println("Error hashing new user password:", err)
		return utils.GeneralErrorResponse(c, utils.InternalServerErr)
	}

	err = p.userService.CreateUser(c.Context(), user)
	if err != nil {
		return utils.GeneralErrorResponse(c, err)
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

	return c.Redirect("/entities")
}
