package middlewares

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func CacheMiddleware(c *fiber.Ctx) error {
	if fiber.IsMethodSafe(c.Method()) {
		return useOrSetCache(c)
	}

	return clearCache(c)
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

func clearCache(c *fiber.Ctx) error {
	// delete all cache keys in the current namespace, not just based on this route

	return c.Next()
}

func genKeys(c *fiber.Ctx) (string, string) {
	userID := c.Locals("UserID").(int)
	route := c.OriginalURL()
	namespace := strings.Split(route, "/")[1]

	bodyKey := fmt.Sprintf("body:%v:%v:%v", userID, route, namespace)
	headersKey := fmt.Sprintf("headers:%v:%v:%v", userID, route, namespace)

	return bodyKey, headersKey
}

func callNextAndSetCache(c *fiber.Ctx) error {
	if err := c.Next(); err != nil {
		return err
	}
	response := c.Response()

	bodyKey, headersKey := genKeys(c)

	err := db.Redis.Set(c.Context(), bodyKey, response.Body(), time.Minute).Err()
	if err != nil {
		log.Println("Error trying to set response body cache:", err)
		return err
	}

	err = db.Redis.Set(c.Context(), headersKey, response.Header.Header(), time.Minute).Err()
	if err != nil {
		log.Println("Error trying to set response headers cache:", err)
		return err
	}

	response.Header.Set("X-Cache", "miss")
	return nil
}
