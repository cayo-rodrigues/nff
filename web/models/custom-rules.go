package models

import "github.com/cayo-rodrigues/safe"

func CfopRule(operationType string) *safe.RuleSet {
	switch operationType {
	case InvoiceOperations.VENDA():
		return safe.OneOf(InvoiceCfops.VENDA[:])
	default:
		return safe.OneOf(InvoiceCfops.REMESSA[:])
	}
}
