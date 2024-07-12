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

	userChannel := fmt.Sprintf("%d:operation-finished", userID)
	sub := redis.Subscribe(c.Context(), userChannel)

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
