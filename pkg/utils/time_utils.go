package utils

import "time"

func GetDatesExclusive(start time.Time, end time.Time) []time.Time {
	if start.After(end) {
		return nil
	}

	days := []time.Time{}
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		days = append(days, Date(d))
	}
	return days
}

func Date(timestamp time.Time) time.Time {
	return time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, time.UTC)
}
