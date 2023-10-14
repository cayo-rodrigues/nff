package bgworkers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/services"
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
	ShouldDownload  bool `json:"should_download"`
	ShouldNotFinish bool `json:"should_not_finish"`
}

type SSAPIInvoiceResponse struct {
	Msg                string              `json:"msg"`
	IsAwaitingAnalisys bool                `json:"is_awaiting_analisys"`
	InvoiceNumber      string              `json:"invoice_id"`
	InvoiceProtocol    string              `json:"invoice_protocol"`
	InvoicePDF         string              `json:"invoice_pdf"`
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

func SiareRequestInvoice(invoice *models.Invoice) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	reqBody := SSAPIInvoiceRequest{
		Invoice:         invoice,
		ShouldDownload:  false,
		ShouldNotFinish: true,
	}

	agent := fiber.Post(SS_API_BASE_URL + "/invoice/request")
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			log.Printf("Something went wrong with the request at /invoice/request for invoice with id %v: %v\n", invoice.Id, err)
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

	err := services.UpdateInvoice(ctx, invoice)
	if err != nil {
		log.Printf("Something went wrong when updating invoice history. Invoice with id %v will be on 'pending' state for ever: %v\n", invoice.Id, err)
	}

	key := fmt.Sprintf("reqstatus:invoice:%v", invoice.Id)
	db.Redis.Set(ctx, key, true, time.Minute)
}

func SiareRequestInvoiceCanceling(invoiceCancel *models.InvoiceCancel) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
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
			log.Printf("Something went wrong with the request at /invoice/cancel for canceling with id %v: %v\n", invoiceCancel.Id, err)
		}
	}

	var response SSAPICancelingResponse
	json.Unmarshal(body, &response)

	if response.Errors != nil {
		response.Msg = fmt.Sprintf("%s\n%s", response.Msg, formatCancelingErrs(&response))
	}

	invoiceCancel.ReqStatus = response.Status
	invoiceCancel.ReqMsg = response.Msg

	err := services.UpdateInvoiceCanceling(ctx, invoiceCancel)
	if err != nil {
		log.Printf("Something went wrong when updating invoice canceling history. Canceling with id %v will be on 'pending' state for ever: %v\n", invoiceCancel.Id, err)
	}

	key := fmt.Sprintf("reqstatus:canceling:%v", invoiceCancel.Id)
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
