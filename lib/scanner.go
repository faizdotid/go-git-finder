package lib

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func NewScanner(c []string) *Scanner {
	file, err := os.OpenFile(
		"results/git-config.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic(err)
	}
	return &Scanner{
		urls: c,
		c: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
		f: file,
	}
}

func (s *Scanner) Scan(url string, wg *sync.WaitGroup, g *GithubTokenValidator) {
	defer wg.Done()

	resp, err := s.c.Get(url)
	if err != nil {
		PrintErr(err)
		return
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		PrintErr(err)
	}
	sCode := strconv.Itoa(resp.StatusCode)

	if !strings.Contains(
		string(buf),
		"[core]",
	) {
		//PRINT NOT VALID .GIT
		fmt.Printf(
			"[ %s%s%s ] - %s%s%s\n",
			Yellow,
			sCode,
			Reset,
			Blue,
			url,
			Reset,
		)
		return
	}
	_, err = s.f.WriteString(url + "\n")
	if err != nil {
		PrintErr(err)
	}
	tokens := GithubRegex.FindAllString(string(buf), -1)
	for _, token := range tokens {
		g.Validate(token)
	}

	if len(tokens) > 0 {
		fmt.Printf(
			"[ %s%s%s ] - [ %s%d%s ] - %s%s%s\n",
			Green,
			sCode,
			Reset,
			Green,
			len(tokens),
			Reset,
			Blue,
			url,
			Reset,
		)
	} else {
		fmt.Printf(
			"[ %s%s%s ] - %s%s%s\n",
			Green,
			sCode,
			Reset,
			Blue,
			url,
			Reset,
		)
	}
}

func (s *Scanner) Run(thread int) {
	// make sure to close the file
	defer s.f.Close()

	// make sure to close the file
	g := NewGithubTokenValidator()
	defer g.Close()

	// create a wait group
	var wg sync.WaitGroup

	threadChannel := make(chan struct{}, thread)
	for _, url := range s.urls {
		wg.Add(1)
		threadChannel <- struct{}{}
		go func(url string) {
			s.Scan(url, &wg, g)
			<-threadChannel
		}(url)
	}
	wg.Wait()
}
