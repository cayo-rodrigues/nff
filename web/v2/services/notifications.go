package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/models"
	"github.com/cayo-rodrigues/nff/web/utils"
)

func GetNotifications(ctx context.Context) []*models.Notification {
	redis := database.GetDB().Redis
	userID := utils.GetUserData(ctx).ID

	key := fmt.Sprintf("%d:notification-queue", userID)
	rawItems, err := redis.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		log.Printf("Error reading notification queue for user with id %d. Err: %v\n", userID, err)
		return nil
	}

	notifications := []*models.Notification{}

	for _, rawItem := range rawItems {
		n := new(models.Notification)
		err := json.Unmarshal([]byte(rawItem), &n)
		if err != nil {
			log.Printf("Error unmarshaling notification json data from queue for user with id %d.\nRaw notification data: %s.\nErr: %v\n", userID, rawItem, err)
			continue
		}
		notifications = append(notifications, n)
	}

	return notifications
}

func ClearNotifications(ctx context.Context) error {
	redis := database.GetDB().Redis
	userID := utils.GetUserData(ctx).ID

	key := fmt.Sprintf("%d:notification-queue", userID)

	err := redis.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error clearing notification queue for user with id %d. Err: %v\n", userID, err)
	}

	return err
}

func GetLatestNotification(ctx context.Context) (*models.Notification, int){
	notifications := GetNotifications(ctx)
	notificationsCount := len(notifications)
	if notificationsCount > 0 {
		return notifications[0], notificationsCount
	}
	return new(models.Notification), notificationsCount
}
