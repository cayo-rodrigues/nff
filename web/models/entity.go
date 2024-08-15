package models

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cayo-rodrigues/safe"
	"github.com/gofiber/fiber/v2"
)

type Entity struct {
	ID        int                `json:"-"`
	Name      string             `json:"-"`
	UserType  string             `json:"user_type"`
	Ie        string             `json:"ie"`
	OtherIes  []string           `json:"other_ies"`
	CpfCnpj   string             `json:"cpf_cnpj"`
	Email     string             `json:"email"`
	Password  string             `json:"password"`
	CreatedBy int                `json:"-"`
	CreatedAt time.Time          `json:"-"`
	UpdatedAt time.Time          `json:"-"`
	Errors    safe.ErrorMessages `json:"-"`
	*Address
}

type Address struct {
	PostalCode   string `json:"postal_code"`
	Neighborhood string `json:"neighborhood"`
	StreetType   string `json:"street_type"`
	StreetName   string `json:"street_name"`
	Number       string `json:"number"`
}

func (a *Address) Values() []any {
	return []any{a.PostalCode, a.Neighborhood, a.StreetType, a.StreetName, a.Number}
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


	entity := &Entity{
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

	re := regexp.MustCompile(`[ ./]`)
	entity.Ie = re.ReplaceAllString(entity.Ie, "")

	ies := c.FormValue("other_ies")
	if ies != "" {
		for _, ie := range strings.Split(ies, ",") {
			ie = re.ReplaceAllString(ie, "")
			entity.OtherIes = append(entity.OtherIes, ie)
		}
	}

	return entity
}

func (e *Entity) IsValid() bool {
	fields := safe.Fields{
		{
			Name:  "Name",
			Value: e.Name,
			Rules: safe.Rules{safe.Required(), safe.Max(128)},
		},
		{
			Name:  "UserType",
			Value: e.UserType,
			Rules: safe.Rules{safe.Required(), safe.OneOf(EntityUserTypes[:])},
		},
		{
			Name:  "CpfCnpj",
			Value: e.CpfCnpj,
			Rules: safe.Rules{safe.CpfCnpj()},
		},
		{
			Name:  "Ie",
			Value: e.Ie,
			Rules: safe.Rules{
				safe.Match(IEMGRegex),
				safe.RequiredUnless(safe.All(e.Address.Values()...)),
			},
		},
		{
			Name:  "Email",
			Value: e.Email,
			Rules: safe.Rules{safe.Email(), safe.Max(128)},
		},
		{
			Name:  "PostalCode",
			Value: e.PostalCode,
			Rules: safe.Rules{safe.Match(safe.CepRegex)},
		},
		{
			Name:  "Neighborhood",
			Value: e.Neighborhood,
			Rules: safe.Rules{safe.Max(64)},
		},
		{
			Name:  "StreetType",
			Value: e.StreetType,
			Rules: safe.Rules{safe.OneOf(EntityAddressStreetTypes[:])},
		},
		{
			Name:  "StreetName",
			Value: e.StreetName,
			Rules: safe.Rules{safe.Max(64)},
		},
		{
			Name:  "Number",
			Value: e.Number,
			Rules: safe.Rules{safe.Match(safe.AddressNumberRegex)},
		},
		{
			Name:  "OtherIes",
			Value: e.OtherIes,
			Rules: safe.Rules{safe.MatchList(IEMGRegex), safe.UniqueList[string]()},
		},
	}
	errors, isValid := safe.Validate(fields)
	e.Errors = errors
	return isValid
}

func (e *Entity) Values() []any {
	return []any{
		&e.ID, &e.Name, &e.UserType, &e.CpfCnpj, &e.Ie, &e.Email, &e.Password,
		&e.PostalCode, &e.Neighborhood, &e.StreetType, &e.StreetName, &e.Number,
		&e.CreatedBy, &e.CreatedAt, &e.UpdatedAt, &e.OtherIes,
	}
}

func (e *Entity) AllIes() []string {
	availableIes := []string{e.Ie}
	for _, ie := range e.OtherIes {
		availableIes = append(availableIes, ie)
	}
	return availableIes
}
