package parseutil

import (
	"fmt"
	"strings"
	"time"
)

var timeLayouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05 -0700 MST",
	"2006-01-02",
}

func ParseTime(raw string) (time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return time.Time{}, fmt.Errorf("time value is required")
	}

	for _, layout := range timeLayouts {
		if t, err := time.Parse(layout, raw); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format: %q", raw)
}
