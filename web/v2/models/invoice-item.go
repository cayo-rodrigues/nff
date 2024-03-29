package models

import (
	"log"
	"time"

	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type InvoiceItem struct {
	ID                 int           `json:"-"`
	Group              string        `json:"group"`
	Description        string        `json:"description"`
	Origin             string        `json:"origin"`
	UnityOfMeasurement string        `json:"unity_of_measurement"`
	Quantity           float64       `json:"quantity"`
	ValuePerUnity      float64       `json:"value_per_unity"`
	InvoiceID          int           `json:"-"`
	CreatedBy          int           `json:"-"`
	CreatedAt          time.Time     `json:"-"`
	UpdatedAt          time.Time     `json:"-"`
	NCM                string        `json:"ncm"`
	Errors             ErrorMessages `json:"-"`
}

func NewInvoiceItem() *InvoiceItem {
	return &InvoiceItem{}
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
		item := NewInvoiceItem()

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
	fields := Fields{
		{
			Name:  "Group",
			Value: i.Group,
			Rules: Rules(OneOf(InvoiceItemGroups)),
		},
		{
			Name:  "NCM",
			Value: i.NCM,
			Rules: Rules(Max(16)),
		},
		{
			Name:  "Description",
			Value: i.Description,
			Rules: Rules(Required, Max(128)),
		},
		{
			Name:  "Origin",
			Value: i.Origin,
			Rules: Rules(OneOf(InvoiceItemOrigins)),
		},
		{
			Name:  "UnityOfMeasurement",
			Value: i.UnityOfMeasurement,
			Rules: Rules(OneOf(InvoiceItemUnitiesOfMeaasurement)),
		},
		{
			Name:  "Quantity",
			Value: i.Quantity,
			Rules: Rules(Required),
		},
		{
			Name:  "ValuePerUnity",
			Value: i.ValuePerUnity,
			Rules: Rules(Required),
		},
	}
	errors, isValid := Validate(fields)
	i.Errors = errors
	return isValid
}

func (i *InvoiceItem) Values() []any {
	return []any{
		&i.ID, &i.Group, &i.Description, &i.Origin,
		&i.UnityOfMeasurement, &i.Quantity, &i.ValuePerUnity, &i.InvoiceID,
		&i.CreatedBy, &i.CreatedAt, &i.UpdatedAt,
		&i.NCM,
	}
}
