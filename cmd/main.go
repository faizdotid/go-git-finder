package main

import (
	"flag"
	"fmt"
	"go-git-finder/lib"
	"log"
	"os"
	"strings"
)

func createResultFolder() {
	if _, err := os.Stat("results"); os.IsNotExist(err) {
		os.Mkdir("results", 0755)
	}
}

func main() {
	createResultFolder() // creating results folder
	filename := flag.String("f", "", "File containing urls to scan")
	threads := flag.Int("t", 10, "Number of threads to use")
	flag.Parse()
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	buf, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatalf("Error reading file: %s", err)
	}

	// creata a slice of urls
	urls := make([]string, 0, len(strings.Split(string(buf), "\n")))
	for _, url := range strings.Split(string(buf), "\n") {
		if url == "" {
			continue
		}
		urls = append(urls, lib.ParseURL(url))
	}

	fmt.Printf(`%s
█▀▀ █ ▀█▀ ▄▄ █▀▀ █ █▄░█ █▀▄ █▀▀ █▀█
█▄█ █ ░█░ ░░ █▀░ █ █░▀█ █▄▀ ██▄ █▀▄%s
%sScanning %s%d%s urls with %s%d%s threads%s

`,
		lib.Blue,
		lib.Reset,
		lib.White,
		lib.Blue,
		len(urls),
		lib.White,
		lib.Blue,
		*threads,
		lib.White,
		lib.Reset,
	)

	// create a file to write the results
	

	// create a scanner
	scanner := lib.NewScanner(urls)
	scanner.Run(*threads)
}
