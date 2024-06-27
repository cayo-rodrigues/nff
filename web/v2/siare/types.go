package siare

import "github.com/cayo-rodrigues/nff/web/models"

// INVOICE ISSUING

type SSApiInvoiceRequest struct {
	*models.Invoice
	ShouldDownload bool `json:"should_download"`
}

type SSApiInvoiceResponse struct {
	Msg                string              `json:"msg"`
	IsAwaitingAnalisys bool                `json:"is_awaiting_analisys"`
	InvoiceNumber      string              `json:"invoice_number"`
	InvoiceProtocol    string              `json:"invoice_protocol"`
	InvoicePDF         string              `json:"invoice_pdf"`
	FileName           string              `json:"file_name"`
	Status             string              `json:"status"`
	Errors             *SSApiErrorResponse `json:"errors"`
}

// INVOICE CANCELING

type SSApiCancelingRequest struct {
	*models.InvoiceCancel
}

type SSApiCancelingResponse struct {
	Msg    string              `json:"msg"`
	Status string              `json:"status"`
	Errors *SSApiErrorResponse `json:"errors"`
}

// INVOICE PRINTING

type SSApiPrintingRequest struct {
	*models.InvoicePrint
}

type SSApiPrintingResponse struct {
	Msg           string              `json:"msg"`
	InvoiceNumber string              `json:"invoice_id"` // MUDAR NA SS-API PARA invoice_number
	InvoicePDF    string              `json:"invoice_pdf"`
	FileName      string              `json:"file_name"`
	Status        string              `json:"status"`
	Errors        *SSApiErrorResponse `json:"errors"`
}

// METRICS

type SSApiMetricsRequest struct {
	Body  *SSApiMetricsRequestBody
	Query *SSApiMetricsRequestQuery
}

type SSApiMetricsRequestQuery struct {
	StartDate      string
	EndDate        string
	IncludeRecords bool
}

type SSApiMetricsRequestBody struct {
	*models.Entity `json:"entity"`
}

type SSApiMetricsResponse struct {
	*models.MetricsResult
	Errors *SSApiErrorResponse `json:"errors"`
}

// ERROR RESPONSES

type SSApiErrorResponse struct {
	MissingFields []string           `json:"missing_fields"`
	InvalidFields []string           `json:"invalid_fields"`
	Entity        *SSApiEntityErrors `json:"entity"`
	Sender        *SSApiEntityErrors `json:"sender"`
	Recipient     *SSApiEntityErrors `json:"recipient"`
}

type SSApiEntityErrors struct {
	MissingFields []string `json:"missing_fields"`
	InvalidFields []string `json:"invalid_fields"`
}
