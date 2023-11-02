package models

import (
	"fmt"

	"github.com/cayo-rodrigues/nff/web/internal/utils"
)

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
	ShouldCheckStatus bool
}

// req may be either *Invoice, *InvoiceCancel or *MetricsQuery
func NewRequestCard(req any) *RequestCard {
	var prefix, reqStatus, reqMsg, from, to, overviewType, resourceName, targetForm string
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
		IsFinished:        reqStatus != "pending",
		IsDownloadable:    isDownloadable,
		ShouldCheckStatus: reqStatus == "pending",
	}
}
