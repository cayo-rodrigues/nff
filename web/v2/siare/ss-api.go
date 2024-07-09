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

var instance *SSApiClient

func NewSSApiClient() *SSApiClient {
	instance = &SSApiClient{
		BaseUrl: os.Getenv("SS_API_BASE_URL"),
		DB:      database.GetDB(),
		Endpoints: &SSApiEndpoints{
			IssueInvoice:  "/invoice/request", // AJUSTAR NOME DO ENDPOINT NA SS-API para /invoice/issue
			CancelInvoice: "/invoice/cancel",
			PrintInvoice:  "/invoice/print",
			Metrics:       "/metrics",
		},
	}

	return instance
}

func GetSSApiClient() *SSApiClient {
	if instance == nil {
		return NewSSApiClient()
	}

	return instance
}

func (c *SSApiClient) IssueInvoice(invoice *models.Invoice) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	resourceName := "invoice-issue"
	defer finishOperation(ctx, c.DB.Redis, resourceName, invoice.CreatedBy, invoice)

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
		return
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

	storage.UpdateInvoice(ctx, invoice)
}

func (c *SSApiClient) CancelInvoice(invoiceCancel *models.InvoiceCancel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	resourceName := "invoice-cancel"
	defer finishOperation(ctx, c.DB.Redis, resourceName, invoiceCancel.CreatedBy, invoiceCancel)

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
		return err
	}

	if response.Errors != nil {
		response.Msg = fmt.Sprintf("%s\n%s", response.Msg, formatErrResponse(response.Errors))
	}

	invoiceCancel.ReqStatus = response.Status
	invoiceCancel.ReqMsg = response.Msg

	err = storage.UpdateInvoiceCanceling(ctx, invoiceCancel)

	// TODO
	// DEIXAR BONITO?
	go func(invoiceCancel *models.InvoiceCancel) {
		if invoiceCancel.ReqStatus == "success" {
			invoice, err := storage.RetrieveInvoice(context.Background(), 0, invoiceCancel.CreatedBy, invoiceCancel.InvoiceNumber)
			if err != nil {
				return
			}

			invoice.ReqStatus = "canceled"
			invoice.ReqMsg = fmt.Sprintf(
				"Nota Fiscal havia sido emitida com sucesso, porém foi cancelada em %s.\nJustificativa: %s",
				utils.FormatDatetimeAsBR(time.Now()),
				invoiceCancel.Justification,
			)
			err = storage.UpdateInvoice(context.Background(), invoice)
			if err != nil {
				return
			}
			innerResourceName := "invoice-cancel-from-invoice-card"
			notifyOperationResult(ctx, c.DB.Redis, innerResourceName, invoiceCancel.CreatedBy)
		}
	}(invoiceCancel)

	return err
}

func (c *SSApiClient) PrintInvoiceFromMetricsRecord(p *models.InvoicePrint, recordID, userID int) error {
	resourceName := "invoice-print-from-record"
	defer notifyOperationResult(context.Background(), c.DB.Redis, resourceName, userID)

	err := c.PrintInvoice(p)
	if err != nil {
		return err
	}

	record := models.NewMetricsResult()
	record.InvoicePDF = p.InvoicePDF
	record.ID = recordID
	record.CreatedBy = userID

	return storage.UpdateMetricsResultRecord(context.Background(), record)
}

