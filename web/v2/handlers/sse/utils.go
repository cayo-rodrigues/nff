package sse

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// Message represents an event message sent to the client
type Message struct {
	Event string
	Data  string
}

// sendEvent formats and sends an SSE message to the client
func sendEvent(c *fiber.Ctx, message Message) error {
	// usar sse como trigger para requests!!!
	if _, err := fmt.Fprintf(c, "event: %s\n", message.Event); err != nil {
		return err
	}
	if message.Data != "" {
		if _, err := fmt.Fprintf(c, "data: %s\n\n", message.Data); err != nil {
			return err
		}
	}
	return nil
}
