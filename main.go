package main

import (
	"flag"
	"fmt"
	"go-git-finder/lib"
	"os"
	"strings"
	"time"
)

func CreateResultFolder() {
	if _, err := os.Stat("result"); os.IsNotExist(err) {
		os.Mkdir("result", 0755)
	}
}

func main() {
	CreateResultFolder()
	var (
		filename string
		thread   int
	)
	flag.StringVar(&filename, "f", "", "File containing list of urls")
	flag.IntVar(&thread, "t", 0, "Number of threads")
	flag.Parse()
	if filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if thread == 0 {
		thread = 10
		fmt.Printf(
			"%sUsing default thread: %s%d%s\n",
			lib.WHITE,
			lib.BLUE,
			thread,
			lib.RESET,
		)
	}
	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	urls := strings.Split(string(file), "\n")

	fmt.Printf(`%s
█▀▀ █ ▀█▀ ▄▄ █▀▀ █ █▄░█ █▀▄ █▀▀ █▀█
█▄█ █ ░█░ ░░ █▀░ █ █░▀█ █▄▀ ██▄ █▀▄%s
%sScanning %s%d%s urls with %s%d%s threads%s

`,
		lib.BLUE,
		lib.RESET,
		lib.WHITE,
		lib.BLUE,
		len(urls),
		lib.WHITE,
		lib.BLUE,
		thread,
		lib.WHITE,
		lib.RESET,
	)
	time.Sleep(
		2 * time.Second,
	)
	scanner := lib.NewScanner(urls)

	scanner.Start(thread)

}
