package replace

import (
	"os"
	"path"
	"testing"

	"github.com/gabrieloliverio/fds/input"
)

func TestReplaceInFile_ConfirmAll(t *testing.T) {
	tempDir := t.TempDir()

	inputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	defer inputFile.Close()
	defer stdin.Close()

	flags := map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, flags)

	confirm := input.ConfirmAnswer('a')
	outputFile, fileChanged, err := fileReplacer.Replace(stdin, &confirm)
	defer os.Remove(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace: %s", err)
	}

	wantText := "this is some replacement\nthis is some other replacement\n"
	wantChanged := true

	if string(result) != wantText || fileChanged != wantChanged {
		t.Errorf(`ReplaceInFile(%s, %s) = %q, %t, want %q, %t`, search, replace, result, fileChanged, wantText, wantChanged)
	}
}

func TestReplaceInFile_ConfirmNo(t *testing.T) {
	tempDir := t.TempDir()

	inputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	defer inputFile.Close()
	defer stdin.Close()

	flags := map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, flags)

	confirm := input.ConfirmAnswer('n')
	outputFile, fileChanged, err := fileReplacer.Replace(stdin, &confirm)
	defer os.Remove(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace: %s", err)
	}

	wantText := "this is some text\nthis is some other text\n"
	wantChanged := false

	if string(result) != wantText || fileChanged != wantChanged {
		t.Errorf(`ReplaceInFile(%s, %s) = %q, %t, want %q, %t`, search, replace, result, fileChanged, wantText, wantChanged)
	}
}

func TestReplaceInFile_ConfirmQuit(t *testing.T) {
	tempDir := t.TempDir()

	inputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	defer inputFile.Close()
	defer stdin.Close()

	flags := map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, flags)

	confirm := input.ConfirmAnswer('q')
	outputFile, fileChanged, err := fileReplacer.Replace(stdin, &confirm)
	defer os.Remove(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace: %s", err)
	}

	wantText := "this is some text\nthis is some other text\n"
	wantChanged := false

	if string(result) != wantText || fileChanged != wantChanged {
		t.Errorf(`ReplaceInFile(%s, %s) = %q, %t, want %q, %t`, search, replace, result, fileChanged, wantText, wantChanged)
	}
}
