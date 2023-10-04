package models

import (
	"strconv"
	"sync"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
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
	isValid := true

	mandatoryFieldMsg := "Campo obrigatório"
	invalidFormatMsg := "Formato inválido"
	eitherIeOrAddrMsg := "Ie OU endereço completo obrigatórios"
	unacceptableValueMsg := "Valor inaceitável"
	validationsCount := 8

	var wg sync.WaitGroup
	wg.Add(validationsCount)
	ch := make(chan bool, validationsCount)

	go utils.ValidateField(e.Name == "", &e.Errors.Name, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(e.Ie == "" && (e.Address.PostalCode == "" || e.Address.Neighborhood == "" || e.Address.StreetType == "" || e.Address.StreetName == "" || e.Address.Number == ""), &e.Errors.Ie, &eitherIeOrAddrMsg, ch, &wg)
	go utils.ValidateField(e.Ie != "" && !globals.ReIeMg.MatchString(e.Ie), &e.Errors.Ie, &invalidFormatMsg, ch, &wg)
	go utils.ValidateField(e.CpfCnpj != "" && !globals.ReCpf.MatchString(e.CpfCnpj) && !globals.ReCnpj.MatchString(e.CpfCnpj), &e.Errors.CpfCnpj, &invalidFormatMsg, ch, &wg)
	go utils.ValidateField(e.Email != "" && !globals.ReEmail.MatchString(e.Email), &e.Errors.Email, &invalidFormatMsg, ch, &wg)
	go utils.ValidateField(e.Address.PostalCode != "" && !globals.RePostalCode.MatchString(e.Address.PostalCode), &e.Errors.PostalCode, &invalidFormatMsg, ch, &wg)

	go utils.ValidateListField(e.UserType, globals.EntityUserTypes[:], &e.Errors.UserType, &unacceptableValueMsg, ch, &wg)
	go utils.ValidateListField(e.Address.StreetType, globals.EntityAddressStreetTypes[:], &e.Errors.StreetType, &unacceptableValueMsg, ch, &wg)

	wg.Wait()
	close(ch)

	for i := 0; i < validationsCount; i++ {
		if validationPassed := <-ch; !validationPassed {
			isValid = false
			break
		}
	}

	return isValid
}
