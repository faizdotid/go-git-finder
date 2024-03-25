package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

const (
	REGEX_GITHUB_TOKEN = `ghp_[a-zA-Z0-9]{30,}`
)

func CheckGithubToken(token string) {
	defer RecoverIfError()
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	request, err := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
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
	request.Header.Set(
		"Authorization",
		"Bearer "+token,
	)
	request.Header.Set(
		"Accept",
		"application/vnd.github+json",
	)
	request.Header.Set(
		"X-GitHub-Api-Version",
		"2022-11-28",
	)
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf(
			"%sToken: %s%s%s - %s%s%s\n",
			WHITE,
			BLUE,
			token,
			RESET,
			YELLOW,
			err.Error(),
			RESET,
		)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Printf(
			"%sToken: %s%s%s - %sInvalid%s\n",
			WHITE,
			BLUE,
			token,
			RESET,
			RED,
			RESET,
		)
		return
	}

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
	var user map[string]interface{}
	err = json.Unmarshal(responseBody, &user)
	if err != nil {
		fmt.Printf(
			"%sError: %s%s%s\n",
			WHITE,
			YELLOW,
			err.Error(),
			RESET,
		)
	}
	fmt.Printf(
		"%sToken: %s%s%s - %sValid%s\n",
		WHITE,
		BLUE,
		token,
		RESET,
		GREEN,
		RESET,
	)
	fileOpen, err := os.OpenFile(
		"result/valid_tokens.txt",
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
	_, err = fileOpen.WriteString(token + "\n")
	if err != nil {
		fmt.Printf(
			"%sError: %s%s%s\n",
			WHITE,
			YELLOW,
			err.Error(),
			RESET,
		)
	}
}

func ParseToken(body string) string {
	reg := regexp.MustCompile(REGEX_GITHUB_TOKEN)
	tokens := reg.FindAllString(
		body,
		-1,
	)
	if len(tokens) == 0 {
		return fmt.Sprintf(
			"%sNo tokens found",
			RED,
		)
	}
	for _, token := range tokens {
		fileOpen, err := os.OpenFile(
			"result/tokens.txt",
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
		_, err = fileOpen.WriteString(token + "\n")
		if err != nil {
			fmt.Printf(
				"%sError: %s%s%s\n",
				WHITE,
				YELLOW,
				err.Error(),
				RESET,
			)
		}
		CheckGithubToken(token)
	}
	return fmt.Sprintf(
		"%sFound %d tokens",
		GREEN,
		len(tokens),
	)
}
