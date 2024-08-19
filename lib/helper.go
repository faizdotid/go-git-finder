package lib

import (
	"fmt"
)

func RecoverIfError() {
	defer func() {
		if r := recover(); r != nil {
			PrintErr(fmt.Errorf("%v", r))
		}
	}()
}

func PrintErr(err error) {
	if err != nil {
		fmt.Printf(
			"[ %sERR%s ] - %s\n",
			Red, Reset, err.Error(),
		)
	}
}
