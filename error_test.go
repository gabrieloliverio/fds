package fds

import (
	"regexp"
	"strings"
	"testing"
)

func TestInputErrorContainsUsageText(t *testing.T) {
	err := InputError{message: "Some error"}

	if !strings.HasPrefix(err.Error(), Usage) {
		t.Errorf(`Error.Error() = %q, does not start with usage text`, err.Error())
	}
}

func TestErrorDoesNotContainUsageText(t *testing.T) {
	err := Error{message: "Some error"}

	if strings.HasPrefix(err.Error(), Usage) {
		t.Errorf(`Error.Error() = %q, starts with usage text and it should not`, err.Error())
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

func TestNewFileReadError(t *testing.T) {
	err := NewFileReadError("/file/path")
	want := regexp.MustCompile(`Failed to read file "/file/path"`)
	code := 46

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewFileReadError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}

	if err.Code != code {
		t.Errorf(`NewFileReadError().Code = %d, want %d`, err.Code, code)
	}
}

func TestNewFileWriteError(t *testing.T) {
	err := NewFileWriteError("/file/path")
	want := regexp.MustCompile(`Failed to write file "/file/path"`)
	code := 47

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewFileWriteError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}

	if err.Code != code {
		t.Errorf(`NewFileWriteError().Code = %d, want %d`, err.Code, code)
	}
}

func TestNewTempFileWriteError(t *testing.T) {
	err := NewTempFileWriteError("/file/path")
	want := regexp.MustCompile(`Failed to write temporary file`)
	code := 48

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewTempFileWriteError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}

	if err.Code != code {
		t.Errorf(`NewTempFileWriteError().Code = %d, want %d`, err.Code, code)
	}
}

func TestNewStdinReadError(t *testing.T) {
	err := NewStdinReadError()
	want := regexp.MustCompile(`Failed to read from Stdin`)
	code := 49

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewStdinReadError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}

	if err.Code != code {
		t.Errorf(`NewStdinReadError().Code = %d, want %d`, err.Code, code)
	}
}

func TestNewRenameFileError(t *testing.T) {
	err := NewRenameFileError("/file/path")
	want := regexp.MustCompile(`Failed to rename temp file into original file \"/file/path\"`)
	code := 50

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewRenameFileError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}

	if err.Code != code {
		t.Errorf(`NewRenameFileError().Code = %d, want %d`, err.Code, code)
	}
}

func TestNewAbortedOperationError(t *testing.T) {
	err := NewAbortedOperationError()
	want := regexp.MustCompile(`Aborted operation`)
	code := 51

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewAbortedOperationError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}

	if err.Code != code {
		t.Errorf(`NewAbortedOperationError().Code = %d, want %d`, err.Code, code)
	}
}

func TestNewDirectoryReadError(t *testing.T) {
	err := NewDirectoryReadError("/file/path")
	want := regexp.MustCompile(`Failed to read directory "/file/path"`)
	code := 52

	if !want.MatchString(err.Error()) {
		t.Errorf(`NewDirectoryReadError().Error() = %q, does not match RegExp %q`, err.Error(), want)
	}

	if err.Code != code {
		t.Errorf(`NewDirectoryReadError().Code = %d, want %d`, err.Code, code)
	}
}
