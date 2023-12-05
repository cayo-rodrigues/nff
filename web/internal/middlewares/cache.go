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

	err := db.Redis.Set(c.Context(), bodyKey, response.Body(), time.Minute).Err()
	if err != nil {
		log.Println("Error trying to set response body cache:", err)
	}

	err = db.Redis.Set(c.Context(), headersKey, response.Header.Header(), time.Minute).Err()
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
	keysPattern := fmt.Sprintf("*:%v:*:%v", userID, namespace)
	cacheStatus := "cleared"

	keys, err := db.Redis.Keys(c.Context(), keysPattern).Result()
	if err != nil {
		log.Println("Error geting cache keys to clear:", err)
		cacheStatus = "stale"
	}

	if keys != nil && len(keys) > 0 {
		err := db.Redis.Del(c.Context(), keys...).Err()
		if err != nil {
			log.Println("Error clearing cache keys:", err)
			cacheStatus = "stale"
		}
	}

	c.Response().Header.Set("X-Cache", cacheStatus)
	return nil
}

func getKeyFactors(c *fiber.Ctx) (int, string, string) {
	userID := c.Locals("UserID").(int)
	route := c.OriginalURL()
	namespace := strings.Split(route, "/")[1]

	return userID, route, namespace
}

func genKeys(c *fiber.Ctx) (string, string) {
	userID, route, namespace := getKeyFactors(c)

	bodyKey := fmt.Sprintf("body:%v:%v:%v", userID, route, namespace)
	headersKey := fmt.Sprintf("headers:%v:%v:%v", userID, route, namespace)

	return bodyKey, headersKey
}
