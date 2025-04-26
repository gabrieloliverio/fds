package input

import "fmt"

type Error struct {
	message string
}

func (u Error) Error() string {
	return fmt.Sprintf("%s\n\n%s", Usage, u.message)
}

func NewInvalidRegExpError() Error {
	return Error{message: "subject is not a valid Regular Expression"}
}

func NewInvalidArgumentsError() Error {
	return Error{message: "Invalid arguments"}
}

func NewInvalidArgumentsErrorFileNotFound(filePath string) Error {
	return Error{message: fmt.Sprintf("File '%s' could not be found", filePath)}
}

func NewLiteralInsensitiveError() Error {
	return Error{message: "[-l, --literal] cannot be used along with [ -i, --insensitive ]"}
}

func NewConfirmNotOnFileError() Error {
	return Error{message: "[-c, --confirm] can only be used when files are supplied, not with STDIN nor positional arguments"}
}

type ConfirmError struct {
	input rune
	message string
}

func (e ConfirmError) Error() string {
	return fmt.Sprintf("%s: %c", e.message, e.input)
}

func NewInvalidConfirmInputError(input rune) ConfirmError {
	return ConfirmError{
		message: "Invalid input",
		input: input,
	}
}
