package bgworkers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/services"
	"github.com/gofiber/fiber/v2"
)

var SS_API_BASE_URL = ""

type SSAPIInvoiceRequest struct {
	*models.Invoice
	ShouldDownload  bool `json:"should_download"`
	ShouldNotFinish bool `json:"should_not_finish"`
}

type SSAPIInvoiceResponse struct {
	Msg                string `json:"msg"`
	IsAwaitingAnalisys bool   `json:"is_awaiting_analisys"`
	InvoiceNumber      string `json:"invoice_id"`
	InvoiceProtocol    string `json:"invoice_protocol"`
	InvoicePDF         string `json:"invoice_pdf"`
	Errors             string `json:"errors"`
	Status             string `json:"status"`
}

type SSAPICancelingRequest struct {
	*models.InvoiceCancel
}

type SSAPICancelingResponse struct {
	Msg    string `json:"msg"`
	Errors string `json:"errors"`
	Status string `json:"status"`
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
	statusCode, body, errs := agent.JSON(reqBody).Bytes()
	// TODO handle errors key on body appropriately
	fmt.Println(statusCode)
	fmt.Println(string(body))
	fmt.Println(errs)

	var response SSAPIInvoiceResponse
	json.Unmarshal(body, &response)

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
	statusCode, body, errs := agent.JSON(reqBody).Bytes()
	// TODO handle errors key on body appropriately
	fmt.Println(statusCode)
	fmt.Println(string(body))
	fmt.Println(errs)

	var response SSAPICancelingResponse
	json.Unmarshal(body, &response)

	invoiceCancel.ReqStatus = response.Status
	invoiceCancel.ReqMsg = response.Msg

	err := services.UpdateInvoiceCanceling(ctx, invoiceCancel)
	if err != nil {
		log.Printf("Something went wrong when updating invoice canceling history. Canceling with id %v will be on 'pending' state for ever: %v\n", invoiceCancel.Id, err)
	}

	key := fmt.Sprintf("reqstatus:canceling:%v", invoiceCancel.Id)
	db.Redis.Set(ctx, key, true, time.Minute)
}
