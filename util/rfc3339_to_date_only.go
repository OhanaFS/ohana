package util

import "time"

// Converts a RFC3339 date to a date only (no time) with UTC.
func Rfc3339ToDateOnly(rfc3339 string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, rfc3339)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC), nil
}
