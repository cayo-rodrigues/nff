package models

import "net/http"

type Address struct {
	PostalCode   string
	Neighborhood string
	StreetType   string
	StreetName   string
	Number       string
}

type Entity struct {
	Id       int
	Name     string
	UserType string
	Ie       string
	CpfCnpj  string
	Address  *Address
	Email    string
	Password string
}

func NewAddress(postalCode string, neighborhood string, streetType string, streetName string, number string) *Address {
	return &Address{
		PostalCode:   postalCode,
		Neighborhood: neighborhood,
		StreetType:   streetType,
		StreetName:   streetName,
		Number:       number,
	}
}

func NewEntity(name string, userType string, cpfCnpj string, ie string, email string, password string, address *Address) *Entity {
	return &Entity{
		Name:     name,
		UserType: userType,
		CpfCnpj:  cpfCnpj,
		Ie:       ie,
		Email:    email,
		Password: password,
		Address:  address,
	}
}

func NewEntityFromForm(r *http.Request) *Entity {
	return NewEntity(
		r.PostFormValue("name"),
		r.PostFormValue("user_type"),
		r.PostFormValue("cpf_cnpj"),
		r.PostFormValue("ie"),
		r.PostFormValue("email"),
		r.PostFormValue("password"),
		NewAddress(
			r.PostFormValue("postal_code"),
			r.PostFormValue("neighborhood"),
			r.PostFormValue("street_type"),
			r.PostFormValue("street_name"),
			r.PostFormValue("number"),
		),
	)
}
