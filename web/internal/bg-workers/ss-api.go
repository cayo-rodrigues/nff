package bgworkers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/interfaces"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

var SS_API_BASE_URL = ""

type SSAPIEntityErrors struct {
	MissingFields []string `json:"missing_fields"`
	InvalidFields []string `json:"invalid_fields"`
}

type SSAPIInvoiceErrors struct {
	MissingFields []string           `json:"missing_fields"`
	InvalidFields []string           `json:"invalid_fields"`
	Sender        *SSAPIEntityErrors `json:"sender"`
	Recipient     *SSAPIEntityErrors `json:"recipient"`
}

type SSAPIInvoiceRequest struct {
	*models.Invoice
	ShouldDownload bool `json:"should_download"`
}

type SSAPIInvoiceResponse struct {
	Msg                string              `json:"msg"`
	IsAwaitingAnalisys bool                `json:"is_awaiting_analisys"`
	InvoiceNumber      string              `json:"invoice_id"`
	InvoiceProtocol    string              `json:"invoice_protocol"`
	InvoicePDF         string              `json:"invoice_pdf"`
	FileName           string              `json:"file_name"`
	Status             string              `json:"status"`
	Errors             *SSAPIInvoiceErrors `json:"errors"`
}

type SSAPICancelingErrors struct {
	MissingFields []string           `json:"missing_fields"`
	InvalidFields []string           `json:"invalid_fields"`
	Entity        *SSAPIEntityErrors `json:"entity"`
}

type SSAPICancelingRequest struct {
	*models.InvoiceCancel
}

type SSAPICancelingResponse struct {
	Msg    string                `json:"msg"`
	Status string                `json:"status"`
	Errors *SSAPICancelingErrors `json:"errors"`
}

type SSAPIMetricsErrors struct {
	MissingFields []string           `json:"missing_fields"`
	InvalidFields []string           `json:"invalid_fields"`
	Entity        *SSAPIEntityErrors `json:"entity"`
}

type SSAPIMetricsRequest struct {
	Body  *SSAPIMetricsRequestBody
	Query *SSAPIMetricsRequestQuery
}

type SSAPIMetricsRequestQuery struct {
	StartDate string
	EndDate   string
}

type SSAPIMetricsRequestBody struct {
	*models.Entity `json:"entity"`
}

type SSAPIMetricsResponse struct {
	*models.MetricsResult
	Errors *SSAPIMetricsErrors `json:"errors"`
}

type SSAPIPrintingErrors struct {
	MissingFields []string           `json:"missing_fields"`
	InvalidFields []string           `json:"invalid_fields"`
	Entity        *SSAPIEntityErrors `json:"entity"`
}

type SSAPIPrintingRequest struct {
	*models.InvoicePrint
}

type SSAPIPrintingResponse struct {
	Msg        string               `json:"msg"`
	InvoiceId  string               `json:"invoice_id"`
	InvoicePDF string               `json:"invoice_pdf"`
	FileName   string               `json:"file_name"`
	Status     string               `json:"status"`
	Errors     *SSAPIPrintingErrors `json:"errors"`
}

type SiareBGWorker struct {
	invoiceService   interfaces.InvoiceService
	cancelingService interfaces.CancelingService
	metricsService   interfaces.MetricsService
	printingService  interfaces.PrintingService
}

func NewSiareBGWorker(
	invoiceService interfaces.InvoiceService,
	cancelingService interfaces.CancelingService,
	metricsService interfaces.MetricsService,
	printingService interfaces.PrintingService,
) *SiareBGWorker {
	return &SiareBGWorker{
		invoiceService:   invoiceService,
		cancelingService: cancelingService,
		metricsService:   metricsService,
		printingService:  printingService,
	}
}

func (w *SiareBGWorker) RequestInvoice(invoice *models.Invoice) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	reqBody := SSAPIInvoiceRequest{
		Invoice:        invoice,
		ShouldDownload: true,
	}

	agent := fiber.Post(SS_API_BASE_URL + "/invoice/request")
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			log.Printf("Something went wrong with the request at /invoice/request for invoice with id %v: %v\n", invoice.ID, err)
		}
	}

	var response SSAPIInvoiceResponse
	json.Unmarshal(body, &response)

	if response.Errors != nil {
		response.Msg = fmt.Sprintf("%s\n%s", response.Msg, formatInvoiceErrs(&response))
	}

	invoice.Number = response.InvoiceNumber
	invoice.Protocol = response.InvoiceProtocol
	invoice.ReqMsg = response.Msg
	invoice.ReqStatus = response.Status
	invoice.PDF = response.InvoicePDF
	if response.FileName != "" {
		invoice.CustomFileName = response.FileName
	}

	err := w.invoiceService.UpdateInvoice(ctx, invoice)
	if err != nil {
		log.Printf("Something went wrong when updating invoice history. Invoice with id %v will be on 'pending' state for ever: %v\n", invoice.ID, err)
	}

	key := fmt.Sprintf("reqstatus:invoice:%v", invoice.ID)
	db.Redis.Set(ctx, key, true, time.Minute)
}

