package main

import (
	"bufio"
	"flag"
	"fmt"
	"go-git-finder/lib"
	"os"
)

func main() {
	filename := flag.String("f", "", "File containing urls to scan")
	threads := flag.Int("t", 10, "Number of threads to use")
	flag.Parse()

	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer f.Close()

	var urls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if url := lib.ParseURL(scanner.Text()); url != "" {
			urls = append(urls, url)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	lib.PrintBanner(len(urls), *threads)

	s, err := lib.NewScanner(urls)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating scanner: %s\n", err)
		os.Exit(1)
	}

	g, err := lib.NewGithubTokenValidator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating validator: %s\n", err)
		os.Exit(1)
	}

	if err := s.Run(*threads, g); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
