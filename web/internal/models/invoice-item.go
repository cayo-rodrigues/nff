package models

import (
	"log"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/sql"
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
}

type InvoiceItemFormSelectFields struct {
	Groups               *[82]string
	Origins              *[3]string
	UnitiesOfMeasurement *[23]string
}

type InvoiceItem struct {
	Id                 int
	Group              string
	Description        string
	Origin             string
	UnityOfMeasurement string
	Quantity           float64
	ValuePerUnity      float64
	InvoiceId          int
	Errors             *InvoiceItemFormError
}

func NewEmptyInvoiceItem() *InvoiceItem {
	return &InvoiceItem{
		Errors: &InvoiceItemFormError{},
	}
}

func NewInvoiceItemsFromForm(c *fiber.Ctx) (*[]InvoiceItem, error) {
	var err error

	items := []InvoiceItem{}
	postArgs := c.Request().PostArgs()

	groups := postArgs.PeekMulti("group")
	descriptions := postArgs.PeekMulti("description")
	origins := postArgs.PeekMulti("origin")
	unitiesOfMeasurement := postArgs.PeekMulti("unity_of_measurement")
	quantities := postArgs.PeekMulti("quantity")
	valuesPerUnity := postArgs.PeekMulti("value_per_unity")

	itemsQuantity := len(groups) // it could be any field
	for i := 0; i < itemsQuantity; i++ {
		item := NewEmptyInvoiceItem()

		item.Group = string(groups[i])
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

		items = append(items, *item)
	}

	return &items, nil
}

func (i *InvoiceItem) IsValid() bool {
	isValid := true
	if i.Group == "" {
		i.Errors.Group = "Campo obrigatório"
		isValid = false
	}

	if i.Description == "" {
		i.Errors.Description = "Campo obrigatório"
		isValid = false
	}

	if i.Origin == "" {
		i.Errors.Origin = "Campo obrigatório"
		isValid = false
	}

	if i.UnityOfMeasurement == "" {
		i.Errors.UnityOfMeasurement = "Campo obrigatório"
		isValid = false
	}

	if i.Quantity == 0.0 {
		i.Errors.Quantity = "Valor inaceitável"
		isValid = false
	}

	if i.ValuePerUnity == 0.0 {
		i.Errors.ValuePerUnity = "Valor inaceitável"
		isValid = false
	}

	if i.Group != "" {
		hasValidOption := false
		for _, option := range &globals.InvoiceItemGroups {
			if i.Group == option {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.Group = "Valor inaceitável"
			isValid = false
		}
	}

	if i.Origin != "" {
		hasValidOption := false
		for _, option := range &globals.InvoiceItemOrigins {
			if i.Origin == option {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.Origin = "Valor inaceitável"
			isValid = false
		}
	}

	if i.UnityOfMeasurement != "" {
		hasValidOption := false
		for _, option := range &globals.InvoiceItemUnitiesOfMeaasurement {
			if i.UnityOfMeasurement == option {
				hasValidOption = true
				break
			}
		}
		if !hasValidOption {
			i.Errors.UnityOfMeasurement = "Valor inaceitável"
			isValid = false
		}
	}

	return isValid
}

func (i *InvoiceItem) Scan(rows sql.Scanner) error {
	return rows.Scan(
		&i.Id, &i.Group, &i.Description, &i.Origin,
		&i.UnityOfMeasurement, &i.Quantity, &i.ValuePerUnity, &i.InvoiceId,
	)
}
