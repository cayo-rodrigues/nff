package siare

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/storage"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

type SSApiClient struct {
	BaseUrl string
	DB      *database.Database
}

func NewSSApiClient(baseUrl string) *SSApiClient {
	return &SSApiClient{
		BaseUrl: baseUrl,
		DB:      database.GetDB(),
	}
}

func (c *SSApiClient) RequestInvoice(invoice *models.Invoice) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	reqBody := SSApiInvoiceRequest{
		Invoice:        invoice,
		ShouldDownload: true,
	}

	agent := fiber.Post(c.BaseUrl + "/invoice/issue")
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			log.Printf("Something went wrong with the request at /invoice/issue for invoice with id %v: %v\n", invoice.ID, err)
		}
	}

	var response SSApiInvoiceResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Something went wrong trying to unmarshal ss-api /invoice/request json response. Invoice with id %v will be on 'pending' state for ever: %v\n", invoice.ID, err)
	}

	if response.Errors != nil {
		response.Msg = fmt.Sprintf("%s\n%s", response.Msg, formatErrResponse(response.Errors))
	}

	invoice.Number = response.InvoiceNumber
	invoice.Protocol = response.InvoiceProtocol
	invoice.ReqMsg = response.Msg
	invoice.ReqStatus = response.Status
	invoice.PDF = response.InvoicePDF
	invoice.FileName = response.FileName

	if err == nil {
		storage.UpdateInvoice(ctx, invoice)
	}

	c.DB.Redis.ClearCache(ctx, invoice.CreatedBy, "invoices")

	key := fmt.Sprintf("reqstatus:invoice:%v", invoice.ID)
	c.DB.Redis.Set(ctx, key, true, time.Minute)
}

func (c *SSApiClient) RequestInvoiceCanceling(invoiceCancel *models.InvoiceCancel) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	reqBody := SSApiCancelingRequest{
		InvoiceCancel: invoiceCancel,
	}

	agent := fiber.Post(c.BaseUrl + "/invoice/cancel")
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			log.Printf("Something went wrong with the request at /invoice/cancel for canceling with id %v: %v\n", invoiceCancel.ID, err)
		}
	}

	var response SSApiCancelingResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Something went wrong trying to unmarshal ss-api /invoice/cancel json response. Canceling with id %v will be on 'pending' state for ever: %v\n", invoiceCancel.ID, err)
	}

	if response.Errors != nil {
		response.Msg = fmt.Sprintf("%s\n%s", response.Msg, formatErrResponse(response.Errors))
	}

	invoiceCancel.ReqStatus = response.Status
	invoiceCancel.ReqMsg = response.Msg

	if err == nil {
		storage.UpdateInvoiceCanceling(ctx, invoiceCancel)
	}

	c.DB.Redis.ClearCache(ctx, invoiceCancel.CreatedBy, "invoices-cancel")

	key := fmt.Sprintf("reqstatus:canceling:%v", invoiceCancel.ID)
	c.DB.Redis.Set(ctx, key, true, time.Minute)
}

func (c *SSApiClient) RequestInvoicePrinting(invoicePrint *models.InvoicePrint) {
	ctx, print := context.WithTimeout(context.Background(), time.Minute*10)
	defer print()

	reqBody := SSApiPrintingRequest{
		InvoicePrint: invoicePrint,
	}

	agent := fiber.Post(c.BaseUrl + "/invoice/print")
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			log.Printf("Something went wrong with the request at /invoice/print for printing with id %v: %v\n", invoicePrint.ID, err)
		}
	}

	var response SSApiPrintingResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Something went wrong trying to unmarshal ss-api /invoice/print json response. Printing with id %v will be on 'pending' state for ever: %v\n", invoicePrint.ID, err)
	}

	if response.Errors != nil {
		response.Msg = fmt.Sprintf("%s\n%s", response.Msg, formatErrResponse(response.Errors))
	}

	invoicePrint.ReqStatus = response.Status
	invoicePrint.ReqMsg = response.Msg
	invoicePrint.InvoicePDF = response.InvoicePDF
	invoicePrint.FileName = response.FileName

	if err == nil {
		storage.UpdateInvoicePrinting(ctx, invoicePrint)
	}

	c.DB.Redis.ClearCache(ctx, invoicePrint.CreatedBy, "invoices-print")

	key := fmt.Sprintf("reqstatus:printing:%v", invoicePrint.ID)
	c.DB.Redis.Set(ctx, key, true, time.Minute)
}

