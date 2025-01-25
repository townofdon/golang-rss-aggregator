package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Extract API key from a request header
// Example:
// Authorization: ApiKey <VALUE>
func ParseApiKeyFromHeader(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("missing `Authorization` header")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", fmt.Errorf("malformed `Authorization` header: %v", val)
	}

	authType := vals[0]
	authVal := vals[1]

	switch authType {
	case "ApiKey":
		return authVal, nil
	default:
		return "", fmt.Errorf("auth type not supported: %v", authType)
	}
}
