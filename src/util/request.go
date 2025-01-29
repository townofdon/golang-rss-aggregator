package util

import (
	"net/http"
	"strconv"
)

const DefaultLimit int = 100

func GetLimitOffsetFromUrlQuery(r *http.Request) (limit int, offset int) {
	query := r.URL.Query()
	limitVal := query["limit"]
	offsetVal := query["offset"]

	limit = parseIntWithDefault(limitVal, 100)
	offset = parseIntWithDefault(offsetVal, 0)
	return
}

func parseIntWithDefault(num []string, defaultValue int) int {
	if len(num) == 0 {
		return defaultValue
	}
	if num[0] == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(num[0])
	if err != nil {
		return defaultValue
	}
	return parsed
}
