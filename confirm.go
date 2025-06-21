package fds

import (
	"bufio"
	"fmt"
	"io"
	"slices"
)

const (
	enter = 10
)

type ConfirmAnswer rune

func Confirm(stdin io.Reader, text string, valid []rune) (rune, error) {
	var input rune
	var err error

	reader := bufio.NewReader(stdin)

	for !slices.Contains(valid, input) {
		if input != enter {
			fmt.Print(text + ": ")
		}

		input, _, err = reader.ReadRune()

		if err != nil {
			return 0, err
		}
	}

	return input, nil
}
