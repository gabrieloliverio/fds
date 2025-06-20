package fds

import "fmt"

type Error struct {
	message string
	Code    int
}

func (u Error) Error() string {
	return fmt.Sprintf("%s\n\n%s", Usage, u.message)
}

func NewInvalidRegExpError() Error {
	return Error{message: "subject is not a valid Regular Expression", Code: 42}
}

func NewInvalidArgumentsError() Error {
	return Error{message: "Invalid arguments", Code: 43}
}

func NewInvalidArgumentsErrorFileNotFound(filePath string) Error {
	return Error{message: fmt.Sprintf("File '%s' could not be found", filePath), Code: 44}
}

func NewLiteralInsensitiveError() Error {
	return Error{message: "[-l, --literal] cannot be used along with [ -i, --insensitive ]", Code: 45}
}

func NewConfirmNotOnFileError() Error {
	return Error{message: "[-c, --confirm] can only be used when files are supplied, not with STDIN nor positional arguments", Code: 45}
}

type ConfirmError struct {
	input   rune
	message string
	Code    int
}

func (e ConfirmError) Error() string {
	return fmt.Sprintf("%s: %c", e.message, e.input)
}

func NewInvalidConfirmInputError(input rune) ConfirmError {
	return ConfirmError{
		message: "Invalid input",
		input:   input,
		Code:    46,
	}
}
