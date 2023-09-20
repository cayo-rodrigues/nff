package handlers

import (
	"html/template"
	"net/http"

	"github.com/cayo-rodrigues/nff/web/internal/models"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/cayo-rodrigues/nff/web/internal/workers"
)

type InvoicesPage struct {
	tmpl *template.Template
}

type InvoicesPageData struct {
	IsAuthenticated bool
	Invoices        *[]models.Invoice
	GeneralError    string
	FormMsg         string
	FormSuccess     bool
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
	}
	invoice := *models.NewEmptyInvoice()
	entities, err := workers.ListEntities(r.Context())
	if err != nil {
		data.GeneralError = err.Error()
		utils.ErrorResponse(w, "general-error", page.tmpl, "layout", data)
		return
	}
	invoice.FormSelectFields.Entities = entities

	data.Invoices = &[]models.Invoice{invoice}
	page.tmpl.ExecuteTemplate(w, "layout", data)
}

func (page *InvoicesPage) GetInvoiceFormSection(w http.ResponseWriter, r *http.Request) {
	invoice := *models.NewEmptyInvoice()
	entities, err := workers.ListEntities(r.Context())
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}
	invoice.FormSelectFields.Entities = entities
	page.tmpl.ExecuteTemplate(w, "invoice-form-section", invoice)
}

func (page *InvoicesPage) RequireInvoices(w http.ResponseWriter, r *http.Request) {
	data := &InvoicesPageData{
		FormMsg:     "Requerimento efetuado com sucesso! Acompanhe o progresso clicando no botão ao lado.",
		FormSuccess: true,
	}
	entities, err := workers.ListEntities(r.Context())
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}

	invoices, err := models.NewInvoiceListFromForm(r, entities)
	if err != nil {
		utils.GeneralErrorResponse(w, err, page.tmpl)
		return
	}

	for _, invoice := range *invoices {
		if !invoice.IsValid() && data.FormSuccess {
			data.FormMsg = "Corrija os campos abaixo."
			data.FormSuccess = false
		}
	}

	// i would call ss-api here in case data.FormSuccess == true

	data.Invoices = invoices
	page.tmpl.ExecuteTemplate(w, "main", data)
}
