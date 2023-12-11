package utils

import (
	"math"
	"strings"
	"sync"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
)

type Field struct {
	ErrCondition bool
	ErrField     *string
	ErrMsg       *string
}

func ValidateField(field *Field, isValid *bool, wg *sync.WaitGroup) {
	defer wg.Done()
	if field.ErrCondition {
		*field.ErrField = *field.ErrMsg
		*isValid = false
	}
}

func ValidateListField[T string | int](val T, options []T, errField *string, errMsg *string, isValid *bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var zeroVal T
	if val == zeroVal {
		*errField = *errMsg
		*isValid = false
		return
	}

	for _, option := range options {
		if val == option {
			return
		}
	}

	*errField = *errMsg
	*isValid = false
}

func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func FormatDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("2006-01-02")
}

func FormatDateAsBR(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("02/01/2006")
}

func NDaysBefore(now time.Time, days int) time.Time {
	return now.Add(-30 * 24 * time.Hour)
}

func FormatedNDaysBefore(now time.Time, days int) string {
	return FormatDate(NDaysBefore(now, 30))
}

func RoundToTwoDecimalPlaces(num float64) float64 {
	return math.Round(num*100) / 100
}

func GetReqCardErrSummary(reqMsg string) string {
	errSummary, _, _ := strings.Cut(reqMsg, "\n")
	return errSummary
}

func GetInvoiceItemSelectFields() *globals.InvoiceItemFormSelectFields {
	return &globals.InvoiceItemFormSelectFields{
		Groups:               &globals.InvoiceItemGroups,
		Origins:              &globals.InvoiceItemOrigins,
		UnitiesOfMeasurement: &globals.InvoiceItemUnitiesOfMeaasurement,
	}
}
