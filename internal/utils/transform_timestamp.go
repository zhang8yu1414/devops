package utils

import (
	"fmt"
	"time"
)

func TransformTimestamp(timestamp time.Time) string {
	now := time.Now().Unix()
	durationTimestamp := now - timestamp.Unix()
	switch {
	case durationTimestamp < 60:
		return fmt.Sprintf("%ds", durationTimestamp)
	case durationTimestamp >= 60 && durationTimestamp < 60*60:
		minute := durationTimestamp / 60
		second := durationTimestamp % 60
		return fmt.Sprintf("%dm%ds", minute, second)
	case durationTimestamp >= 60*60 && durationTimestamp < 60*60*24:
		hour := durationTimestamp / (60 * 60)
		return fmt.Sprintf("%dh", hour)
	case durationTimestamp >= 60*60*24:
		day := durationTimestamp / (60 * 60 * 24)
		hour := durationTimestamp % (60 * 60 * 24) / (60 * 60)
		return fmt.Sprintf("%dd%dh", day, hour)
	}
	return ""
}