func (c *SSApiClient) PrintInvoice(invoicePrint *models.InvoicePrint) error {
	ctx, print := context.WithTimeout(context.Background(), time.Minute*5)
	defer print()

	resourceName := "invoice-print"
	defer finishOperation(ctx, c.DB.Redis, resourceName, invoicePrint.CreatedBy, invoicePrint)

	reqBody := SSApiPrintingRequest{
		InvoicePrint: invoicePrint,
	}

	agent := fiber.Post(c.BaseUrl + c.Endpoints.PrintInvoice)
	// TEMP!
	agent.InsecureSkipVerify()
	_, body, errs := agent.JSON(reqBody).Bytes()

	for _, err := range errs {
		if err != nil {
			requestErrLog(c.Endpoints.PrintInvoice, resourceName, invoicePrint.ID, err)
		}
	}

	var response SSApiPrintingResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		jsonUnmarchalErrLog(c.Endpoints.PrintInvoice, resourceName, invoicePrint.ID, err)
		return err
	}

	if response.Errors != nil {
		response.Msg = fmt.Sprintf("%s\n%s", response.Msg, formatErrResponse(response.Errors))
	}

	invoicePrint.ReqStatus = response.Status
	invoicePrint.ReqMsg = response.Msg
	invoicePrint.InvoicePDF = response.InvoicePDF
	invoicePrint.FileName = response.FileName

	return storage.UpdateInvoicePrinting(ctx, invoicePrint)
}

func (c *SSApiClient) GetMetrics(metrics *models.Metrics) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	resourceName := "metrics"
	defer finishOperation(ctx, c.DB.Redis, resourceName, metrics.CreatedBy, metrics)

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

	for _, err := range errs {
		if err != nil {
			requestErrLog(c.Endpoints.Metrics, resourceName, metrics.ID, err)
		}
	}

	var response SSApiMetricsResponse
	err := json.Unmarshal(body, &response)
	if err != nil {
		jsonUnmarchalErrLog(c.Endpoints.Metrics, resourceName, metrics.ID, err)
		return
	}

	if response.Errors != nil {
		response.ReqMsg = fmt.Sprintf("%s\n%s", response.ReqMsg, formatErrResponse(response.Errors))
	}

	metrics.MetricsResult = response.MetricsResult
	metrics.ReqStatus = response.ReqStatus
	metrics.ReqMsg = response.ReqMsg

	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		defer wg.Done()
		storage.UpdateMetrics(ctx, metrics)
	}()
	go func() {
		defer wg.Done()
		if metrics.MetricsResult.Total == nil {
			return
		}
		storage.CreateMetricsResult(ctx, metrics.MetricsResult.Total, "total", metrics.ID, metrics.CreatedBy, metrics.Entity.ID)
	}()
	go func() {
		defer wg.Done()
		if metrics.MetricsResult.Months == nil {
			return
		}
		storage.BulkCreateMetricsResults(ctx, metrics.MetricsResult.Months, "month", metrics.ID, metrics.CreatedBy, metrics.Entity.ID)
	}()
	go func() {
		defer wg.Done()
		if metrics.MetricsResult.Records == nil {
			return
		}
		storage.BulkCreateMetricsResults(ctx, metrics.MetricsResult.Records, "record", metrics.ID, metrics.CreatedBy, metrics.Entity.ID)
	}()

	wg.Wait()
	return
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

func finishOperation(ctx context.Context, redisClient *database.Redis, resourceName string, userID int, n models.Notification) {
	// primeiro, adicionar Notification em uma fila no redis para que seja consumida por outro endpoint (que será chamado quando o sininho receber notificação)
	pushOperationToNotificationQueue(ctx, redisClient, userID, n)
	go redisClient.ClearCache(ctx, userID, resourceName)
	go notifyOperationResult(ctx, redisClient, resourceName, userID)
}

func notifyOperationResult(ctx context.Context, redisClient *database.Redis, resourceName string, userID int) {
	channel := fmt.Sprintf("%d:%s-finished", userID, resourceName)
	redisClient.Publish(ctx, channel, nil)
}

func pushOperationToNotificationQueue(ctx context.Context, redisClient *database.Redis, userID int, n models.Notification) {
	queueName := fmt.Sprintf("notification-queue:%d", userID)
	notification := models.JsonSerializableNotification(n, userID)

	notificationData, err := json.Marshal(notification)
	if err != nil {
		log.Println("Failed to marshal notification:", err)
		return
	}

	err = redisClient.RPush(ctx, queueName, notificationData).Err()
	if err != nil {
		log.Println("Failed to enqueue notification:", err)
	}
}
