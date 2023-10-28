package models

type RequestCard struct {
	From              string
	To                string
	ToPrefix          string
	FeedbackTextColor string
	ReqStatus         string
	ResourceName      string
	TargetForm        string
	OverviewType      string
	HasItems          bool
	IsFinished        bool
	IsDownloadable    bool
	ShouldCheckStatus bool
}

// req may be either *Invoice or *InvoiceCancel
func NewRequestCard(req any) *RequestCard {
	var prefix, reqStatus, from, to, overviewType, resourceName, targetForm string
	var hasItems, isDownloadable bool

	switch r := req.(type) {
	case *Invoice:
		from = r.Sender.Name
		to = r.Recipient.Name
		reqStatus = r.ReqStatus
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
		resourceName = "invoices/cancel"
		targetForm = "#invoice-cancel-form"
		overviewType = "canceling"
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
		ResourceName:      resourceName,
		TargetForm:        targetForm,
		OverviewType:      overviewType,
		HasItems:          hasItems,
		IsFinished:        reqStatus != "pending",
		IsDownloadable:    isDownloadable,
		ShouldCheckStatus: reqStatus == "pending",
	}
}
