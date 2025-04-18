package middlewares

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/cayo-rodrigues/nff/web/database"
	"github.com/cayo-rodrigues/nff/web/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

var editEntityPageRegex = regexp.MustCompile("^/entities/([1-9]\\d*)(/?)$")

func CacheManagementMiddleware(c *fiber.Ctx) error {
	if shouldSkip(c) {
		return c.Next()
	}

	if fiber.IsMethodSafe(c.Method()) {
		return useOrSetCache(c)
	}

	return callNextAndClearCache(c)
}

func shouldSkip(c *fiber.Ctx) bool {
	reqPath := c.Path()
	return strings.HasPrefix(reqPath, "/sse") || editEntityPageRegex.MatchString(reqPath)
}

func useOrSetCache(c *fiber.Ctx) error {
	bodyKey, headersKey := genKeys(c)

	db := database.GetDB()

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

	db := database.GetDB()

	err := db.Redis.Set(c.Context(), bodyKey, bodyBytes, time.Hour).Err()
	if err != nil {
		log.Println("Error trying to set response body cache:", err)
	}

	err = db.Redis.Set(c.Context(), headersKey, headersBytes, time.Hour).Err()
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
	if namespace == "entities" {
		namespace = "*" // all other routes use entities, so we gotta clear all pages cache in this case
	}

	db := database.GetDB()

	cacheStatus := db.Redis.ClearCache(c.Context(), userID, namespace)

	c.Response().Header.Set("X-Cache", cacheStatus)

	return nil
}

func getKeyFactors(c *fiber.Ctx) (int, string, string) {
	userID := utils.GetUserData(c.Context()).ID
	route := c.OriginalURL()

	var namespace string

	switch {
	case strings.HasPrefix(route, "/invoices/cancel"):
		namespace = "invoice-cancel"
	case strings.HasPrefix(route, "/invoices/print"):
		namespace = "invoice-print"
	case strings.HasPrefix(route, "/invoices"):
		namespace = "invoice-issue"
	case strings.HasPrefix(route, "/reauthenticate"):
		namespace = "*"
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
