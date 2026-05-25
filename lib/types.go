package lib

import (
	"net/http"
	"regexp"
)

// Scanner probes URLs for exposed .git/config files.
type Scanner struct {
	urls   []string
	client *http.Client
	writer *ResultWriter
}

// GithubTokenValidator checks extracted GitHub tokens against the API.
type GithubTokenValidator struct {
	all    *ResultWriter
	valid  *ResultWriter
	client *http.Client
}

// GithubRegex matches personal access tokens and OAuth tokens.
var GithubRegex = regexp.MustCompile(`(ghp_|gho_|ghu_|ghs_|ghr_|github_pat)[a-zA-Z0-9_-]{30,}`)
