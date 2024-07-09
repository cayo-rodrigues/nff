package models

import "time"

type Notification interface {
	GetStatus() string
	GetOperationType() string
	GetCreatedAt() time.Time
	GetID() int
	GetPageEndpoint() string
}

func JsonSerializableNotification(n Notification, userID int) map[string]any {
	operation := map[string]interface{}{
		"id":             n.GetID(),
		"status":         n.GetStatus(),
		"operation_type": n.GetOperationType(),
		"page_endpoint":  n.GetPageEndpoint(),
		"created_at":     n.GetCreatedAt(),
		"user_id":        userID,
	}
	return operation
}
