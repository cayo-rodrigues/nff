package sse

import (
	"log"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/gofiber/fiber/v2"
)

func NotifyOperationsResults(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	db := database.GetDB()

	sub := db.Redis.Subscribe(c.Context(), "invoice-issue", "invoice-cancel", "invoice-print", "metrics")

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
				Event: msg.Payload,
				Data:  msg.Payload,
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
