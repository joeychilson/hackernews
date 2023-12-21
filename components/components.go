package components

import (
	"fmt"
	"time"
)

func TimeAgo(postTime int64) string {
	now := time.Now()
	postTimeUTC := time.Unix(postTime, 0)
	duration := now.Sub(postTimeUTC)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		if duration/time.Minute == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", duration/time.Minute)
	} else if duration < time.Hour*24 {
		if duration/time.Hour == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", duration/time.Hour)
	} else if duration < time.Hour*24*31 {
		if duration/(time.Hour*24) == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", duration/(time.Hour*24))
	} else if duration < time.Hour*24*365 {
		months := duration / (time.Hour * 24 * 30)
		if months <= 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	years := duration / (time.Hour * 24 * 365)
	if years <= 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}
