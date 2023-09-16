package models

type InvoiceItem struct {
	Group              string
	Description        string
	Origin             string
	UnityOfMeasurement string
	Quantity           int
	ValuePerUnity      float64
}

type Invoice struct {
	Id                 int
	Number             string
	Protocol           string
	Operation          string
	Cfop               string
	IsFinalCustomer    string
	IsIcmsContributor  string
	Shipping           float64
	AddShippingToTotal string
	Gta                string
	Sender             *Entity
	Recipient          *Entity
	Items              *[]InvoiceItem
}
