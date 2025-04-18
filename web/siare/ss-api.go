package siare

import (
	"context"
	"encoding/hex"
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
	"github.com/cayo-rodrigues/nff/web/utils/cryptoutils"
	"github.com/gofiber/fiber/v2"
)

type SSApiClient struct {
	BaseUrl       string
	DB            *database.Database
	Endpoints     *SSApiEndpoints
	DecryptionKey []byte
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

func (c *SSApiClient) WithDecryptionKey(key []byte) *SSApiClient {
	c.DecryptionKey = key
	return c
}

func (c *SSApiClient) IssueInvoice(invoice *models.Invoice) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	resourceName := "invoice-issue"
	defer finishOperation(ctx, c.DB.Redis, resourceName, invoice.CreatedBy, invoice)

	pwd, err := hex.DecodeString(invoice.Sender.Password)
	if err != nil {
		log.Println("Failed to decode invoice.Sender.Password hexadecimal string:", err)
		return
	}

	decryptedPwd, err := cryptoutils.Decrypt(c.DecryptionKey, pwd)
	if err != nil {
		log.Println("Failed to decrypt invoice.Sender.Password:", err)
		return
	}
	invoice.Sender.Password = string(decryptedPwd)


	reqBody := SSApiInvoiceRequest{
		Invoice:            invoice,
		ShouldDownload:     true,
		ShouldAbortMission: false,
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
	err = json.Unmarshal(body, &response)
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

	pwd, err := hex.DecodeString(invoiceCancel.Entity.Password)
	if err != nil {
		log.Println("Failed to decode invoiceCancel.Entity.Password hexadecimal string:", err)
		return err
	}

	decryptedPwd, err := cryptoutils.Decrypt(c.DecryptionKey, pwd)
	if err != nil {
		log.Println("Failed to decrypt invoiceCancel.Entity.Password:", err)
		return err
	}
	invoiceCancel.Entity.Password = string(decryptedPwd)

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
	err = json.Unmarshal(body, &response)
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
		if invoiceCancel.ReqStatus != "success" {
			return
		}
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

		// TODO
		// pode melhorar?
		c.DB.Redis.ClearCache(ctx, invoiceCancel.CreatedBy, "invoice-issue")
		notifyOperationResult(ctx, c.DB.Redis, invoiceCancel.CreatedBy, invoiceCancel.ID)
	}(invoiceCancel)

	return err
}

func (c *SSApiClient) PrintInvoiceFromMetricsRecord(p *models.InvoicePrint, recordID, userID int) error {
	// TODO
	// pode melhorar?
	ctx := context.Background()
	defer c.DB.Redis.ClearCache(ctx, userID, "invoice-print")
	defer c.DB.Redis.ClearCache(ctx, userID, "metrics")
	defer notifyOperationResult(ctx, c.DB.Redis, userID, p.ID)

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

	pwd, err := hex.DecodeString(invoicePrint.Entity.Password)
	if err != nil {
		log.Println("Failed to decode invoicePrint.Entity.Password hexadecimal string:", err)
		return err
	}

	decryptedPwd, err := cryptoutils.Decrypt(c.DecryptionKey, pwd)
	if err != nil {
		log.Println("Failed to decrypt invoicePrint.Entity.Password:", err)
		return err
	}
	invoicePrint.Entity.Password = string(decryptedPwd)

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
	err = json.Unmarshal(body, &response)
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

	pwd, err := hex.DecodeString(metrics.Entity.Password)
	if err != nil {
		log.Println("Failed to decode metrics.Entity.Password hexadecimal string:", err)
		return
	}

	decryptedPwd, err := cryptoutils.Decrypt(c.DecryptionKey, pwd)
	if err != nil {
		log.Println("Failed to decrypt metrics.Entity.Password:", err)
		return 
	}
	metrics.Entity.Password = string(decryptedPwd)

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
	err = json.Unmarshal(body, &response)
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

	otherEntitiesIes := []any{}
	for _, record := range metrics.MetricsResult.Records {
		if record.IsPositive {
			record.InvoiceSender = fmt.Sprintf("%s - %s", record.InvoiceSender, metrics.Entity.Name)
		} else {
			otherEntitiesIes = append(otherEntitiesIes, record.InvoiceSender)
		}
	}
	f := models.NewFilters().Where("created_by = ").Placeholder(metrics.CreatedBy)
	f.And("ie").In(otherEntitiesIes)
	for _, ie := range otherEntitiesIes {
		f.Or("other_ies").ILike().WildPlaceholder(ie)
	}

	query := new(strings.Builder)
	query.WriteString("SELECT name, ie, other_ies FROM entities")
	query.WriteString(f.String())

	db := database.GetDB()

	rows, err := db.SQLite.QueryContext(ctx, query.String(), f.Values()...)
	if err != nil {
		log.Printf("Error on get metrics post request, aborting operation. Metrics with id %v will be on 'pending' state for ever. Error: %v\n", metrics.ID, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var entityName, entityIe string
		var entityOtherIesJSON []byte
		var entityOtherIes []string
		if err := rows.Scan(&entityName, &entityIe, &entityOtherIesJSON); err != nil {
			log.Printf("Error scanning entity rows, skipping entity. Error: %v", err)
			continue
		}

		err = json.Unmarshal(entityOtherIesJSON, &entityOtherIes)
		if err != nil {
			log.Println("Error unmarshaling entity other ies, skipping entity. Error: ", err)
			continue
		}

		for _, record := range metrics.MetricsResult.Records {
			if record.InvoiceSender == entityIe {
				record.InvoiceSender = fmt.Sprintf("%s - %s", entityIe, entityName)
				break
			}
			foundSecondaryIe := false
			for _, secondaryIe := range entityOtherIes {
				if record.InvoiceSender == secondaryIe {
					record.InvoiceSender = fmt.Sprintf("%s - %s", secondaryIe, entityName)
					foundSecondaryIe = true
					break
				}
			}

			if foundSecondaryIe {
				break
			}
		}
	}

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
}
