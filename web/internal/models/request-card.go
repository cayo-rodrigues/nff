package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/globals"
	"github.com/cayo-rodrigues/nff/web/internal/utils"
	"github.com/gofiber/fiber/v2"
)

type ReqCardFilters struct {
	FromDate time.Time
	ToDate   time.Time
}

type RequestCard struct {
	From              string
	To                string
	ToPrefix          string
	FeedbackTextColor string
	ReqStatus         string
	ReqMsg            string
	ResourceName      string
	TargetForm        string
	OverviewType      string
	HasItems          bool
	IsFinished        bool
	IsDownloadable    bool
	DownloadLink      string
	ShouldCheckStatus bool
}

// req may be either *Invoice, *InvoiceCancel, *MetricsQuery or *InvoicePrint
func NewRequestCard(req any) *RequestCard {
	var prefix, reqStatus, reqMsg, from, to, overviewType, resourceName, targetForm, downloadLink string
	var hasItems, isDownloadable bool

	switch r := req.(type) {
	case *Invoice:
		from = r.Sender.Name
		to = r.Recipient.Name
		reqStatus = r.ReqStatus
		reqMsg = r.ReqMsg
		resourceName = "invoices"
		targetForm = "#invoice-form"
		overviewType = "invoice"
		hasItems = true
		isDownloadable = true
		downloadLink = r.PDF
	case *InvoiceCancel:
		from = r.Entity.Name
		to = r.Number
		prefix = "NFA-"
		reqStatus = r.ReqStatus
		reqMsg = r.ReqMsg
		resourceName = "invoices/cancel"
		targetForm = "#invoice-cancel-form"
		overviewType = "canceling"
	case *MetricsQuery:
		from = r.Entity.Name
		to = fmt.Sprintf("%v - %v", utils.FormatDateAsBR(r.StartDate), utils.FormatDateAsBR(r.EndDate))
		reqStatus = r.Results.ReqStatus
		reqMsg = r.Results.ReqMsg
		resourceName = "metrics"
		targetForm = "#metrics-form"
		overviewType = "metrics"
	case *InvoicePrint:
		from = r.Entity.Name
		to = r.InvoiceID
		if r.InvoiceIDType == "Número da NFA" {
			prefix = "NFA-"
		}
		reqStatus = r.ReqStatus
		reqMsg = r.ReqMsg
		resourceName = "invoices/print"
		targetForm = "#invoice-print-form"
		overviewType = "printing"
		isDownloadable = true
		downloadLink = r.InvoicePDF
	}

	feedbackColor := ""
	switch reqStatus {
	case "success":
		feedbackColor = "green"
	case "warning":
		feedbackColor = "yellow"
	case "error":
		feedbackColor = "red"
	case "pending":
		feedbackColor = "sky"
	}

	isFinished := reqStatus != "pending"

	return &RequestCard{
		From:              from,
		To:                to,
		ToPrefix:          prefix,
		FeedbackTextColor: feedbackColor,
		ReqStatus:         reqStatus,
		ReqMsg:            reqMsg,
		ResourceName:      resourceName,
		TargetForm:        targetForm,
		OverviewType:      overviewType,
		HasItems:          hasItems,
		IsFinished:        isFinished,
		IsDownloadable:    isDownloadable,
		DownloadLink:      downloadLink,
		ShouldCheckStatus: !isFinished,
	}
}

func NewRequestCardFilters() *ReqCardFilters {
	now := time.Now()
	return &ReqCardFilters{
		FromDate: utils.NDaysBefore(now, globals.DefaultFiltersDaysRange),
		ToDate:   now,
	}
}

func NewRawFiltersFromForm(c *fiber.Ctx) map[string]string {
	return map[string]string{
		"from_date": strings.TrimSpace(c.FormValue("from_date")),
		"to_date":   strings.TrimSpace(c.FormValue("to_date")),
	}
}
