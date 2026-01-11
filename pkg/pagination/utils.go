package pagination

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CursorSeparator is used to separate timestamp and ID in cursor encoding
// Using "||" to avoid confusion with timestamp format (RFC3339) which contains colons
const CursorSeparator = "||"

// EncodeCursor creates a base64-encoded cursor from timestamp and ID
// Format: "timestamp||id" -> base64("2026-01-10T10:00:00Z||12345")
func EncodeCursor(timestamp time.Time, id int64) string {
	rawCursor := fmt.Sprintf("%s%s%d", timestamp.Format(time.RFC3339), CursorSeparator, id)
	return base64.URLEncoding.EncodeToString([]byte(rawCursor))
}

// DecodeCursor parses a base64-encoded cursor back to timestamp and ID
// Returns empty string and 0 if cursor is invalid (graceful degradation)
func DecodeCursor(encodedCursor string) (string, int64, error) {
	if encodedCursor == "" {
		return "", 0, nil
	}

	decoded, err := base64.URLEncoding.DecodeString(encodedCursor)
	if err != nil {
		return "", 0, fmt.Errorf("failed to decode cursor: %w", err)
	}

	// split into timestamp and ID using CursorSeparator
	parts := strings.Split(string(decoded), CursorSeparator)
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid cursor format")
	}

	timestampStr := parts[0]
	idStr := parts[1]

	// validate timestamp format
	_, err = time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return "", 0, fmt.Errorf("invalid timestamp in cursor: %w", err)
	}

	// parse ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return "", 0, fmt.Errorf("invalid ID in cursor: %w", err)
	}

	return timestampStr, id, nil
}

// ReverseOrder returns the opposite of the given order string ("asc" <-> "desc")
func ReverseOrder(order string) string {
	if order == "asc" {
		return "desc"
	}
	return "asc"
}

// ReverseSlice reverses a slice of interface{} items in-place and returns it
func ReverseSlice(items []interface{}) []interface{} {
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	return items
}
