package replace

import (
	"os"
	"path"
	"testing"

	"github.com/gabrieloliverio/fds/config"
	"github.com/gabrieloliverio/fds/input"
)

func TestReplaceInFile_ConfirmAll(t *testing.T) {
	tempDir := t.TempDir()

	inputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	defer inputFile.Close()
	defer stdin.Close()

	config := config.NewConfig()
	config.Flags = map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, config)

	confirm := input.ConfirmAnswer('a')
	outputFile, err := fileReplacer.Replace(stdin, &confirm)
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

	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	defer inputFile.Close()
	defer stdin.Close()

	config := config.NewConfig()
	config.Flags = map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, config)

	confirm := input.ConfirmAnswer('n')
	outputFile, err := fileReplacer.Replace(stdin, &confirm)

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

	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	defer inputFile.Close()
	defer stdin.Close()

	config := config.NewConfig()
	config.Flags = map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, config)

	confirm := input.ConfirmAnswer('q')
	outputFile, err := fileReplacer.Replace(stdin, &confirm)

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
