package lib

import (
	"strings"
)

func ParseURL(rawUrl string) string {
	url := strings.TrimSpace(rawUrl) // remove leading and trailing spaces

	if !strings.Contains(url, "http") {
		url = "http://" + url
	}

	url = strings.TrimSuffix(url, "/")

	return url + "/.git/config"
}
