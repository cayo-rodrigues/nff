package handlers

import (
	"fmt"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/services"
	"github.com/cayo-rodrigues/nff/web/ui/forms"
	"github.com/cayo-rodrigues/nff/web/ui/layouts"
	"github.com/cayo-rodrigues/nff/web/ui/pages"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/cayo-rodrigues/nff/web/utils/cryptoutils"
	"github.com/gofiber/fiber/v2"
)

func LoginPage(c *fiber.Ctx) error {
	user := models.NewUser()
	c.Append("HX-Trigger-After-Settle", "highlight-current-page")
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

	key, err := cryptoutils.DeriveKey([]byte(user.Password), user.Salt)
	if err != nil {
		return err
	}

	if err = services.SaveEncryptionKeyInSession(c, key); err != nil {
		return err
	}

	if err = services.SaveUserSession(c, user.ID); err != nil {
		return err
	}

	redis := database.GetDB().Redis
	redis.Publish(c.Context(), "0:operation-finished", 0)

	return RetargetToPageHandler(c, "/entities", EntitiesPage)
}

func LogoutUser(c *fiber.Ctx) error {
	err := services.DestroyUserSession(c)
	if err != nil {
		return err
	}

	userID := utils.GetUserData(c.Context()).ID
	redis := database.GetDB().Redis
	redis.Publish(c.Context(), fmt.Sprintf("%d:operation-finished", userID), 0)

	return RetargetToPageHandler(c, "/login", LoginPage)
}

func ReauthUser(c *fiber.Ctx) error  {
	user := models.NewUserFromForm(c)
	user.ID = utils.GetUserID(c.Context())

	if !services.IsReauthDataValid(c.Context(), user) {
		return RetargetToForm(c, "reauth", forms.ReauthenticateForm(user))
	}

	key, err := cryptoutils.DeriveKey([]byte(user.Password), user.Salt)
	if err != nil {
		return err
	}

	if err = services.SaveEncryptionKeyInSession(c, key); err != nil {
		return err
	}

	c.Append("HX-Trigger-After-Settle", "close-reauth-form-dialog")
	return Render(c, forms.ReauthenticateForm(user))
}
