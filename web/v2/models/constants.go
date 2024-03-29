package models

import "regexp"

var EmailRegex = regexp.MustCompile(`[^@ \t\r\n]+@[^@ \t\r\n]+\.[^@ \t\r\n]+`)
var PhoneRegex = regexp.MustCompile(`(?:(?:\+|00)?(55)\s?)?(?:\(?([1-9][0-9])\)?\s?)(?:((?:9\d|[2-9])\d{3})\-?(\d{4}))`)
var WhateverRegex = regexp.MustCompile(`.*`)
var IEMGRegex = regexp.MustCompile(`^\d{3}.?\d{3}.?\d{3}\/?\d{4}$`)
var CPFRegex = regexp.MustCompile(`^\d{3}.?\d{3}.?\d{3}\-?\d{2}$`)
var CNPJRegex = regexp.MustCompile(`^(\d{2}.?\d{3}.?\d{3}\/?\d{4}\-?\d{2})$`)
var PostalCodeRegex = regexp.MustCompile(`(^\d{5})\-?(\d{3}$)`)
var AddressNumberRegex = regexp.MustCompile(`^(?:s\/n|S\/n|S\/N|s\/N)|^(\d)*$`)
var SiareNFANumberRegex = regexp.MustCompile(`^(\d{3}\.\d{3}\.\d{3}|\d{9})$`)
var SiareNFAProtocolRegex = regexp.MustCompile(`^(\d{13})$`)
var GTARegex = regexp.MustCompile(`^([a-zA-Z]-\d{1,6}(;\s*[a-zA-Z]-\d{1,6})*$)`)

// ENTITIES

type SiareUserTypes [3]string

var EntityUserTypes = SiareUserTypes{
	"Produtor Rural",
	"Inscrição Estadual",
	"Apenas Destinatário",
}

type SiareAddressStreetTypes [3]string

var EntityAddressStreetTypes = SiareAddressStreetTypes{
	"Rua",
	"Estrada",
	"Avenida",
}

// INVOICE

type SiareInvoiceOperations [2]string

var InvoiceOperations = SiareInvoiceOperations{
	"VENDA",
	"REMESSA",
}

type SiareInvoiceCfops [14]int

var InvoiceCfops = SiareInvoiceCfops{
	5101,
	5102,
	5103,
	5105,
	5111,
	5113,
	5116,
	5118,
	5122,
	5159,
	5160,
	5401,
	5402,
	5551,
}

type SiareInvoiceIcmsOptions [3]string

var InvoiceIcmsOptions = SiareInvoiceIcmsOptions{
	"Sim",
	"Não",
	"Isento",
}

type BooleanField [2]string

func (f *BooleanField) Reverse() *BooleanField {
	return &BooleanField{f[1], f[0]}
}

var InvoiceBooleanField = BooleanField{
	"Sim",
	"Não",
}

type SiareInvoiceIDTypes [2]string

var InvoiceIDTypes = SiareInvoiceIDTypes{
	"Número da NFA",
	"Protocolo",
}

// INVOICE ITEMS

type SiareInvoiceItemGroups [82]string

var InvoiceItemGroups = SiareInvoiceItemGroups{
	"Adubo",
	"Algodão",
	"Animais silvestres",
	"Apicultura",
	"Aquicultura e pesca",
	"Avicultura - ovos",
	"Avicultura - reprodutor",
	"Avicultura para corte",
	"Avicultura para recria",
	"Café",
	"Cana de açucar",
	"Carvão de floresta nativa",
	"Carvão de floresta plantada",
	"Carvão mineral",
	"Cereais",
	"Combustíveis",
	"Derivados do leite",
	"Dormentes",
	"Embrião asinino",
	"Embrião bovino",
	"Embrião bufalino",
	"Embrião caprino",
	"Embrião equino",
	"Embrião muar",
	"Embrião ovino",
	"Embrião suíno",
	"Embrião taurino",
	"Esterco animal",
	"Farelos",
	"Feno",
	"Flores",
	"Gado asinino - reprodutor",
	"Gado asinino para corte",
	"Gado asinino para recria",
	"Gado asinino para serviço",
	"Gado bovino - reprodutor",
	"Gado bovino para corte",
	"Gado bovino para recria",
	"Gado bovino para serviço",
	"Gado bufalino - reprodutor",
	"Gado bufalino para corte",
	"Gado bufalino para recria",
	"Gado bufalino para serviço",
	"Gado caprinos vivos",
	"Gado equino - reprodutor",
	"Gado equino para corte",
	"Gado equino para recria",
	"Gado equino para serviço",
	"Gado muar - reprodutor",
	"Gado muar para corte",
	"Gado muar para recria",
	"Gado muar para serviço",
	"Gado suino - reprodutor",
	"Gado suino para corte",
	"Gado suino para recria",
	"Gado taurino - reprodutor",
	"Gado taurino para corte",
	"Gado taurino para recria",
	"Gado taurino para serviço",
	"Hortifrutigranjeiros",
	"Leite",
	"Lenha - floresta nativa",
	"Lenha - floresta plantada",
	"Madeira",
	"Milho",
	"Minerais",
	"Mudança",
	"Mudas e sementes",
	"Outros",
	"Ovinocultura",
	"Palha de café",
	"Piscicultura - reprodutor",
	"Piscicultura para corte",
	"Piscicultura para recria",
	"Prestação de serviço de transporte",
	"Queijo artesanal",
	"Ração",
	"Resíduos",
	"Semem",
	"Soja",
	"Sorgos",
	"Suplementos",
}

type SiareInvoiceItemOrigins [3]string

var InvoiceItemOrigins = SiareInvoiceItemOrigins{
	"Nacional",
	"Estrangeira - Importação direta",
	"Estrangeira - Adquirida no mercado interno",
}

type SiareUnitiesOfMeasurement [23]string

var InvoiceItemUnitiesOfMeaasurement = SiareUnitiesOfMeasurement{
	"CB",
	"CT",
	"CX",
	"DUZIA",
	"EST",
	"G",
	"JOGO",
	"KG",
	"KM",
	"LT",
	"M2",
	"M3",
	"MDC",
	"METRO",
	"MI",
	"PARES",
	"PC",
	"PT",
	"QUILAT",
	"SC",
	"ST",
	"TON",
	"UN",
}
