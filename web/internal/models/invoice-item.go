package models

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

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

func NewInvoiceItemsFromForm(c *fiber.Ctx) ([]*InvoiceItem, error) {
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

		item.Group = string(groups[i])
		item.NCM = string(ncms[i])
		item.Description = string(descriptions[i])
		item.Origin = string(origins[i])
		item.UnityOfMeasurement = string(unitiesOfMeasurement[i])
		item.Quantity, err = strconv.ParseFloat(string(quantities[i]), 64)
		if err != nil {
			log.Printf("Error converting invoice item %d quantity from string to float64: %v", i, err)
			return nil, utils.InternalServerErr
		}
		item.ValuePerUnity, err = strconv.ParseFloat(string(valuesPerUnity[i]), 64)
		if err != nil {
			log.Printf("Error converting invoice item %d value_per_unity from string to float64: %v", i, err)
			return nil, utils.InternalServerErr
		}

		items = append(items, item)
	}

	return items, nil
}

func (i *InvoiceItem) IsValid() bool {
	isValid := true

	mandatoryFieldMsg := "Campo obrigatório"
	unacceptableValueMsg := "Valor inaceitável"
	validationsCount := 9

	var wg sync.WaitGroup
	wg.Add(validationsCount)
	ch := make(chan bool, validationsCount)

	go utils.ValidateField(i.Group == "", &i.Errors.Group, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.Description == "", &i.Errors.Description, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.Origin == "", &i.Errors.Origin, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.UnityOfMeasurement == "", &i.Errors.UnityOfMeasurement, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.Quantity == 0.0, &i.Errors.Quantity, &mandatoryFieldMsg, ch, &wg)
	go utils.ValidateField(i.ValuePerUnity == 0.0, &i.Errors.ValuePerUnity, &mandatoryFieldMsg, ch, &wg)

	go utils.ValidateListField(i.Group, globals.InvoiceItemGroups[:], &i.Errors.Group, &unacceptableValueMsg, ch, &wg)
	go utils.ValidateListField(i.Origin, globals.InvoiceItemOrigins[:], &i.Errors.Origin, &unacceptableValueMsg, ch, &wg)
	go utils.ValidateListField(i.UnityOfMeasurement, globals.InvoiceItemUnitiesOfMeaasurement[:], &i.Errors.UnityOfMeasurement, &unacceptableValueMsg, ch, &wg)

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

func (i *InvoiceItem) Scan(rows db.Scanner) error {
	return rows.Scan(
		&i.ID, &i.Group, &i.Description, &i.Origin,
		&i.UnityOfMeasurement, &i.Quantity, &i.ValuePerUnity, &i.InvoiceID,
		&i.CreatedBy, &i.CreatedAt, &i.UpdatedAt,
		&i.NCM,
	)
}
