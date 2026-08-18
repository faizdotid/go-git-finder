package lib

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// NewScanner creates a Scanner with a buffered result writer.
func NewScanner(urls []string) (*Scanner, error) {
	writer, err := NewResultWriter("results/git-config.txt")
	if err != nil {
		return nil, err
	}
	return &Scanner{
		urls: urls,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		writer: writer,
	}, nil
}

// Close flushes and closes the underlying writer.
func (s *Scanner) Close() error {
	return s.writer.Close()
}

// Run starts a worker pool and processes all URLs.
func (s *Scanner) Run(workers int, g *GithubTokenValidator) error {
	defer s.Close()
	if g != nil {
		defer g.Close()
	}

	if workers < 1 {
		workers = 1
	}

	jobs := make(chan string, workers)
	var wg sync.WaitGroup
	ctx := context.Background()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range jobs {
				s.scanURL(ctx, url, g)
			}
		}()
	}

	for _, url := range s.urls {
		jobs <- url
	}
	close(jobs)
	wg.Wait()
	return nil
}

func (s *Scanner) scanURL(ctx context.Context, url string, g *GithubTokenValidator) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		PrintErr(err)
		return
	}

	resp, err := s.client.Do(req)
	if err != nil {
		PrintErr(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		PrintErr(err)
		return
	}

	statusCode := resp.StatusCode
	bodyStr := string(body)

	if !strings.Contains(bodyStr, "[core]") {
		fmt.Printf("[ %s%d%s ] - %s%s%s\n", Yellow, statusCode, Reset, Blue, url, Reset)
		return
	}

	if err := s.writer.WriteLine(url); err != nil {
		PrintErr(err)
	}

	tokens := GithubRegex.FindAllString(bodyStr, -1)
	if g != nil {
		for _, token := range tokens {
			g.Validate(ctx, url, token)
		}
	}

	if len(tokens) > 0 {
		fmt.Printf("[ %s%d%s ] - [ %s%d tokens%s ] - %s%s%s\n",
			Green, statusCode, Reset,
			Green, len(tokens), Reset,
			Blue, url, Reset)
	} else {
		fmt.Printf("[ %s%d%s ] - %s%s%s\n", Green, statusCode, Reset, Blue, url, Reset)
	}
}
