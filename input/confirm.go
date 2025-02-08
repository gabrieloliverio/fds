package input

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

const (
    confirmText = "[y]es [n]o [a]ll q[uit]: "

    ConfirmYes = 'y'
    ConfirmNo = 'n'
    ConfirmAll = 'a'
    ConfirmQuit = 'q'
)

func Confirm(stdin *os.File) (rune, error) {
    valid := []rune{'y', 'n', 'a', 'q'}

    reader := bufio.NewReader(stdin)

    // TODO: how to read a rune without pressing <ENTER>?
    fmt.Print(confirmText)
    input, _, err := reader.ReadRune()

    if err != nil {
        return 0, err
    }

    if slices.Contains(valid, input) {
        return input, nil
    }

    return 0, NewInvalidConfirmInputError(input)
}
