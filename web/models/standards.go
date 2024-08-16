package models

import (
	"regexp"
)

var IEMGRegex = regexp.MustCompile(`^\d{3}.?\d{3}.?\d{3}\/?\d{4}$`)
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

func (o *SiareInvoiceOperations) VENDA() string {
	return o[0]
}

func (o *SiareInvoiceOperations) REMESSA() string {
	return o[1]
}

type SiareInvoiceCfops struct {
	VENDA   [14]string
	REMESSA [22]string
}

func (cfops *SiareInvoiceCfops) ByOperation(invoiceOperation string) []string {
	switch invoiceOperation {
	case InvoiceOperations.VENDA():
		return cfops.VENDA[:]
	default:
		return cfops.REMESSA[:]
	}
}

var InvoiceCfops = &SiareInvoiceCfops{
	VENDA: [14]string{
		"5101 - Venda de produção do estabelecimento",
		"5102 - Venda de mercadoria adquirida ou recebida de terceiros, ou qualquer venda de mercadoria efetuada pelo MEI com exceção das saídas classificadas nos códigos 5.501, 5.502, 5.504 e 5.505.",
		"5103 - Venda de produção do estabelecimento, efetuada fora do estabelecimento",
		"5105 - Venda de produção do estabelecimento que não deva por ele transitar",
		"5111 - Venda de produção do estabelecimento remetida anteriormente em consignação industrial",
		"5113 - Venda de produção do estabelecimento remetida anteriormente em consignação mercantil",
		"5116 - Venda de produção do estabelecimento originada de encomenda para entrega futura",
		"5118 - Venda de produção do estabelecimento entregue ao destinatário por conta e ordem do adquirente originário, em venda à ordem",
		"5122 - Venda de produção do estabelecimento remetida para industrialização, por conta e ordem do adquirente, sem transitar pelo estabelecimento do adquirente",
		"5159 - Fornecimento de produção do estabelecimento de ato cooperativo",
		"5160 - Fornecimento de mercadoria adquirida ou recebida de terceiros de ato cooperativo",
		"5401 - Venda de produção do estabelecimento em operação com produto sujeito ao regime de substituição tributária, na condição de contribuinte substituto",
		"5402 - Venda de produção do estabelecimento de produto sujeito ao regime de substituição tributária, em operação entre contribuintes substitutos do mesmo produto",
		"5551 - Venda de bem do ativo imobilizado",
	},
	REMESSA: [22]string{
		"5131 - Remessa de produção do estabelecimento, com previsão de posterior ajuste ou fixação de preço, de ato cooperativo",
		"5132 - Fixação de preço de produção do estabelecimento, inclusive quando remetidas anteriormente com previsão de posterior ajuste ou fixação de preço de ato cooperativo",
		"5414 - Remessa de produção do estabelecimento para venda fora do estabelecimento em operação com produto sujeito ao regime de substituição tributária",
		"5415 - Remessa de mercadoria adquirida ou recebida de terceiros para venda fora do estabelecimento, em operação com mercadoria sujeita ao regime de substituição tributária",
		"5451 - Remessa de animal e de insumo para estabelecimento produtor",
		"5452 - Remessa de insumo - Sistema de Integração e Parceria Rural",
		"5501 - Remessa de produção do estabelecimento, com fim específico de exportação",
		"5554 - Remessa de bem do ativo imobilizado para uso fora do estabelecimento",
		"5901 - Remessa para industrialização por encomenda",
		"5904 - Remessa para venda fora do estabelecimento, ou qualquer remessa efetuada pelo MEI com exceção das classificadas nos códigos 5.502 e 5.505.",
		"5905 - Remessa para depósito fechado ou armazém geral",
		"5908 - Remessa de bem por conta de contrato de comodato",
		"5910 - Remessa em bonificação, doação ou brinde",
		"5911 - Remessa de amostra grátis",
		"5912 - Remessa de mercadoria ou bem para demonstração",
		"5914 - Remessa de mercadoria ou bem para exposição ou feira",
		"5915 - Remessa de mercadoria ou bem para conserto ou reparo",
		"5917 - Remessa de mercadoria em consignação mercantil ou industrial",
		"5920 - Remessa de vasilhame ou sacaria",
		"5923 - Remessa de mercadoria por conta e ordem de terceiros, em venda à ordem",
		"5924 - Remessa para industrialização por conta e ordem do adquirente da mercadoria, quando esta não transitar pelo estabelecimento do adquirente",
		"5934 - Remessa simbólica de mercadoria depositada em armazém geral ou depósito fechado",
	},
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

func (t *SiareInvoiceIDTypes) NFANumber() string {
	return t[0]
}

func (t *SiareInvoiceIDTypes) NFAProtocol() string {
	return t[1]
}

var InvoiceIDTypes = SiareInvoiceIDTypes{
	"Número da NFA",
	"Protocolo",
}

// INVOICE ITEMS

var InvoiceItemDefaultNCM = "94019900"

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
