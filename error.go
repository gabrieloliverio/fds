package fds

import "fmt"

type InputError struct {
	message string
	Code    int
}

func (e InputError) Error() string {
	return fmt.Sprintf("%s\n\n%s", Usage, e.message)
}

type Error struct {
	message string
	Code    int
}

func (e Error) Error() string {
	return e.message
}

func NewInvalidRegExpError() InputError {
	return InputError{message: "subject is not a valid Regular Expression", Code: 42}
}

func NewInvalidArgumentsError() InputError {
	return InputError{message: "Invalid arguments", Code: 43}
}

func NewInvalidArgumentsErrorFileNotFound(filePath string) InputError {
	return InputError{message: fmt.Sprintf("File '%s' could not be found", filePath), Code: 44}
}

func NewLiteralInsensitiveError() InputError {
	return InputError{message: "[-l, --literal] cannot be used along with [ -i, --insensitive ]", Code: 45}
}

func NewConfirmNotOnFileError() InputError {
	return InputError{message: "[-c, --confirm] can only be used when files are supplied, not with STDIN nor positional arguments", Code: 45}
}

func NewFileReadError(file string) Error {
	return Error{message: fmt.Sprintf("Failed to read file %q. Do you have permission to read it?", file), Code: 46}
}

func NewFileWriteError(file string) Error {
	return Error{message: fmt.Sprintf("Failed to write file %q. Do you have permission to write in directory?", file), Code: 47}
}

func NewTempFileWriteError(dir string) Error {
	return Error{message: fmt.Sprintf("Failed to write temporary file. Do you have permission to write in directory %q?", dir), Code: 48}
}

func NewStdinReadError() Error {
	return Error{message: "Failed to read from Stdin", Code: 49}
}

func NewRenameFileError(file string) Error {
	return Error{message: fmt.Sprintf("Failed to rename temp file into original file %q", file), Code: 50}
}

func NewAbortedOperationError() Error {
	return Error{message: "Aborted operation", Code: 51}
}

func NewDirectoryReadError(dir string) Error {
	return Error{message: fmt.Sprintf("Failed to read directory %q. Do you have permission to read it?", dir), Code: 52}
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
