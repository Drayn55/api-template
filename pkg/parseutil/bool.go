package parseutil

import (
	"strconv"
	"strings"
)

func ParseBool(val string) bool {
	if val == "" {
		return false
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return val == "1"
	}
	return b
}

func ParseInt(val string) int64 {
	v := strings.TrimSpace(val)
	if v == "" {
		return 0
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

