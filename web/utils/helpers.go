package utils

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
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

func ParseDateAsBR(date string) (time.Time, error) {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load location: %w", err)
	}

	parsedDate, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %w", err)
	}

	return parsedDate, nil
}

func FormatDate(date time.Time) string {
	if date.IsZero() {
		return ""
	}
	return date.Format("2006-01-02")
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
	return dt.Format("02/01/2006 às 15:04:05")
}

func IsTodayBR(date string) bool {
	today := FormatDateAsBR(time.Now())
	return today == date
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

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

type ReqUserData struct {
	ID              int
	IsAuthenticated bool
}

func GetUserData(ctx context.Context) *ReqUserData {
	if userData, ok := ctx.Value("UserData").(*ReqUserData); ok {
		return userData
	}
	return new(ReqUserData)
}

func GetUserID(ctx context.Context) int {
	return GetUserData(ctx).ID
}

func Concurrent(tasks ...func() error) error {
	var wg sync.WaitGroup
	errors := make(chan error, len(tasks))

	for _, task := range tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := task(); err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	if len(errors) > 0 {
		return <-errors
	}

	return nil
}
