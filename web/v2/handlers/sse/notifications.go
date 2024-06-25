package sse

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func NotifyOperationsResults(c *fiber.Ctx) error {
	// 1. subscribe em todos os tópicos das operações (invoice-issuing, invoice-cancel, invoice-print, metrics)
	// 2. ao receber uma notificação, enviar resposta
	// 3. defer unsubscribe (se for o caso)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	timeChan := make(chan time.Time)

	go func() {
		time.Sleep(2 * time.Second)
		timeChan<- time.Now()
	}()

	for {
		select {
		case t := <- timeChan:
			message := Message{
				Event: "time",
				Data:  t.Format("2006-01-02 15:04:05"),
			}
			if err := sendEvent(c, message); err != nil {
				log.Println("Error sending event:", err)
				return err
			}
			return nil
		case <-c.Context().Done():
			log.Println("Client closed connection")
			return nil
		}
	}
}
