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
	_, err := fmt.Fprintf(c, "event: %s\ndata: %s\n\n", message.Event, message.Data)
	return err
}
