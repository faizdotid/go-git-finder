package lib

import (
	"net/http"
	"os"
	"regexp"
)

// scanner struct
type Scanner struct {
	urls []string
	c    *http.Client
	f    *os.File
}

// github token check struct
type GithubTokenValidator struct {
	all   *os.File
	valid *os.File
	c     *http.Client
}

var (
	// regexp for github token
	GithubRegex = regexp.MustCompile(`(ghp_|gho_|ghu_|ghs_|ghr_|github_pat)[a-zA-Z0-9]{30,}`)
)
