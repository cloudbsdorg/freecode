package util

import (
	"fmt"
)

func FormatDuration(secs int) string {
	if secs <= 0 {
		return ""
	}
	if secs < 60 {
		return fmt.Sprintf("%ds", secs)
	}
	if secs < 3600 {
		mins := secs / 60
		remaining := secs % 60
		if remaining > 0 {
			return fmt.Sprintf("%dm %ds", mins, remaining)
		}
		return fmt.Sprintf("%dm", mins)
	}
	if secs < 86400 {
		hours := secs / 3600
		remaining := (secs % 3600) / 60
		if remaining > 0 {
			return fmt.Sprintf("%dh %dm", hours, remaining)
		}
		return fmt.Sprintf("%dh", hours)
	}
	if secs < 604800 {
		days := secs / 86400
		if days == 1 {
			return "~1 day"
		}
		return fmt.Sprintf("~%d days", days)
	}
	weeks := secs / 604800
	if weeks == 1 {
		return "~1 week"
	}
	return fmt.Sprintf("~%d weeks", weeks)
}