func (c *SSApiClient) GetMetrics(metrics *models.Metrics) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	reqData := SSApiMetricsRequest{
		Body: &SSApiMetricsRequestBody{
			Entity: metrics.Entity,
		},
		Query: &SSApiMetricsRequestQuery{
			StartDate:      utils.FormatDateAsBR(metrics.StartDate),
			EndDate:        utils.FormatDateAsBR(metrics.EndDate),
			IncludeRecords: true,
		},
	}

	queryString := fmt.Sprintf(
		"start_date=%v&end_date=%v",
		reqData.Query.StartDate, reqData.Query.EndDate,
	)

	if reqData.Query.IncludeRecords {
		queryString += fmt.Sprintf("&include_records=%v", reqData.Query.IncludeRecords)
	}

	agent := fiber.Get(c.BaseUrl + "/metrics")
	agent.QueryString(queryString)
	agent.InsecureSkipVerify() // TEMP!
	_, body, errs := agent.JSON(reqData.Body).Bytes()

	for _, err := range errs {
		if err != nil {
			log.Printf("Something went wrong with the request at /metrics for metrics query with id %v: %v\n", metrics.ID, err)
		}
	}

	var response SSApiMetricsResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Something went wrong trying to unmarshal ss-api /metrics json response. Metrics query with id %v will be on 'pending' state for ever: %v\n", metrics.ID, err)
	}

	if response.Errors != nil {
		response.ReqMsg = fmt.Sprintf("%s\n%s", response.ReqMsg, formatErrResponse(response.Errors))
	}

	metrics.MetricsResult = response.MetricsResult

	if err == nil {
		var wg *sync.WaitGroup

		wg.Add(3)

		go func() {
			defer wg.Done()
			storage.UpdateMetrics(ctx, metrics)
		}()
		go func() {
			defer wg.Done()
			storage.BulkCreateMetricsResults(ctx, metrics.Months, "month", metrics.ID, metrics.CreatedBy, metrics.Entity.ID)
		}()
		go func() {
			defer wg.Done()
			storage.BulkCreateMetricsResults(ctx, metrics.Records, "record", metrics.ID, metrics.CreatedBy, metrics.Entity.ID)
		}()

		wg.Wait()
	}

	c.DB.Redis.ClearCache(ctx, metrics.CreatedBy, "metrics")

	key := fmt.Sprintf("reqstatus:metrics:%v", metrics.ID)
	c.DB.Redis.Set(ctx, key, true, time.Minute)
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

func formatErrResponse(response *SSApiErrorResponse) string {
	var clientErrs strings.Builder

	clientErrs.WriteString(formatErrResponseField(response.InvalidFields, InvalidFields))
	clientErrs.WriteString(formatErrResponseField(response.MissingFields, MissingFields))
	if response.Entity != nil {
		clientErrs.WriteString(formatErrResponseField(response.Entity.InvalidFields, EntityInvalidFields))
		clientErrs.WriteString(formatErrResponseField(response.Entity.MissingFields, EntityMissingFields))
	}
	if response.Sender != nil {
		clientErrs.WriteString(formatErrResponseField(response.Sender.InvalidFields, SenderInvalidFields))
		clientErrs.WriteString(formatErrResponseField(response.Sender.MissingFields, SenderMissingFields))
	}
	if response.Recipient != nil {
		clientErrs.WriteString(formatErrResponseField(response.Recipient.InvalidFields, RecipientInvalidFields))
		clientErrs.WriteString(formatErrResponseField(response.Recipient.MissingFields, RecipientMissingFields))
	}

	return clientErrs.String()
}
