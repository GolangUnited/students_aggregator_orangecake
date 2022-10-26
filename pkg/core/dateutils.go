package core

import (
	"errors"
	"time"
)

func NormalizeDate(aDate time.Time) time.Time {
	lDate := aDate.UTC()
	return time.Date(lDate.Year(), lDate.Month(), lDate.Day(), 0, 0, 0, 0, time.UTC)
}

func ParseDate(aLayout, aValue string) (time.Time, error) {
	if len(aValue) <= 0 {
		return NormalizeDate(time.Now()), errors.New("empty Date")
	}

	CI Test test compiler error

	// try parse date or use Now if failed
	lDate, lErr := time.Parse(aLayout, aValue)
	if lErr != nil {
		return NormalizeDate(time.Now()), errors.New("invalid Date format")
	}

	return NormalizeDate(lDate), nil
}
