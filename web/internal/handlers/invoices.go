package handlers

import (
	"html/template"
	"net/http"

	"github.com/cayo-rodrigues/nff/web/internal/models"
)

type InvoicesPage struct {
	tmpl *template.Template
}

type InvoicesPageData struct {
	IsAuthenticated bool
	Invoices        *[]models.Invoice
	GeneralError    string
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
	page.tmpl.ExecuteTemplate(w, "layout", &InvoicesPageData{
		IsAuthenticated: true,
		Invoices:        &[]models.Invoice{},
	})
}
