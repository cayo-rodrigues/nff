package utils

import (
	"bytes"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func TrimSpaceBytes(b []byte) string {
	return string(bytes.TrimSpace(b))
}

func TrimSpaceInt(i string) (int, error) {
	return strconv.Atoi(strings.TrimSpace(i))
}

func TrimSpaceFloat64(f string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(f), 64)
}

func TrimSpaceFromBytesToFloat64(f []byte) (float64, error) {
	return TrimSpaceFloat64(TrimSpaceBytes(f))
}

func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func FormatDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("2006-01-02")
}

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("15:04:05")
}

func FormatDateAsBR(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("02/01/2006")
}

func FormatDatetimeAsBR(dt time.Time) string {
	if dt.IsZero() {
		return ""
	}
	return dt.Format("02/01/2006 às 15:04")
}

func IsToday(date string) bool {
	today := FormatDate(time.Now())
	return today == date
}

func IsTodayBR(date string) bool {
	today := FormatDateAsBR(time.Now())
	return today == date
}

func IsYesterday(date string) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return date == FormatDate(yesterday)
}

func IsYesterdayBR(date string) bool {
	yesterday := time.Now().AddDate(0, 0, -1)
	return date == FormatDateAsBR(yesterday)
}

func NDaysBefore(now time.Time, days int) time.Time {
	return now.Add(-time.Duration(days) * 24 * time.Hour)
}

func FormatedNDaysBefore(now time.Time, days int) string {
	return FormatDate(NDaysBefore(now, days))
}

func GetCurrentUserID(c *fiber.Ctx) int {
	userID := c.Locals("UserID").(int)
	return userID
}

func IsAuthenticated(c *fiber.Ctx) bool {
	isAuthenticated, ok := c.Locals("IsAuthenticated").(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
