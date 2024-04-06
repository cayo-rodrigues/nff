package handlers

import (
	"github.com/cayo-rodrigues/nff/web/components/layouts"
	"github.com/cayo-rodrigues/nff/web/components/pages"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/gofiber/fiber/v2"
)

func EntitiesPage(c *fiber.Ctx) error {
	entities := []*models.Entity{
		{Name: "Kira", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Produtor Rural"},
		{Name: "Limão", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Produtor Rural"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Ivy", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Apenas Destinatário"},
		{Name: "Cay", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Produtor Rural"},
		{Name: "Gui", CpfCnpj: "123123123", Ie: "99999990088712", UserType: "Produtor Rural"},
	}
	return Render(c, layouts.Base(pages.EntitiesPage(entities)))
}
