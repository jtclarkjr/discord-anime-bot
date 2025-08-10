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

// FormatDate formats a date to user-friendly format (Discord timestamp)
// Format: <t:timestamp:D> - Shows date only (e.g., "December 25, 2023")
func FormatDate(date time.Time) string {
	return fmt.Sprintf("<t:%d:D>", date.Unix())
}

// FormatTime formats time to user-friendly format (Discord timestamp)
// Format: <t:timestamp:t> - Shows time only (e.g., "3:30 PM")
func FormatTime(date time.Time) string {
	return fmt.Sprintf("<t:%d:t>", date.Unix())
}

// FormatAirDate formats date and time for air date display (Discord timestamp)
// Format: <t:timestamp:F> - Full date and time (e.g., "Monday, December 25, 2023 3:30 PM")
func FormatAirDate(date time.Time) string {
	return fmt.Sprintf("<t:%d:F>", date.Unix())
}

// FormatCompactDateTime formats compact date and time for lists (Discord timestamp)
// Format: <t:timestamp:f> - Short date and time (e.g., "December 25, 2023 3:30 PM")
func FormatCompactDateTime(date time.Time) string {
	return fmt.Sprintf("<t:%d:f>", date.Unix())
}

// FormatRelativeTimestamp formats a time as Discord relative timestamp
// Format: <t:timestamp:R> - Relative time (e.g., "in 2 hours", "3 days ago")
func FormatRelativeTimestamp(date time.Time) string {
	return fmt.Sprintf("<t:%d:R>", date.Unix())
}

// FormatRelativeTime formats a time relative to now (e.g., "in 2 hours", "tomorrow")
func FormatRelativeTime(t time.Time) string {
	now := time.Now()
	duration := t.Sub(now)

	if duration < 0 {
		return "already aired"
	}

	// Less than a minute
	if duration < time.Minute {
		return "in less than a minute"
	}

	// Less than an hour
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes == 1 {
			return "in 1 minute"
		}
		return fmt.Sprintf("in %d minutes", minutes)
	}

	// Less than a day
	if duration < 24*time.Hour {
		hours := int(duration.Hours())
		if hours == 1 {
			return "in 1 hour"
		}
		return fmt.Sprintf("in %d hours", hours)
	}

	// Less than a week
	if duration < 7*24*time.Hour {
		days := int(duration.Hours() / 24)
		if days == 1 {
			return "tomorrow"
		}
		return fmt.Sprintf("in %d days", days)
	}

	// More than a week
	weeks := int(duration.Hours() / (24 * 7))
	if weeks == 1 {
		return "in 1 week"
	}
	return fmt.Sprintf("in %d weeks", weeks)
}
