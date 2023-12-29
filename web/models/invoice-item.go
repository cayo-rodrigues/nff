package models

import (
	"log"
	"time"
	"unicode/utf8"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/globals"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoiceItemctFields struct {
	Groups               *[82]string
	Origins              *[3]string
	UnitiesOfMeasurement *[23]string
}

func NewInvoiceItemSelectFields() *InvoiceItemctFields {
	return &InvoiceItemctFields{
		Groups:               &globals.InvoiceItemGroups,
		Origins:              &globals.InvoiceItemOrigins,
		UnitiesOfMeasurement: &globals.InvoiceItemUnitiesOfMeaasurement,
	}
}

type InvoiceItemFormError struct {
	Group              string
	Description        string
	Origin             string
	UnityOfMeasurement string
	Quantity           string
	ValuePerUnity      string
	NCM                string
}

type InvoiceItem struct {
	ID                 int                   `json:"-"`
	Group              string                `json:"group"`
	Description        string                `json:"description"`
	Origin             string                `json:"origin"`
	UnityOfMeasurement string                `json:"unity_of_measurement"`
	Quantity           float64               `json:"quantity"`
	ValuePerUnity      float64               `json:"value_per_unity"`
	InvoiceID          int                   `json:"-"`
	Errors             *InvoiceItemFormError `json:"-"`
	CreatedBy          int                   `json:"-"`
	CreatedAt          time.Time             `json:"-"`
	UpdatedAt          time.Time             `json:"-"`
	NCM                string                `json:"ncm"`
}

func NewEmptyInvoiceItem() *InvoiceItem {
	return &InvoiceItem{
		Errors: &InvoiceItemFormError{},
	}
}

func NewInvoiceItemsFromForm(c *fiber.Ctx) []*InvoiceItem {
	var err error

	items := []*InvoiceItem{}
	postArgs := c.Request().PostArgs()

	groups := postArgs.PeekMulti("group")
	ncms := postArgs.PeekMulti("ncm")
	descriptions := postArgs.PeekMulti("description")
	origins := postArgs.PeekMulti("origin")
	unitiesOfMeasurement := postArgs.PeekMulti("unity_of_measurement")
	quantities := postArgs.PeekMulti("quantity")
	valuesPerUnity := postArgs.PeekMulti("value_per_unity")

	itemsQuantity := len(groups) // it could be any field
	for i := 0; i < itemsQuantity; i++ {
		item := NewEmptyInvoiceItem()

		item.Group = utils.TrimSpaceBytes(groups[i])
		item.NCM = utils.TrimSpaceBytes(ncms[i])
		item.Description = utils.TrimSpaceBytes(descriptions[i])
		item.Origin = utils.TrimSpaceBytes(origins[i])
		item.UnityOfMeasurement = utils.TrimSpaceBytes(unitiesOfMeasurement[i])
		item.Quantity, err = utils.TrimSpaceFromBytesToFloat64(quantities[i])
		if err != nil {
			log.Printf("Error converting invoice item %d quantity from string to float64: %v", i, err)
		}
		item.ValuePerUnity, err = utils.TrimSpaceFromBytesToFloat64(valuesPerUnity[i])
		if err != nil {
			log.Printf("Error converting invoice item %d value_per_unity from string to float64: %v", i, err)
		}

		items = append(items, item)
	}

	return items
}

func (i *InvoiceItem) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := globals.MandatoryFieldMsg
	unacceptableValueMsg := globals.UnacceptableValueMsg
	valueTooLongMsg := globals.ValueTooLongMsg

	hasDescription := i.Description != ""
	hasQuantity := i.Quantity != 0.0
	hasValuePerUnity := i.ValuePerUnity != 0.0
	hasNCM := i.NCM != ""

	descriptionTooLong := utf8.RuneCount([]byte(i.Description)) > 128
	ncmTooLong := utf8.RuneCount([]byte(i.NCM)) > 16

	fields := [5]*utils.Field{
		{ErrCondition: !hasDescription, ErrField: &i.Errors.Description, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasQuantity, ErrField: &i.Errors.Quantity, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: !hasValuePerUnity, ErrField: &i.Errors.ValuePerUnity, ErrMsg: &mandatoryFieldMsg},
		{ErrCondition: descriptionTooLong, ErrField: &i.Errors.Description, ErrMsg: &valueTooLongMsg},
		{ErrCondition: hasNCM && ncmTooLong, ErrField: &i.Errors.NCM, ErrMsg: &valueTooLongMsg},
	}

	for _, field := range fields {
		utils.ValidateField(field, &isValid)
	}

	utils.ValidateListField(i.Group, globals.InvoiceItemGroups[:], &i.Errors.Group, &unacceptableValueMsg, &isValid)
	utils.ValidateListField(i.Origin, globals.InvoiceItemOrigins[:], &i.Errors.Origin, &unacceptableValueMsg, &isValid)
	utils.ValidateListField(i.UnityOfMeasurement, globals.InvoiceItemUnitiesOfMeaasurement[:], &i.Errors.UnityOfMeasurement, &unacceptableValueMsg, &isValid)

	return isValid
}

func (i *InvoiceItem) Scan(rows db.Scanner) error {
	return rows.Scan(
		&i.ID, &i.Group, &i.Description, &i.Origin,
		&i.UnityOfMeasurement, &i.Quantity, &i.ValuePerUnity, &i.InvoiceID,
		&i.CreatedBy, &i.CreatedAt, &i.UpdatedAt,
		&i.NCM,
	)
}