func (w *SiareBGWorker) RequestInvoiceCanceling(invoiceCancel *models.InvoiceCancel) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	reqBody := SSAPICancelingRequest{
		InvoiceCancel: invoiceCancel,
	}

	agent := fiber.Post(SS_API_BASE_URL + "/invoice/cancel")
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			log.Printf("Something went wrong with the request at /invoice/cancel for canceling with id %v: %v\n", invoiceCancel.ID, err)
		}
	}

	var response SSAPICancelingResponse
	json.Unmarshal(body, &response)

	if response.Errors != nil {
		response.Msg = fmt.Sprintf("%s\n%s", response.Msg, formatCancelingErrs(&response))
	}

	invoiceCancel.ReqStatus = response.Status
	invoiceCancel.ReqMsg = response.Msg

	err := w.cancelingService.UpdateInvoiceCanceling(ctx, invoiceCancel)
	if err != nil {
		log.Printf("Something went wrong when updating invoice canceling history. Canceling with id %v will be on 'pending' state for ever: %v\n", invoiceCancel.ID, err)
	}

	key := fmt.Sprintf("reqstatus:canceling:%v", invoiceCancel.ID)
	db.Redis.Set(ctx, key, true, time.Minute)
}

func (w *SiareBGWorker) GetMetrics(query *models.MetricsQuery) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	reqData := SSAPIMetricsRequest{
		Body: &SSAPIMetricsRequestBody{
			Entity: query.Entity,
		},
		Query: &SSAPIMetricsRequestQuery{
			StartDate: utils.FormatDateAsBR(query.StartDate),
			EndDate:   utils.FormatDateAsBR(query.EndDate),
		},
	}

	agent := fiber.Get(SS_API_BASE_URL + "/invoice/overal-balance")
	agent.QueryString(fmt.Sprintf("start_date=%v&end_date=%v", reqData.Query.StartDate, reqData.Query.EndDate))
	agent.InsecureSkipVerify() // TEMP!
	_, body, errs := agent.JSON(reqData.Body).Bytes()

	for _, err := range errs {
		if err != nil {
			log.Printf("Something went wrong with the request at /invoice/overal-balance for metrics query with id %v: %v\n", query.ID, err)
		}
	}

	var response SSAPIMetricsResponse
	json.Unmarshal(body, &response)

	if response.Errors != nil {
		response.ReqMsg = fmt.Sprintf("%s\n%s", response.ReqMsg, formatMetricsErrs(&response))
	}

	query.Results = response.MetricsResult

	err := w.metricsService.UpdateMetrics(ctx, query)
	if err != nil {
		log.Printf("Something went wrong when updating metrics history. Metrics query with id %v will be on 'pending' state for ever: %v\n", query.ID, err)
	}

	key := fmt.Sprintf("reqstatus:metrics:%v", query.ID)
	db.Redis.Set(ctx, key, true, time.Minute)
}

func (w *SiareBGWorker) RequestInvoicePrinting(invoicePrint *models.InvoicePrint) {
	ctx, print := context.WithTimeout(context.Background(), time.Minute*10)
	defer print()

	reqBody := SSAPIPrintingRequest{
		InvoicePrint: invoicePrint,
	}

	agent := fiber.Post(SS_API_BASE_URL + "/invoice/print")
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			log.Printf("Something went wrong with the request at /invoice/print for printing with id %v: %v\n", invoicePrint.ID, err)
		}
	}

	var response SSAPIPrintingResponse
	json.Unmarshal(body, &response)

	if response.Errors != nil {
		response.Msg = fmt.Sprintf("%s\n%s", response.Msg, formatPrintingErrs(&response))
	}

	invoicePrint.ReqStatus = response.Status
	invoicePrint.ReqMsg = response.Msg
	invoicePrint.InvoicePDF = response.InvoicePDF
	if response.FileName != "" {
		invoicePrint.CustomFileName = response.FileName
	}

	err := w.printingService.UpdateInvoicePrinting(ctx, invoicePrint)
	if err != nil {
		log.Printf("Something went wrong when updating invoice printing history. Printing with id %v will be on 'pending' state for ever: %v\n", invoicePrint.ID, err)
	}

	key := fmt.Sprintf("reqstatus:printing:%v", invoicePrint.ID)
	db.Redis.Set(ctx, key, true, time.Minute)
}

