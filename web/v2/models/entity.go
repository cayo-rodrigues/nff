package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Entity struct {
	ID         int           `json:"-"`
	Name       string        `json:"-"`
	UserType   string        `json:"user_type"`
	Ie         string        `json:"ie"`
	CpfCnpj    string        `json:"cpf_cnpj"`
	Email      string        `json:"email"`
	Password   string        `json:"password"`
	CreatedBy  int           `json:"-"`
	CreatedAt  time.Time     `json:"-"`
	UpdatedAt  time.Time     `json:"-"`
	Errors     ErrorMessages `json:"-"`
	*Address
}

type Address struct {
	PostalCode   string `json:"postal_code"`
	Neighborhood string `json:"neighborhood"`
	StreetType   string `json:"street_type"`
	StreetName   string `json:"street_name"`
	Number       string `json:"number"`
}

func NewEntity() *Entity {
	return &Entity{
		Address: &Address{},
	}
}

func NewEntityFromForm(c *fiber.Ctx) *Entity {
	id, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		id = 0
	}
	return &Entity{
		ID:       id,
		Name:     strings.TrimSpace(c.FormValue("name")),
		UserType: strings.TrimSpace(c.FormValue("user_type")),
		CpfCnpj:  strings.TrimSpace(c.FormValue("cpf_cnpj")),
		Ie:       strings.TrimSpace(c.FormValue("ie")),
		Email:    strings.TrimSpace(c.FormValue("email")),
		Password: strings.TrimSpace(c.FormValue("password")),
		Address: &Address{
			PostalCode:   strings.TrimSpace(c.FormValue("postal_code")),
			Neighborhood: strings.TrimSpace(c.FormValue("neighborhood")),
			StreetType:   strings.TrimSpace(c.FormValue("street_type")),
			StreetName:   strings.TrimSpace(c.FormValue("street_name")),
			Number:       strings.TrimSpace(c.FormValue("number")),
		},
	}
}

func (e *Entity) Validate() bool {
	fields := Fields{
		{
			Name:  "Name",
			Value: e.Name,
			Rules: Rules(Required, Max(128)),
		},
		{
			Name:  "UserType",
			Value: e.UserType,
			Rules: Rules(Required, OneOf(EntityUserTypes)),
		},
		{
			Name:  "CpfCnpj",
			Value: e.CpfCnpj,
			Rules: Rules(Match(CPFRegex, CNPJRegex)),
		},
		{
			Name:  "Ie",
			Value: e.Ie,
			Rules: Rules(Match(IEMGRegex)),
		},
		{
			Name:  "Email",
			Value: e.Email,
			Rules: Rules(Email, Max(128)),
		},
		{
			Name:  "PostalCode",
			Value: e.PostalCode,
			Rules: Rules(Match(PostalCodeRegex)),
		},
		{
			Name:  "Neighborhood",
			Value: e.Neighborhood,
			Rules: Rules(Max(64)),
		},
		{
			Name:  "StreetType",
			Value: e.StreetType,
			Rules: Rules(OneOf(EntityAddressStreetTypes)),
		},
		{
			Name:  "StreetName",
			Value: e.StreetName,
			Rules: Rules(Max(64)),
		},
		{
			Name:  "Number",
			Value: e.Number,
			Rules: Rules(Match(AddressNumberRegex)),
		},
	}
	errors, isValid := Validate(fields)
	e.Errors = errors
	return isValid
}

func (e *Entity) Values() []any {
	return []any{
		&e.ID, &e.Name, &e.UserType, &e.CpfCnpj, &e.Ie, &e.Email, &e.Password,
		&e.PostalCode, &e.Neighborhood, &e.StreetType, &e.StreetName, &e.Number,
		&e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
	}
}
