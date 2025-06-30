package extractor

import (
	"strconv"
	"strings"
)

func TimeAgoToMinutes(timeAgo string) int {
	timeAgo = strings.TrimSpace(timeAgo)

	if strings.HasSuffix(timeAgo, "mo") {
		months, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "mo"))
		return months * 30 * 24 * 60
	} else if strings.HasSuffix(timeAgo, "w") {
		weeks, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "w"))
		return weeks * 7 * 24 * 60
	} else if strings.HasSuffix(timeAgo, "d") {
		days, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "d"))
		return days * 24 * 60
	} else if strings.HasSuffix(timeAgo, "h") {
		hours, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "h"))
		return hours * 60
	} else if strings.HasSuffix(timeAgo, "m") {
		minutes, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "m"))
		return minutes
	} else if strings.HasSuffix(timeAgo, "y") {
		years, _ := strconv.Atoi(strings.TrimSuffix(timeAgo, "y"))
		return years * 365 * 24 * 60 // Approximate a year as 365 days
	}

	return 0 // fallback if format is unknown
}
