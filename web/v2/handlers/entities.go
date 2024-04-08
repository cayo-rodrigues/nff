package handlers

import (
	"github.com/a-h/templ"
	"github.com/cayo-rodrigues/nff/web/components/forms"
	"github.com/cayo-rodrigues/nff/web/components/layouts"
	"github.com/cayo-rodrigues/nff/web/components/pages"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/gofiber/fiber/v2"
)

func EntitiesPage(c *fiber.Ctx) error {
	entities := []*models.Entity{
		{Name: "Kira", CpfCnpj: "139.503.176-27", Ie: "99999990088712", UserType: "Produtor Rural"},
		{Name: "Limão", CpfCnpj: "44.504.044/0001-24", Ie: "99999990088712", UserType: "Produtor Rural"},
		{Name: "Ivy", CpfCnpj: "44504044000124", Ie: "024.263.939/8624", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "13950317627", Ie: "6222304167764", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Emerson Cássio da Silva", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Antonio Francisco Reginaldo", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Benedito Eugenio da Fonseca Junior e outro(s)", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Cay", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Produtor Rural"},
		{Name: "Gui", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Produtor Rural"},
	}
	return Render(c, layouts.Base(pages.EntitiesPage(entities)))
}

func CreateEntityPage(c *fiber.Ctx) error {
	entity := models.NewEntity()
	return Render(c, layouts.Base(pages.EntityFormPage(entity)))
}

func EditEntityPage(c *fiber.Ctx) error {
	entity := models.NewEntity()
	entity.ID = 1
	return Render(c, layouts.Base(pages.EntityFormPage(entity)))
}

func CreateEntity(c *fiber.Ctx) error {
	entity := models.NewEntityFromForm(c)
	if !entity.IsValid() {
		return RetargetToForm(c, "entity", forms.EntityForm(entity))
	}
	return Render(c, forms.EntityForm(entity), templ.WithStatus(fiber.StatusCreated))
}
