package siare

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	BaseUrl   string
	DB        *database.Database
	Endpoints *SSApiEndpoints
}

type SSApiEndpoints struct {
	IssueInvoice  string
	CancelInvoice string
	PrintInvoice  string
	Metrics       string
}

func NewSSApiClient() *SSApiClient {
	return &SSApiClient{
		BaseUrl: os.Getenv("SS_API_BASE_URL"),
		DB:      database.GetDB(),
		Endpoints: &SSApiEndpoints{
			IssueInvoice:  "/invoice/issue", // AJUSTAR NOME DO ENDPOINT NA SS-API
			CancelInvoice: "/invoice/cancel",
			PrintInvoice:  "/invoice/print",
			Metrics:       "/metrics",
		},
	}
}

func (c *SSApiClient) IssueInvoice(invoice *models.Invoice) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	reqBody := SSApiInvoiceRequest{
		Invoice:        invoice,
		ShouldDownload: true,
	}

	agent := fiber.Post(c.BaseUrl + c.Endpoints.IssueInvoice)
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			requestErrLog(c.Endpoints.IssueInvoice, "invoice", invoice.ID, err)
		}
	}

	var response SSApiInvoiceResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		jsonUnmarchalErrLog(c.Endpoints.IssueInvoice, "invoice", invoice.ID, err)
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

func (c *SSApiClient) CancelInvoice(invoiceCancel *models.InvoiceCancel) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	reqBody := SSApiCancelingRequest{
		InvoiceCancel: invoiceCancel,
	}

	agent := fiber.Post(c.BaseUrl + c.Endpoints.CancelInvoice)
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			requestErrLog(c.Endpoints.CancelInvoice, "canceling", invoiceCancel.ID, err)
		}
	}

	var response SSApiCancelingResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		jsonUnmarchalErrLog(c.Endpoints.CancelInvoice, "canceling", invoiceCancel.ID, err)
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

func (c *SSApiClient) PrintInvoice(invoicePrint *models.InvoicePrint) {
	ctx, print := context.WithTimeout(context.Background(), time.Minute*10)
	defer print()

	reqBody := SSApiPrintingRequest{
		InvoicePrint: invoicePrint,
	}

	agent := fiber.Post(c.BaseUrl + c.Endpoints.PrintInvoice)
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			requestErrLog(c.Endpoints.PrintInvoice, "printing", invoicePrint.ID, err)
		}
	}

	var response SSApiPrintingResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		jsonUnmarchalErrLog(c.Endpoints.PrintInvoice, "printing", invoicePrint.ID, err)
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

	agent := fiber.Get(c.BaseUrl + c.Endpoints.Metrics)
	agent.QueryString(queryString)
	agent.InsecureSkipVerify() // TEMP!
	_, body, errs := agent.JSON(reqData.Body).Bytes()

	fmt.Println("response!", string(body))

	for _, err := range errs {
		if err != nil {
			requestErrLog(c.Endpoints.Metrics, "metrics", metrics.ID, err)
		}
	}

	var response SSApiMetricsResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		jsonUnmarchalErrLog(c.Endpoints.Metrics, "metrics", metrics.ID, err)
	}

	if response.Errors != nil {
		response.ReqMsg = fmt.Sprintf("%s\n%s", response.ReqMsg, formatErrResponse(response.Errors))
	}

	metrics.MetricsResult = response.MetricsResult

	if err == nil {
		var wg sync.WaitGroup
		wg.Add(4)

		go func() {
			defer wg.Done()
			storage.UpdateMetrics(ctx, metrics)
		}()
		go func() {
			defer wg.Done()
			storage.CreateMetricsResult(ctx, metrics.MetricsResult.Total, "total", metrics.ID, metrics.CreatedBy, metrics.Entity.ID)
		}()
		go func() {
			defer wg.Done()
			storage.BulkCreateMetricsResults(ctx, metrics.MetricsResult.Months, "month", metrics.ID, metrics.CreatedBy, metrics.Entity.ID)
		}()
		go func() {
			defer wg.Done()
			storage.BulkCreateMetricsResults(ctx, metrics.MetricsResult.Records, "record", metrics.ID, metrics.CreatedBy, metrics.Entity.ID)
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

	clientErrs.WriteString(formatErrResponseField(response.InvalidFields, InvalidFieldsPrefix))
	clientErrs.WriteString(formatErrResponseField(response.MissingFields, MissingFieldsPrefix))
	if response.Entity != nil {
		clientErrs.WriteString(formatErrResponseField(response.Entity.InvalidFields, EntityInvalidFieldsPrefix))
		clientErrs.WriteString(formatErrResponseField(response.Entity.MissingFields, EntityMissingFieldsPrefix))
	}
	if response.Sender != nil {
		clientErrs.WriteString(formatErrResponseField(response.Sender.InvalidFields, SenderInvalidFieldsPrefix))
		clientErrs.WriteString(formatErrResponseField(response.Sender.MissingFields, SenderMissingFieldsPrefix))
	}
	if response.Recipient != nil {
		clientErrs.WriteString(formatErrResponseField(response.Recipient.InvalidFields, RecipientInvalidFieldsPrefix))
		clientErrs.WriteString(formatErrResponseField(response.Recipient.MissingFields, RecipientMissingFieldsPrefix))
	}

	return clientErrs.String()
}

func requestErrLog(endpoint, resourceName string, resourceID int, err error) {
	log.Printf("Something went wrong with the request at %s for %s with id %v: %v\n", endpoint, resourceName, resourceID, err)
}

func jsonUnmarchalErrLog(endpoint, resourceName string, resourceID int, err error) {
	log.Printf("Something went wrong trying to unmarshal ss-api %s json response. %s with id %v will be on 'pending' state for ever: %v\n", endpoint, resourceName, resourceID, err)
}
