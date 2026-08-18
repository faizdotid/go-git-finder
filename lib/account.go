package lib

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// NewGithubTokenValidator creates a validator with buffered writers.
func NewGithubTokenValidator() (*GithubTokenValidator, error) {
	all, err := NewResultWriter("results/tokens.txt")
	if err != nil {
		return nil, err
	}
	valid, err := NewResultWriter("results/valid-tokens.txt")
	if err != nil {
		all.Close()
		return nil, err
	}
	return &GithubTokenValidator{
		all:   all,
		valid: valid,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

// Close flushes and closes both underlying writers.
func (g *GithubTokenValidator) Close() error {
	return errors.Join(g.all.Close(), g.valid.Close())
}

// Validate checks a GitHub token against the GitHub API.
func (g *GithubTokenValidator) Validate(ctx context.Context, url, token string) {
	if err := g.all.WriteLine(token); err != nil {
		PrintErr(err)
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		PrintErr(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := g.client.Do(req)
	if err != nil {
		PrintErr(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("[ %sVALID%s ] - %s%s%s\n", Green, Reset, Blue, token, Reset)
		if err := g.valid.WriteLine(url + "|" + token); err != nil {
			PrintErr(err)
		}
	}
}
