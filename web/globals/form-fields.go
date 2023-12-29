package globals

// ENTITIES

var EntityUserTypes = [2]string{
	"Produtor Rural",
	"Apenas Destinatário",
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

var InvoiceIcmsOptions = [3]string{
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

var InvoiceIDTypes = [2]string{
	"Número da NFA",
	"Protocolo",
}

// INVOICE ITEMS

var InvoiceItemGroups = [82]string{
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
