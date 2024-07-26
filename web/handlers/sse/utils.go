package sse

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Message struct {
	Event string
	Data  string
}

func sendEvent(c *fiber.Ctx, message Message) error {
	if _, err := fmt.Fprintf(c, "event: %s\n", message.Event); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(c, "data: %s\n\n", message.Data); err != nil {
		return err
	}
	return nil
}
