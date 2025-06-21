package fds

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestReplaceInFile_ConfirmAll(t *testing.T) {
	tempDir := t.TempDir()

	inputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	var stdin io.Reader
	var stdout bytes.Buffer

	defer inputFile.Close()

	config := NewConfig()
	config.Flags = map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, config)

	confirm := ConfirmAnswer('a')
	outputFile, err := fileReplacer.Replace(stdin, &stdout, &confirm)
	defer os.Remove(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace: %s", err)
	}

	wantText := "this is some replacement\nthis is some other replacement\n"

	if string(result) != wantText {
		t.Errorf(`ReplaceInFile(%s, %s) = %q, want %q`, search, replace, result, wantText)
	}
}

func TestReplaceInFile_ConfirmNo(t *testing.T) {
	var result []byte
	tempDir := t.TempDir()

	inputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	var stdin = bytes.NewBuffer([]byte{'n'})
	var stdout bytes.Buffer

	defer inputFile.Close()

	config := NewConfig()
	config.Flags = map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, config)

	confirm := ConfirmAnswer('n')
	outputFile, err := fileReplacer.Replace(stdin, &stdout, &confirm)

	if outputFile != nil {
		result, err = os.ReadFile(outputFile.Name())
	}

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	if outputFile != nil {
		t.Errorf(`ReplaceInFile(%s, %s) should have returned nil as output file. File with content returned %s`, search, replace, result)
	}
}

func TestReplaceInFile_ConfirmQuit(t *testing.T) {
	var result []byte
	tempDir := t.TempDir()

	inputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	var stdin = bytes.NewBuffer([]byte{'q'})
	var stdout bytes.Buffer

	defer inputFile.Close()

	config := NewConfig()
	config.Flags = map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, config)

	confirm := ConfirmAnswer('q')
	outputFile, err := fileReplacer.Replace(stdin, &stdout, &confirm)

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	if outputFile != nil {
		result, err = os.ReadFile(outputFile.Name())
	}

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	if outputFile != nil {
		t.Errorf(`ReplaceInFile(%s, %s) should have returned nil as output file. File with content returned %s`, search, replace, result)
	}
}
