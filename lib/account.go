package lib

import (
	"net/http"
	"os"
)

func NewGithubTokenValidator() *GithubTokenValidator {
	all, err := os.OpenFile("results/tokens.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	valid, err := os.OpenFile("results/valid-tokens.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	return &GithubTokenValidator{
		all:   all,
		valid: valid,
	}
}

func (g *GithubTokenValidator) Validate(token string) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		PrintErr(err)
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := g.c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		_, err = g.valid.WriteString(token + "\n")
		if err != nil {
			PrintErr(err)
		}
	}

	_, err = g.all.WriteString(token + "\n")
	if err != nil {
		PrintErr(err)
	}
}

// closing the files
func (g *GithubTokenValidator) Close() {
	g.all.Close()
	g.valid.Close()
}
