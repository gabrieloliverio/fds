package input

import "fmt"

type UsageError struct {
	arg string
	flag string
	message string
}

func (u UsageError) Error() string {
	var message string

	switch {
	case u.arg != "":
		message = fmt.Sprintf("%s\n[ %s ] %s\n", Usage, u.arg, u.message)
	case u.flag != "":
		message = fmt.Sprintf("%s\n[ %s ] %s\n", Usage, u.flag, u.message)
	default:
		message = fmt.Sprintf("%s\n%s\n", Usage, u.message)
	}

	return message
}

func NewInvalidRegExpError() UsageError {
	return UsageError{
		message: "is not a valid Regular Expression",
		arg: "subject",
	}
}

func NewInvalidArgumentsError() UsageError {
	return UsageError{message: "Invalid arguments"}
}

func NewLiteralInsensitiveError() UsageError {
	return UsageError{
		message: "cannot be used along with [ -i, --insensitive ]",
		flag: "-l, --literal",
	}
}
