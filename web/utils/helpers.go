package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
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

func ParseDate(date string, tz string, layout string) (time.Time, error) {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to load location: %w", err)
	}

	parsedDate, err := time.ParseInLocation(layout, date, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date: %w", err)
	}

	return parsedDate, nil
}

func ParseDateWithBRTZ(date string) (time.Time, error) {
	return ParseDate(date, "America/Sao_Paulo", "2006-01-02")
}

func ParseDateWithBRTZAndBRLayout(date string) (time.Time, error) {
	return ParseDate(date, "America/Sao_Paulo", "02/01/2006")
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

var BrazilTZ, _ = time.LoadLocation("America/Sao_Paulo")

func NowBR() time.Time {
	return time.Now().In(BrazilTZ)
}

func InBrazilTZ(t time.Time) time.Time {
	return t.In(BrazilTZ)
}

func IsTodayBR(date string) bool {
	today := FormatDateAsBR(NowBR())
	return today == date
}

func GetWeekDay(date time.Time) string {
	weekDays := map[time.Weekday]string{
		time.Sunday:    "domingo",
		time.Monday:    "segunda-feira",
		time.Tuesday:   "terça-feira",
		time.Wednesday: "quarta-feira",
		time.Thursday:  "quinta-feira",
		time.Friday:    "sexta-feira",
		time.Saturday:  "sábado",
	}

	return weekDays[date.Weekday()]
}

func GetWeekDayFromString(date string) string {
	dt, err := ParseDateWithBRTZAndBRLayout(date)
	if err != nil {
		log.Printf("[GetWeekDayFromString] Failed to parse date %s: %s\n", date, err)
		return ""
	}

	return GetWeekDay(dt)
}

func IsYesterdayBR(date string) bool {
	yesterday := NowBR().AddDate(0, 0, -1)
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

type UserCtxData struct {
	ID                   int
	IsAuthenticated      bool
	IsBlocked            bool
	HasChosenPaymentPlan bool
}

func GetUserData(ctx context.Context) *UserCtxData {
	if userData, ok := ctx.Value("UserData").(*UserCtxData); ok {
		return userData
	}
	return new(UserCtxData)
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
