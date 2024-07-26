package siare

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
)

func formatErrResponseField(fieldErrs []string, prefix string) string {
	var builder strings.Builder

	if fieldErrs != nil {
		builder.WriteString(prefix)
		for _, fieldName := range fieldErrs {
			builder.WriteString(" - ")
			builder.WriteString(fieldName)
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func formatErrResponse(response *SSApiErrorResponse) string {
	var clientErrs strings.Builder

	clientErrs.WriteString(formatErrResponseField(response.InvalidFields, InvalidFieldsPrefix))
	clientErrs.WriteString(formatErrResponseField(response.MissingFields, MissingFieldsPrefix))
	if response.Entity != nil {
		clientErrs.WriteString(formatErrResponseField(response.Entity.InvalidFields, EntityInvalidFieldsPrefix))
		clientErrs.WriteString(formatErrResponseField(response.Entity.MissingFields, EntityMissingFieldsPrefix))
	}
	if response.Sender != nil {
		clientErrs.WriteString(formatErrResponseField(response.Sender.InvalidFields, SenderInvalidFieldsPrefix))
		clientErrs.WriteString(formatErrResponseField(response.Sender.MissingFields, SenderMissingFieldsPrefix))
	}
	if response.Recipient != nil {
		clientErrs.WriteString(formatErrResponseField(response.Recipient.InvalidFields, RecipientInvalidFieldsPrefix))
		clientErrs.WriteString(formatErrResponseField(response.Recipient.MissingFields, RecipientMissingFieldsPrefix))
	}

	return clientErrs.String()
}

func requestErrLog(endpoint, resourceName string, resourceID int, err error) {
	log.Printf("Something went wrong with the request at %s for %s with id %v: %v\n", endpoint, resourceName, resourceID, err)
}

func jsonUnmarchalErrLog(endpoint, resourceName string, resourceID int, err error) {
	log.Printf("Something went wrong trying to unmarshal ss-api %s json response. %s with id %v will be on 'pending' state for ever: %v\n", endpoint, resourceName, resourceID, err)
}

func notifyOperationResult(ctx context.Context, redisClient *database.Redis, userID int, resourceID int) {
	channel := fmt.Sprintf("%d:operation-finished", userID)
	redisClient.Publish(ctx, channel, resourceID)
}

func notifySingleOperationResult(ctx context.Context, redisClient *database.Redis, userID int, resourceID int, operationName string) {
	channel := fmt.Sprintf("%d:%s-operation-finished", userID, operationName)
	redisClient.Publish(ctx, channel, resourceID)
}

func enqueueNotification(ctx context.Context, redisClient *database.Redis, userID int, n *models.Notification) {
	queueName := fmt.Sprintf("%d:notification-queue", userID)

	notificationJSON, err := json.Marshal(n)
	if err != nil {
		log.Println("Failed to marshal notification:", err)
		return
	}

	err = redisClient.LPush(ctx, queueName, notificationJSON).Err()
	if err != nil {
		log.Println("Failed to enqueue notification:", err)
	}
}

func finishOperation(ctx context.Context, redisClient *database.Redis, resourceName string, userID int, n models.Notifiable) {
	notification := n.AsNotification()

	enqueueNotification(ctx, redisClient, userID, notification)
	redisClient.ClearCache(ctx, userID, resourceName)
	notifyOperationResult(ctx, redisClient, userID, notification.ID)
}