func formatErrResponseField(fieldErrs []string, prefix string) string {
	var builder strings.Builder

	if fieldErrs != nil {
		builder.WriteString(prefix)
		for _, fieldName := range fieldErrs {
			builder.WriteString(" - ")
			builder.WriteString(fieldName)
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func formatInvoiceErrs(response *SSAPIInvoiceResponse) string {
	invalidFields := "Campos inválidos:\n"
	missingFields := "Campos faltando:\n"
	senderInvalidFields := "Campos inválidos no remetente:\n"
	senderMissingFields := "Campos faltando no remetente:\n"
	recipientInvalidFields := "Campos inválidos no destinatário:\n"
	recipientMissingFields := "Campos faltando no destinatário:\n"

	var clientErrs strings.Builder

	clientErrs.WriteString(formatErrResponseField(response.Errors.InvalidFields, invalidFields))
	clientErrs.WriteString(formatErrResponseField(response.Errors.MissingFields, missingFields))
	if response.Errors.Sender != nil {
		clientErrs.WriteString(formatErrResponseField(response.Errors.Sender.InvalidFields, senderInvalidFields))
		clientErrs.WriteString(formatErrResponseField(response.Errors.Sender.MissingFields, senderMissingFields))
	}
	if response.Errors.Recipient != nil {
		clientErrs.WriteString(formatErrResponseField(response.Errors.Recipient.InvalidFields, recipientInvalidFields))
		clientErrs.WriteString(formatErrResponseField(response.Errors.Recipient.MissingFields, recipientMissingFields))
	}

	return clientErrs.String()
}

func formatCancelingErrs(response *SSAPICancelingResponse) string {
	invalidFields := "Campos inválidos:\n"
	missingFields := "Campos faltando:\n"
	entityInvalidFields := "Campos inválidos na entidade:\n"
	entityMissingFields := "Campos faltando na entidade:\n"

	var clientErrs strings.Builder

	clientErrs.WriteString(formatErrResponseField(response.Errors.InvalidFields, invalidFields))
	clientErrs.WriteString(formatErrResponseField(response.Errors.MissingFields, missingFields))
	if response.Errors.Entity != nil {
		clientErrs.WriteString(formatErrResponseField(response.Errors.Entity.InvalidFields, entityInvalidFields))
		clientErrs.WriteString(formatErrResponseField(response.Errors.Entity.MissingFields, entityMissingFields))
	}

	return clientErrs.String()
}

func formatMetricsErrs(response *SSAPIMetricsResponse) string {
	invalidFields := "Campos inválidos:\n"
	missingFields := "Campos faltando:\n"
	entityInvalidFields := "Campos inválidos na entidade:\n"
	entityMissingFields := "Campos faltando na entidade:\n"

	var clientErrs strings.Builder

	clientErrs.WriteString(formatErrResponseField(response.Errors.InvalidFields, invalidFields))
	clientErrs.WriteString(formatErrResponseField(response.Errors.MissingFields, missingFields))
	if response.Errors.Entity != nil {
		clientErrs.WriteString(formatErrResponseField(response.Errors.Entity.InvalidFields, entityInvalidFields))
		clientErrs.WriteString(formatErrResponseField(response.Errors.Entity.MissingFields, entityMissingFields))
	}

	return clientErrs.String()
}

func formatPrintingErrs(response *SSAPIPrintingResponse) string {
	invalidFields := "Campos inválidos:\n"
	missingFields := "Campos faltando:\n"
	entityInvalidFields := "Campos inválidos na entidade:\n"
	entityMissingFields := "Campos faltando na entidade:\n"

	var clientErrs strings.Builder

	clientErrs.WriteString(formatErrResponseField(response.Errors.InvalidFields, invalidFields))
	clientErrs.WriteString(formatErrResponseField(response.Errors.MissingFields, missingFields))
	if response.Errors.Entity != nil {
		clientErrs.WriteString(formatErrResponseField(response.Errors.Entity.InvalidFields, entityInvalidFields))
		clientErrs.WriteString(formatErrResponseField(response.Errors.Entity.MissingFields, entityMissingFields))
	}

	return clientErrs.String()
}
