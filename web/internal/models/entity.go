package models

import (
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
	"github.com/gofiber/fiber/v2"
)

type EntityFormSelectFields struct {
	UserTypes   *[14]string
	StreetTypes *[3]string
}

type EntityFormError struct {
	Name       string
	UserType   string
	CpfCnpj    string
	Ie         string
	Email      string
	PostalCode string
	StreetType string
}

type Address struct {
	PostalCode   string
	Neighborhood string
	StreetType   string
	StreetName   string
	Number       string
}

type Entity struct {
	Id         int
	Name       string
	UserType   string
	Ie         string
	CpfCnpj    string
	Email      string
	Password   string
	IsSelected bool
	Address    *Address
	Errors     *EntityFormError
}

func NewEntityFormSelectFields() *EntityFormSelectFields {
	return &EntityFormSelectFields{
		UserTypes:   &globals.EntityUserTypes,
		StreetTypes: &globals.EntityAddressStreetTypes,
	}
}

func NewEmptyEntity() *Entity {
	return &Entity{
		Address: &Address{},
		Errors:  &EntityFormError{},
	}
}

func NewEntityFromForm(c *fiber.Ctx) *Entity {
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		id = 0
	}
	return &Entity{
		Id:       id,
		Name:     c.FormValue("name"),
		UserType: c.FormValue("user_type"),
		CpfCnpj:  c.FormValue("cpf_cnpj"),
		Ie:       c.FormValue("ie"),
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
		Address: &Address{
			PostalCode:   c.FormValue("postal_code"),
			Neighborhood: c.FormValue("neighborhood"),
			StreetType:   c.FormValue("street_type"),
			StreetName:   c.FormValue("street_name"),
			Number:       c.FormValue("number"),
		},
		Errors: &EntityFormError{},
	}
}

func (e *Entity) Scan(rows sql.Scanner) error {
	return rows.Scan(
		&e.Id, &e.Name, &e.UserType, &e.CpfCnpj, &e.Ie, &e.Email, &e.Password,
		&e.Address.PostalCode, &e.Address.Neighborhood, &e.Address.StreetType, &e.Address.StreetName, &e.Address.Number,
	)
}

func (e *Entity) IsValid() bool {
	// TODO async validation, so that all validations happen at once

	isValid := true
	if e.Name == "" {
		e.Errors.Name = "Campo obrigatório"
		isValid = false
	}

	if e.Ie == "" && (e.Address.PostalCode == "" || e.Address.Neighborhood == "" || e.Address.StreetType == "" || e.Address.StreetName == "" || e.Address.Number == "") {
		e.Errors.Ie = "Ie OU endereço completo obrigatórios"
		isValid = false
	}

	if e.Ie != "" && !globals.ReIeMg.MatchString(e.Ie) {
		e.Errors.Ie = "Formato inválido"
		isValid = false
	}

	if e.CpfCnpj != "" && !globals.ReCpf.MatchString(e.CpfCnpj) && !globals.ReCnpj.MatchString(e.CpfCnpj) {
		e.Errors.CpfCnpj = "Formato inválido"
		isValid = false
	}

	if e.Email != "" && !globals.ReEmail.MatchString(e.Email) {
		e.Errors.Email = "Formato inválido"
		isValid = false
	}

	if e.Address.PostalCode != "" && !globals.RePostalCode.MatchString(e.Address.PostalCode) {
		e.Errors.PostalCode = "Formato inválido"
		isValid = false
	}

	if e.UserType != "" {
		hasValidUserType := false
		for _, userType := range globals.EntityUserTypes {
			if e.UserType == userType {
				hasValidUserType = true
				break
			}
		}
		if !hasValidUserType {
			e.Errors.UserType = "Valor inaceitável"
			isValid = false
		}
	}

	if e.Address.StreetType != "" && e.Address.StreetType != "Rua" {
		hasValidStreetType := false
		for _, streetType := range globals.EntityAddressStreetTypes {
			if e.Address.StreetType == streetType {
				hasValidStreetType = true
				break
			}
		}
		if !hasValidStreetType {
			e.Errors.StreetType = "Valor inaceitável"
			isValid = false
		}
	}

	return isValid
}
