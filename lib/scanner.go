package lib

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Scanner struct {
	Urls   []string
	Client *http.Client
}

func NewScanner(urls []string) *Scanner {
	return &Scanner{
		Urls: urls,
		Client: &http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

func (s *Scanner) Scan(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	url = ParseUrl(url)
	response, err := s.Client.Get(url + "/.git/HEAD")

	if err != nil {
		//PRINT ERROR URL
		fmt.Printf(
			"%sUrl: %s%s%s - %s%s%s\n",
			WHITE,
			BLUE,
			url,
			RESET,
			YELLOW,
			err.Error(),
			RESET,
		)
		return
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf(
			"%sError: %s%s%s\n",
			WHITE,
			YELLOW,
			err.Error(),
			RESET,
		)
	}
	if !strings.Contains(
		string(responseBody),
		"refs/heads",
	) {
		//PRINT NOT VALID .GIT
		fmt.Printf(
			"%sUrl: %s%s%s - %sNot Valid Git%s\n",
			WHITE,
			BLUE,
			url,
			RESET,
			RED,
			RESET,
		)
		return
	}

	//PRINT VALID .GIT
	fmt.Printf(
		"%sUrl: %s%s%s - %sValid Git%s\n",
		WHITE,
		BLUE,
		url,
		RESET,
		GREEN,
		RESET,
	)

	fileOpen, err := os.OpenFile(
		"result/valid_git.txt",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		fmt.Printf(
			"%sError: %s%s%s\n",
			WHITE,
			YELLOW,
			err.Error(),
			RESET,
		)
	}
	defer fileOpen.Close()
	_, err = fileOpen.WriteString(url + "/.git/HEAD\n")
	if err != nil {
		fmt.Printf(
			"%sError: %s%s%s\n",
			WHITE,
			YELLOW,
			err.Error(),
			RESET,
		)
	}
	response, err = s.Client.Get(url + "/.git/config")
	if err != nil {
		fmt.Printf(
			"%sError: %s%s%s\n",
			WHITE,
			YELLOW,
			err.Error(),
			RESET,
		)
	}
	defer response.Body.Close()
	responseBody, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf(
			"%sError: %s%s%s\n",
			WHITE,
			YELLOW,
			err.Error(),
			RESET,
		)
	}
	ret := ParseToken(string(responseBody))
	fmt.Printf(
		"%sUrl: %s%s%s - %s%s\n",
		WHITE,
		BLUE,
		url,
		WHITE,
		ret,
		RESET,
	)
}

func (s *Scanner) Start(thread int) {
	var wg sync.WaitGroup
	threadchan := make(chan struct{}, thread)
	for _, url := range s.Urls {
		wg.Add(1)
		threadchan <- struct{}{}
		go func(url string) {
			s.Scan(url, &wg)
			<-threadchan
		}(url)
	}
	wg.Wait()
}
