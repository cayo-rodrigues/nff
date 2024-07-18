package models

import "time"

type Notifiable interface {
	AsNotification() *Notification
}

type Notification struct {
	ID            int       `json:"id"`
	Status        string    `json:"status"`
	OperationType string    `json:"operation_type"`
	PageEndpoint  string    `json:"page_endpoint"`
	InvoicePDF    string    `json:"invoice_pdf"`
	CreatedAt     time.Time `json:"created_at"`
	UserID        int       `json:"user_id"`
}

func (n *Notification) GetStatus() string {
	return n.Status
}
