package globals

// ENTITIES

var EntityUserTypes = [14]string{
	"Produtor Rural",
	"Inscrição Estadual",
	"Protocolo",
	"Contabilista Pessoa Física",
	"Gráfica e Outros - CNPJ",
	"Despachante Aduaneiro Pessoa Física",
	"Recinto Alfandegado Pessoa Jurídica",
	"CERM/TFRM Pessoa Física",
	"CERM/TFRM Pessoa Jurídica",
	"VAF Especial",
	"Contribuinte Interestadual",
	"Pessoa Física Autuada - PTA eletrônico",
	"Responsável Tributário - Instituição Financeira",
	"Conselheiro",
}
var EntityAddressStreetTypes = [3]string{
	"Rua",
	"Estrada",
	"Avenida",
}

// INVOICE

var InvoiceOperations = [2]string{
	"VENDA",
	"REMESSA",
}

var InvoiceCfops = [14]int{
	5102,
	5101,
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

var InvoiceIcmsOptions = [3]string{
	"Sim",
	"Não",
	"Isento",
}

var InvoiceBooleanField = [2]string{
	"Sim",
	"Nâo",
}

var InvoiceItemGroups = [2]string{
	"Bovino",
	"Asinino",
}

var InvoiceItemOrigins = [3]string{
	"Nacional",
	"Estrangeira - Importação direta",
	"Estrangeira - Adquirida no mercado interno",
}

var InvoiceItemUnitiesOfMeaasurement = [23]string{
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
