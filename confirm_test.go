package fds

import (
	"io"
	"os"
	"path"
	"testing"
)

func createFile(tempDir, inputContent string, t *testing.T) *os.File {
	var file *os.File

	file, err := os.Create(path.Join(tempDir, "input"))

	if err != nil {
		t.Fatalf("Failed to open input file")
	}

	_, err = file.WriteString(inputContent)

	if err != nil {
		t.Fatalf("Failed to write input file")
	}

	file.Seek(0, io.SeekStart)

	return file
}

func TestConfirm_Valid(t *testing.T) {
	stdin := createFile(t.TempDir(), "y", t)

	want := 'y'
	result, err := Confirm(stdin, "text", []rune{'y', 'n'})

	if err != nil {
		t.Fatalf("Confirm() does not expect error, got error %s", err)
	}

	if want != result {
		t.Errorf("Confirm() = %c, want %c", result, want)
	}
}

func TestConfirm_Invalid(t *testing.T) {
	stdin := createFile(t.TempDir(), "*", t)

	_, err := Confirm(stdin, "text", []rune{'y', 'n'})

	if err == nil {
		t.Fatalf("Confirm() expects error, did not get error %s", err)
	}
}
