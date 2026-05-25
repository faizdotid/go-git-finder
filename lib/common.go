package lib

import "strings"

// ParseURL normalizes a raw URL and appends the .git/config path.
func ParseURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
		raw = "http://" + raw
	}
	return strings.TrimSuffix(raw, "/") + "/.git/config"
}
