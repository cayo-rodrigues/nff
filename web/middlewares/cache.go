package middlewares

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/db"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func CacheMiddleware(c *fiber.Ctx) error {
	if fiber.IsMethodSafe(c.Method()) {
		return useOrSetCache(c)
	}

	return callNextAndClearCache(c)
}

func useOrSetCache(c *fiber.Ctx) error {
	bodyKey, headersKey := genKeys(c)

	cachedBody, err := db.Redis.Get(c.Context(), bodyKey).Bytes()
	if err == redis.Nil {
		return callNextAndSetCache(c)
	}
	if err != nil {
		log.Println("Error trying to get cached response body:", err)
		return callNextAndSetCache(c)
	}

	cachedHeaders, err := db.Redis.Get(c.Context(), headersKey).Bytes()
	if err == redis.Nil {
		return callNextAndSetCache(c)
	}
	if err != nil {
		log.Println("Error trying to get cached response headers:", err)
		return callNextAndSetCache(c)
	}

	response := c.Response()
	response.Header.Read(bufio.NewReader(bytes.NewReader(cachedHeaders)))
	response.Header.Set("X-Cache", "hit")
	response.Header.Del("Connection")
	response.SetBody(cachedBody)

	return nil
}

func callNextAndSetCache(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		return err
	}
	response := c.Response()

	bodyKey, headersKey := genKeys(c)
	bodyBytes := response.Body()
	headersBytes := response.Header.Header()

	err := db.Redis.Set(c.Context(), bodyKey, bodyBytes, time.Minute).Err()
	if err != nil {
		log.Println("Error trying to set response body cache:", err)
	}

	err = db.Redis.Set(c.Context(), headersKey, headersBytes, time.Minute).Err()
	if err != nil {
		log.Println("Error trying to set response headers cache:", err)
	}

	response.Header.Set("X-Cache", "miss")

	return nil
}

func callNextAndClearCache(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		return err
	}

	userID, _, namespace := getKeyFactors(c)
	cacheStatus := utils.ClearCache(c.Context(), userID, namespace)

	c.Response().Header.Set("X-Cache", cacheStatus)

	return nil
}

func getKeyFactors(c *fiber.Ctx) (int, string, string) {
	userID := c.Locals("UserID").(int)
	route := c.OriginalURL()

	var namespace string

	switch {
	case strings.HasPrefix(route, "/invoices/cancel"):
		namespace = "invoices-cancel"
	case strings.HasPrefix(route, "/invoices/print"):
		namespace = "invoices-print"
	case strings.HasPrefix(route, "/invoices"):
		namespace = "invoices"
	default:
		namespace = strings.Split(route, "/")[1]
	}

	return userID, route, namespace
}

func genKeys(c *fiber.Ctx) (string, string) {
	userID, route, namespace := getKeyFactors(c)

	bodyKey := fmt.Sprintf("body:%v:%v:%v", userID, route, namespace)
	headersKey := fmt.Sprintf("headers:%v:%v:%v", userID, route, namespace)

	return bodyKey, headersKey
}
