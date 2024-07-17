package utils

import "time"

func IsTodayHiddenFromFilters(filters map[string]string) bool {
	filtersExcludeToday := false

	fromDate, err := ParseDateAsBR(filters["from_date"])
	if err != nil {
		filtersExcludeToday = true
	}
	toDate, err := ParseDateAsBR(filters["to_date"])
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
