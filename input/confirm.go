package input

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

const (
    confirmText = "[y]es [n]o [a]ll q[uit]: "
    enter = 10

    ConfirmYes = 'y'
    ConfirmNo = 'n'
    ConfirmAll = 'a'
    ConfirmQuit = 'q'
)

type ConfirmAnswer rune

func Confirm(stdin *os.File) (rune, error) {
    var input rune
    var err error

    valid := []rune{'y', 'n', 'a', 'q'}
    reader := bufio.NewReader(stdin)

    for !slices.Contains(valid, input) {
        if input != enter {
            fmt.Print(confirmText)
        }

        input, _, err = reader.ReadRune()

        if err != nil {
            return 0, err
        }
    }

    return input, nil
}
