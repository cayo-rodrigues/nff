package utils

import (
	"strings"
	"sync"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
)

func ValidateField(errCondition bool, errField *string, errMsg *string, result chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	isValid := true
	if errCondition {
		*errField = *errMsg
		isValid = false
	}
	result <- isValid
}

func ValidateListField[T string | int](val T, options []T, errField *string, errMsg *string, result chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var zeroVal T
	if val == zeroVal {
		*errField = *errMsg
		result <- false
		return
	}
	isValid := false
	for _, option := range options {
		if val == option {
			isValid = true
			break
		}
	}
	if !isValid {
		*errField = *errMsg
	}
	result <- isValid
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
