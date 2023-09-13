package models

import (
	"net/http"

	"github.com/cayo-rodrigues/nff/web/internal/sql"
)

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
	Address    *Address
	Email      string
	Password   string
	IsSelected bool
}

func NewEntityFromForm(r *http.Request) *Entity {
	return &Entity{
		Name:     r.PostFormValue("name"),
		UserType: r.PostFormValue("user_type"),
		CpfCnpj:  r.PostFormValue("cpf_cnpj"),
		Ie:       r.PostFormValue("ie"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
		Address: &Address{
			PostalCode:   r.PostFormValue("postal_code"),
			Neighborhood: r.PostFormValue("neighborhood"),
			StreetType:   r.PostFormValue("street_type"),
			StreetName:   r.PostFormValue("street_name"),
			Number:       r.PostFormValue("number"),
		},
	}
}

func (e *Entity) Scan(rows sql.Scanner) error {
	return rows.Scan(
		&e.Id, &e.Name, &e.UserType, &e.CpfCnpj, &e.Ie, &e.Email, &e.Password,
		&e.Address.PostalCode, &e.Address.Neighborhood, &e.Address.StreetType, &e.Address.StreetName, &e.Address.Number,
	)
}
