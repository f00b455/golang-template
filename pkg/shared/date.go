package shared

import "time"

// FormatDate formats a date to ISO date string (YYYY-MM-DD).
func FormatDate(date time.Time) string {
	return date.Format("2006-01-02")
}
