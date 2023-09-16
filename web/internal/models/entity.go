package models

import (
	"net/http"
	"regexp"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/sql"
)

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

func NewEntityFromForm(r *http.Request) *Entity {
	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil {
		id = 0
	}
	return &Entity{
		Id:         id,
		Name:       r.PostFormValue("name"),
		UserType:   r.PostFormValue("user_type"),
		CpfCnpj:    r.PostFormValue("cpf_cnpj"),
		Ie:         r.PostFormValue("ie"),
		Email:      r.PostFormValue("email"),
		Password:   r.PostFormValue("password"),
		Address: &Address{
			PostalCode:   r.PostFormValue("postal_code"),
			Neighborhood: r.PostFormValue("neighborhood"),
			StreetType:   r.PostFormValue("street_type"),
			StreetName:   r.PostFormValue("street_name"),
			Number:       r.PostFormValue("number"),
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
	if e.Name == "" {
		e.Errors.Name = "Campo obrigatório"
		isValid = false
	}

	if e.Ie == "" && (e.Address.PostalCode == "" || e.Address.Neighborhood == "" || e.Address.StreetType == "" || e.Address.StreetName == "" || e.Address.Number == "") {
		e.Errors.Ie = "Ie OU endereço completo obrigatórios"
		isValid = false
	}

	reIEMG := regexp.MustCompile(`^\d{3}.?\d{3}.?\d{3}\/?\d{4}$`)
	if e.Ie != "" && !reIEMG.MatchString(e.Ie) {
		e.Errors.Ie = "Formato inválido"
		isValid = false
	}

	reCpf := regexp.MustCompile(`^\d{3}.?\d{3}.?\d{3}\-?\d{2}$`)
	reCnpj := regexp.MustCompile(`^(\d{2}.?\d{3}.?\d{3}\/?\d{4}\-?\d{2})$`)
	if e.CpfCnpj != "" && !reCpf.MatchString(e.CpfCnpj) && !reCnpj.MatchString(e.CpfCnpj) {
		e.Errors.CpfCnpj = "Formato inválido"
		isValid = false
	}

	reEmail := regexp.MustCompile(`[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+`)
	if e.Email != "" && !reEmail.MatchString(e.Email) {
		e.Errors.Email = "Formato inválido"
		isValid = false
	}

	rePostalCode := regexp.MustCompile(`(^\d{5})\-?(\d{3}$)`)
	if e.Address.PostalCode != "" && !rePostalCode.MatchString(e.Address.PostalCode) {
		e.Errors.PostalCode = "Formato inválido"
		isValid = false
	}

	if e.UserType != "" && e.UserType != "Produtor Rural" {
		e.Errors.UserType = "Valor inaceitável"
		isValid = false
	}

	if e.Address.StreetType != "" && e.Address.StreetType != "Rua" {
		e.Errors.StreetType = "Valor inaceitável"
		isValid = false
	}

	return isValid
}
