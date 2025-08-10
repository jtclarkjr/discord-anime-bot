package utils

import (
	"fmt"
	"time"
)

// FormatCountdown formats countdown time from seconds to a human readable string
func FormatCountdown(timeUntil int) string {
	days := timeUntil / (24 * 60 * 60)
	hours := (timeUntil % (24 * 60 * 60)) / (60 * 60)
	minutes := (timeUntil % (60 * 60)) / 60

	var timeString string
	if days > 0 {
		if days > 1 {
			timeString += fmt.Sprintf("%d days ", days)
		} else {
			timeString += fmt.Sprintf("%d day ", days)
		}
	}
	if hours > 0 {
		if hours > 1 {
			timeString += fmt.Sprintf("%d hours ", hours)
		} else {
			timeString += fmt.Sprintf("%d hour ", hours)
		}
	}
	if minutes > 0 {
		if minutes > 1 {
			timeString += fmt.Sprintf("%d minutes", minutes)
		} else {
			timeString += fmt.Sprintf("%d minute", minutes)
		}
	}

	if timeString == "" {
		return "less than a minute"
	}

	return timeString
}

// FormatDate formats a date to user-friendly format
func FormatDate(date time.Time) string {
	return date.Format("Monday, January 2, 2006")
}

// FormatTime formats time to user-friendly format (12-hour)
func FormatTime(date time.Time) string {
	return date.Format("3:04 PM")
}

// FormatAirDate formats date and time for air date display
func FormatAirDate(date time.Time) string {
	return fmt.Sprintf("%s at %s", FormatDate(date), FormatTime(date))
}

// FormatCompactDateTime formats compact date and time for lists
func FormatCompactDateTime(date time.Time) string {
	return date.Format("Jan 2, 3:04 PM")
}
