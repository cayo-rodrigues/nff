package sse

import (
	"fmt"
	"log"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
)

func NotifyOperationsResults(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	redis := database.GetDB().Redis

	userID := utils.GetUserData(c.Context()).ID

	channelPatterns := []string{"invoice-issue", "invoice-cancel", "invoice-print", "metrics"}
	userChannels := []string{}
	for _, pattern := range channelPatterns {
		userChannels = append(userChannels, fmt.Sprintf("%d:%s-finished", userID, pattern))
	}
	sub := redis.Subscribe(c.Context(), userChannels...)

	defer sub.Unsubscribe(c.Context())
	defer sub.Close()

	_, err := sub.Receive(c.Context())
	if err != nil {
		log.Println("Could not subscribe: ", err)
		return err
	}

	ch := sub.Channel()

	for {
		select {
		case msg := <-ch:
			message := Message{
				Event: msg.Channel,
			}
			if err := sendEvent(c, message); err != nil {
				log.Println("Error sending event: ", err)
				return err
			}
			return nil
		case <-c.Context().Done():
			log.Println("Client closed connection")
			return nil
		}
	}
}
