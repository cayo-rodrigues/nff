package utils

import "time"

func FiltersExcludeToday(filters map[string]string) bool {
	filtersExcludeToday := false

	fromDate, err := ParseDate(filters["from_date"])
	if err != nil {
		filtersExcludeToday = true
	}
	toDate, err := ParseDate(filters["to_date"])
	if err != nil {
		filtersExcludeToday = true
	}

	today := time.Now().Truncate(time.Hour * 24)
	fromDate = fromDate.Truncate(time.Hour * 24)
	toDate = toDate.Truncate(time.Hour * 24)
	if fromDate.After(today) || toDate.Before(today) {
		filtersExcludeToday = true
	}

	return filtersExcludeToday
}
