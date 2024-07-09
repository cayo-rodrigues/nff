package services

import (
	"time"

	// "github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
)

func GetNotificationListItems() []models.Notification {
	// items will be stored in redis
	// redis := database.GetDB().Redis

	m := models.NewMetrics()
	m.ReqStatus = "success"
	m.CreatedAt = time.Now()
	i := models.NewInvoice()
	i.ReqStatus = "error"
	i.CreatedAt = time.Now()
	c := models.NewInvoiceCancel()
	c.ReqStatus = "warning"
	c.CreatedAt = time.Now()
	p := models.NewInvoicePrint()
	p.ReqStatus = "success"
	p.CreatedAt = time.Now()
	return []models.Notification{m, i, c, p}
}
