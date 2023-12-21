package models

import (
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type EntityFormSelectFields struct {
	UserTypes   *[2]string
	StreetTypes *[3]string
}

type EntityFormError struct {
	Name         string
	UserType     string
	CpfCnpj      string
	Ie           string
	Email        string
	PostalCode   string
	Neighborhood string
	StreetType   string
	StreetName   string
	Number       string
}

type Address struct {
	PostalCode   string `json:"postal_code"`
	Neighborhood string `json:"neighborhood"`
	StreetType   string `json:"street_type"`
	StreetName   string `json:"street_name"`
	Number       string `json:"number"`
}

type Entity struct {
	ID         int              `json:"-"`
	Name       string           `json:"-"`
	UserType   string           `json:"user_type"`
	Ie         string           `json:"ie"`
	CpfCnpj    string           `json:"cpf_cnpj"`
	Email      string           `json:"email"`
	Password   string           `json:"password"`
	IsSelected bool             `json:"-"`
	CreatedBy  int              `json:"-"`
	CreatedAt  time.Time        `json:"-"`
	UpdatedAt  time.Time        `json:"-"`
	Errors     *EntityFormError `json:"-"`
	*Address
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
		ID:       id,
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

func (e *Entity) Scan(rows db.Scanner) error {
	return rows.Scan(
		&e.ID, &e.Name, &e.UserType, &e.CpfCnpj, &e.Ie, &e.Email, &e.Password,
		&e.Address.PostalCode, &e.Address.Neighborhood, &e.Address.StreetType, &e.Address.StreetName, &e.Address.Number,
		&e.CreatedBy, &e.CreatedAt, &e.UpdatedAt,
	)
}

func (e *Entity) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := globals.MandatoryFieldMsg
	invalidFormatMsg := globals.InvalidFormatMsg
	mustHaveIeOrAddressMsg := globals.MustHaveIeOrAddressMsg
	unacceptableValueMsg := globals.UnacceptableValueMsg
	valueTooLongMsg := globals.ValueTooLongMsg

	hasName := e.Name != ""
	hasIe := e.Ie != ""
	hasAddress := e.Address.PostalCode != "" && e.Address.Neighborhood != "" && e.Address.StreetType != "" && e.Address.StreetName != "" && e.Address.Number != ""
	hasCpfCnpj := e.CpfCnpj != ""
	hasEmail := e.Email != ""
	hasPostalCode := e.Address.PostalCode != ""
	hasNumber := e.Address.Number != ""

	hasValidIeFormat := globals.ReIeMg.MatchString(e.Ie)
	hasValidCpfCnpjFormat := globals.ReCpf.MatchString(e.CpfCnpj) && globals.ReCnpj.MatchString(e.CpfCnpj)
	hasValidEmailFormat := globals.ReEmail.MatchString(e.Email)
	hasValidPostalCodeFormat := globals.RePostalCode.MatchString(e.Address.PostalCode)
	hasValidNumberFormat := globals.ReAddressNumber.MatchString(e.Address.Number)

	nameTooLong := utf8.RuneCount([]byte(e.Name)) > 128
	neighborhoodTooLong := utf8.RuneCount([]byte(e.Address.Neighborhood)) > 64
	streetNameTooLong := utf8.RuneCount([]byte(e.Address.StreetName)) > 64
	numberTooLong := utf8.RuneCount([]byte(e.Address.Number)) > 6

	fields := [11]*utils.Field{
		{ErrCondition: !hasName, ErrField: &e.Errors.Name, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasIe && !hasAddress, ErrField: &e.Errors.Ie, ErrMsg: &mustHaveIeOrAddressMsg},
		{ErrCondition: hasIe && !hasValidIeFormat, ErrField: &e.Errors.Ie, ErrMsg: &invalidFormatMsg},
		{ErrCondition: hasCpfCnpj && !hasValidCpfCnpjFormat, ErrField: &e.Errors.CpfCnpj, ErrMsg: &invalidFormatMsg},
		{ErrCondition: hasEmail && !hasValidEmailFormat, ErrField: &e.Errors.Email, ErrMsg: &invalidFormatMsg},
		{ErrCondition: hasNumber && !hasValidNumberFormat, ErrField: &e.Errors.Number, ErrMsg: &invalidFormatMsg},
		{ErrCondition: hasPostalCode && !hasValidPostalCodeFormat, ErrField: &e.Errors.PostalCode, ErrMsg: &invalidFormatMsg},
		{ErrCondition: nameTooLong, ErrField: &e.Errors.Name, ErrMsg: &valueTooLongMsg},
		{ErrCondition: neighborhoodTooLong, ErrField: &e.Errors.Neighborhood, ErrMsg: &valueTooLongMsg},
		{ErrCondition: streetNameTooLong, ErrField: &e.Errors.StreetName, ErrMsg: &valueTooLongMsg},
		{ErrCondition: hasNumber && numberTooLong, ErrField: &e.Errors.Number, ErrMsg: &valueTooLongMsg},
	}

	var wg sync.WaitGroup
	for _, field := range fields {
		wg.Add(1)
		go utils.ValidateField(field, &isValid, &wg)
	}

	wg.Add(2)
	go utils.ValidateListField(e.UserType, globals.EntityUserTypes[:], &e.Errors.UserType, &unacceptableValueMsg, &isValid, &wg)
	go utils.ValidateListField(e.Address.StreetType, globals.EntityAddressStreetTypes[:], &e.Errors.StreetType, &unacceptableValueMsg, &isValid, &wg)

	wg.Wait()

	return isValid
}
