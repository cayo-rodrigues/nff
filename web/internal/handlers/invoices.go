package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/models"
	// "github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/cayo-rodrigues/nff/web/internal/workers"
)

type InvoicesPage struct {
	tmpl *template.Template
}

type InvoicesPageData struct {
	IsAuthenticated  bool
	Invoice          *models.Invoice
	GeneralError     string
	FormMsg          string
	FormSuccess      bool
	FormSelectFields *models.InvoiceFormSelectFields
}

func NewInvoicesPage() *InvoicesPage {
	invoicesPage := &InvoicesPage{}
	invoicesPage.ParseTemplates()
	return invoicesPage
}

func (page *InvoicesPage) ParseTemplates() {
	page.tmpl = template.Must(template.ParseFiles(
		"internal/templates/layout.html",
		"internal/templates/invoices.html",
	))
}

func (page *InvoicesPage) Render(w http.ResponseWriter, r *http.Request) {
	data := &InvoicesPageData{
		IsAuthenticated: true,
		FormSelectFields: &models.InvoiceFormSelectFields{
			Operations:   &globals.InvoiceOperations,
			Cfops:        &globals.InvoiceCfops,
			BooleanField: &globals.InvoiceBooleanField,
			IcmsOptions:  &globals.InvoiceIcmsOptions,
		},
	}

	entities, err := workers.ListEntities(r.Context())
	if err != nil {
		data.GeneralError = err.Error()
		// utils.ErrorResponse(w, "general-error", page.tmpl, "layout", data)
		return
	}

	data.FormSelectFields.Entities = entities
	data.Invoice = models.NewEmptyInvoice()

	page.tmpl.ExecuteTemplate(w, "layout", data)
}

func (page *InvoicesPage) RequireInvoice(w http.ResponseWriter, r *http.Request) {
	data := &InvoicesPageData{
		FormMsg:     "Requerimento efetuado com sucesso! Acompanhe o progresso na sessão abaixo.",
		FormSuccess: true,
		FormSelectFields: &models.InvoiceFormSelectFields{
			Operations:   &globals.InvoiceOperations,
			Cfops:        &globals.InvoiceCfops,
			BooleanField: &globals.InvoiceBooleanField,
			IcmsOptions:  &globals.InvoiceIcmsOptions,
		},
	}

	entities, err := workers.ListEntities(r.Context())
	if err != nil {
		// utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}
	data.FormSelectFields.Entities = entities

	senderId, err := strconv.Atoi(r.PostFormValue("sender"))
	if err != nil {
		log.Println("Error converting sender id from string to int: ", err)
		// utils.GeneralErrorResponse(w, utils.InternalServerErr, page.tmpl)
		return
	}
	recipientId, err := strconv.Atoi(r.PostFormValue("recipient"))
	if err != nil {
		log.Println("Error converting recipient id from string to int: ", err)
		// utils.GeneralErrorResponse(w, utils.InternalServerErr, page.tmpl)
		return
	}

	sender, err := workers.RetrieveEntity(r.Context(), senderId)
	if err != nil {
		// utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}
	recipient, err := workers.RetrieveEntity(r.Context(), recipientId)
	if err != nil {
		// utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}

	invoice, err := models.NewInvoiceFromForm(r)
	if err != nil {
		// utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}

	invoice.Sender = sender
	invoice.Recipient = recipient

	if !invoice.IsValid() {
		data.FormMsg = "Corrija os campos abaixo."
		data.FormSuccess = false
	}

	data.Invoice = invoice

	// i would call ss-api here in case data.FormSuccess == true

	page.tmpl.ExecuteTemplate(w, "main", data)
}

func (page *InvoicesPage) GetItemFormSection(w http.ResponseWriter, r *http.Request) {
	item := *models.NewEmptyInvoiceItem()
	page.tmpl.ExecuteTemplate(w, "invoice-item-form-section", item)
}
