// internal/utils/time.go
package utils

import (
    "time"
)

const (
    DateFormat     = "2006-01-02"
    TimeFormat     = "15:04:05"
    DateTimeFormat = "2006-01-02 15:04:05"
)

// FormatDate formats time.Time to YYYY-MM-DD
func FormatDate(t time.Time) string {
    return t.Format(DateFormat)
}

// FormatDateTime formats time.Time to YYYY-MM-DD HH:MM:SS
func FormatDateTime(t time.Time) string {
    return t.Format(DateTimeFormat)
}

// ParseDate parses a date string in YYYY-MM-DD format
func ParseDate(date string) (time.Time, error) {
    return time.Parse(DateFormat, date)
}

// ParseDateTime parses a datetime string in YYYY-MM-DD HH:MM:SS format
func ParseDateTime(datetime string) (time.Time, error) {
    return time.Parse(DateTimeFormat, datetime)
}

// BeginningOfDay returns the beginning of the day for the given time
func BeginningOfDay(t time.Time) time.Time {
    year, month, day := t.Date()
    return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day for the given time
func EndOfDay(t time.Time) time.Time {
    year, month, day := t.Date()
    return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// IsWeekend returns true if the given time is a weekend
func IsWeekend(t time.Time) bool {
    return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

// AddWorkdays adds the specified number of workdays to the time
func AddWorkdays(t time.Time, days int) time.Time {
    for days > 0 {
        t = t.AddDate(0, 0, 1)
        if !IsWeekend(t) {
            days--
        }
    }
    return t
}