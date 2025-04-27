package input

import (
	"regexp"
	"strings"
	"testing"
)

func TestErrorContainsUsageText(t *testing.T) {
    err := Error{message: "Some error"}

	if !strings.HasPrefix(err.Error(), Usage) {
		t.Errorf(`Error.Error() = %q, does not start with usage text`, err.Error())
	}
}

func TestNewInvalidRegExpError(t *testing.T) {
    err := NewInvalidRegExpError()
    want := regexp.MustCompile("subject is not a valid Regular Expression")

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewInvalidRegExpError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}
}

func TestNewInvalidArgumentsError(t *testing.T) {
    err := NewInvalidArgumentsError()
    want := regexp.MustCompile("Invalid arguments")

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewInvalidArgumentsError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}
}

func TestNewInvalidArgumentsErrorFileNotFound(t *testing.T) {
    err := NewInvalidArgumentsErrorFileNotFound("foo")
    want := regexp.MustCompile("File 'foo' could not be found")

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewInvalidArgumentsErrorFileNotFound().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}
}

func TestNewLiteralInsensitiveError(t *testing.T) {
    err := NewLiteralInsensitiveError()
    want := regexp.MustCompile("cannot be used along with")

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewLiteralInsensitiveError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}
}

func TestNewConfirmNotOnFileError(t *testing.T) {
    err := NewConfirmNotOnFileError()
    want := regexp.MustCompile("can only be used when files are supplied")

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewConfirmNotOnFileError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}
}

func TestNewInvalidConfirmInputError(t *testing.T) {
    err := NewInvalidConfirmInputError('t')
    want := regexp.MustCompile("Invalid input: t")

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewInvalidConfirmInputError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}
}
