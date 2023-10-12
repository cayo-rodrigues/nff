package bgworkers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/services"
)

func SiareRequestInvoice(invoice *models.Invoice) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	// i would call ss-api here
	time.Sleep(time.Second * 30)

	invoice.Number = "1234"
	invoice.Protocol = "9876"
	invoice.ReqStatus = "success"
	invoice.ReqMsg = "Requerimento efetuado com sucesso!"
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

	// do the thing
	time.Sleep(time.Second * 15)

	invoiceCancel.ReqStatus = "success"
	invoiceCancel.ReqMsg = "Cancelamento efetuado com sucesso!"
	err := services.UpdateInvoiceCanceling(ctx, invoiceCancel)
	if err != nil {
		log.Printf("Something went wrong when updating invoice canceling history. Canceling with id %v will be on 'pending' state for ever: %v\n", invoiceCancel.Id, err)
	}

	key := fmt.Sprintf("reqstatus:canceling:%v", invoiceCancel.Id)
	db.Redis.Set(ctx, key, true, time.Minute)
}
