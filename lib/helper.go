package lib

import (
	"fmt"
	"os"
)

// PrintErr prints an error message to stderr.
func PrintErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ %sERR%s ] - %s%s%s\n", Red, Reset, Red, err, Reset)
	}
}

// PrintBanner displays the scanner banner.
func PrintBanner(urls, threads int) {
	fmt.Printf(`%s
	█▀▀ █ ▀█▀ ▄▄ █▀▀ █ █▄░█ █▀▄ █▀▀ █▀█
	█▄█ █ ░█░ ░░ █▀░ █ █░▀█ █▄▀ ██▄ █▀▄%s
	%sScanning %s%d%s urls with %s%d%s threads%s

`,
		Blue,
		Reset,
		White,
		Blue,
		urls,
		White,
		Blue,
		threads,
		White,
		Reset,
	)
}
