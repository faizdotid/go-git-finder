package lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// NewGithubTokenValidator creates a validator with buffered writers.
func NewGithubTokenValidator(proxyURL string) (*GithubTokenValidator, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	retryProxy := proxyURL != ""
	if retryProxy {
		parsedProxy, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("parse proxy URL: %w", err)
		}
		if parsedProxy.Scheme == "" || parsedProxy.Host == "" {
			return nil, fmt.Errorf("invalid proxy URL %q", proxyURL)
		}
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.Proxy = http.ProxyURL(parsedProxy)
		client.Transport = transport
	}

	all, err := NewResultWriter("results/tokens.txt")
	if err != nil {
		return nil, err
	}
	valid, err := NewResultWriter("results/valid-tokens.txt")
	if err != nil {
		return nil, errors.Join(err, all.Close())
	}
	return &GithubTokenValidator{
		all:        all,
		valid:      valid,
		client:     client,
		retryProxy: retryProxy,
	}, nil
}

// Close flushes and closes both underlying writers.
func (g *GithubTokenValidator) Close() error {
	return errors.Join(g.all.Close(), g.valid.Close())
}

// Validate checks a GitHub token against the GitHub API.
func (g *GithubTokenValidator) Validate(ctx context.Context, url, token string) error {
	if err := g.all.WriteLine(token); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		PrintErr(err)
		return nil
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	attempts := 1
	if g.retryProxy {
		attempts = 4
	}
	var resp *http.Response
	for attempt := 0; attempt < attempts; attempt++ {
		resp, err = g.client.Do(req.Clone(ctx))
		if err == nil {
			break
		}
	}
	if err != nil {
		PrintErr(err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var user struct {
			Login string `json:"login"`
		}
		body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		if err != nil {
			PrintErr(err)
			return nil
		}
		if err := json.Unmarshal(body, &user); err != nil {
			PrintErr(err)
			return nil
		}
		fmt.Printf("[ %sVALID%s ] - %s%s%s ( %s%s%s )\n", Green, Reset, Blue, token, Reset, Green, user.Login, Reset)
		if err := g.valid.WriteLine(url + "|" + token + ":" + user.Login); err != nil {
			return err
		}
	}
	return nil
}
