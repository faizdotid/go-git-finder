package lib

import (
	"fmt"
	"strings"
)

const (
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	BLUE   = "\033[34m"
	YELLOW = "\033[33m"
	WHITE  = "\033[37m"
	RESET  = "\033[0m"
)

func ParseUrl(url string) (string) {
	url = strings.ReplaceAll(url, "\r", "")
	if !strings.Contains(url, "http") {
		url = "https://" + url
	}
	url = strings.TrimSuffix(url, "/")
	return url
}

func RecoverIfError() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf(
				"%sError%s: %s%s%s\n",
				YELLOW,
				BLUE,
				WHITE,
				r,
				RESET,
			)
		}
	}()
}
