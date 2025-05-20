package replace

import (
	"os"
	"path"
	"testing"

	"github.com/gabrieloliverio/fds/input"
)

func TestFileReplacer_ReplaceInFile_SingleLine(t *testing.T) {
	tempDir := t.TempDir()

	inputFile, outputFile := createFiles(tempDir, "this is some text", t)

	defer inputFile.Close()
	defer outputFile.Close()

	flags := map[string]bool{"insensitive": false, "confirm": false, "literal": false}
	search := "text"
	replace :=	"replacement"

	fileReplacer := NewFileReplacer(inputFile, outputFile, flags)
	pattern := fileReplacer.CompilePattern(search)

	err := fileReplacer.ReplaceInFile(pattern, replace, nil, nil)

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(path.Join(tempDir, "output"))

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace")
	}

	want := "this is some replacement"

	if string(result) != want {
		t.Errorf(`ReplaceInFile(%s, %s) = "%q", want "%q"`, search, replace, result, want)
	}
}

func TestReplaceInFile_Multiline(t *testing.T) {
	tempDir := t.TempDir()

	inputFile, outputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	defer inputFile.Close()
	defer outputFile.Close()

	flags := map[string]bool{"insensitive": false, "confirm": false, "literal": false}

	search := "text"
	replace :=	"replacement"

	fileReplacer := NewFileReplacer(inputFile, outputFile, flags)
	pattern := fileReplacer.CompilePattern(search)

	err := fileReplacer.ReplaceInFile(pattern, replace, nil, nil)

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(path.Join(tempDir, "output"))

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace")
	}

	want := "this is some replacement\nthis is some other replacement\n"

	if string(result) != want {
		t.Errorf(`ReplaceInFile(%s, %s) = "%q", want "%q"`, search, replace, result, want)
	}
}

func TestReplaceInFile_NotFound(t *testing.T) {
	tempDir := t.TempDir()

	inputFile, outputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	defer inputFile.Close()
	defer outputFile.Close()

	flags := map[string]bool{"insensitive": false, "confirm": false, "literal": false}

	search := "foo"
	replace :=	"replacement"

	fileReplacer := NewFileReplacer(inputFile, outputFile, flags)
	pattern := fileReplacer.CompilePattern(search)

	err := fileReplacer.ReplaceInFile(pattern, replace, nil, nil)

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(path.Join(tempDir, "output"))

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace")
	}

	want := "this is some text\nthis is some other text\n"

	if string(result) != want {
		t.Errorf(`ReplaceInFile(%s, %s) = "%q", want "%q"`, search, replace, result, want)
	}
}

func TestReplaceInFile_ConfirmAll(t *testing.T) {
	tempDir := t.TempDir()

	inputFile, outputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	defer inputFile.Close()
	defer outputFile.Close()

	flags := map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace :=	"replacement"

	fileReplacer := NewFileReplacer(inputFile, outputFile, flags)
	pattern := fileReplacer.CompilePattern(search)

	confirm := input.ConfirmAnswer('a')
	err := fileReplacer.ReplaceInFile(pattern, replace, nil, &confirm)

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(path.Join(tempDir, "output"))

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace")
	}

	want := "this is some replacement\nthis is some other replacement\n"

	if string(result) != want {
		t.Errorf(`ReplaceInFile(%s, %s) = "%q", want "%q"`, search, replace, result, want)
	}
}

func TestReplaceInFile_ConfirmNo(t *testing.T) {
	tempDir := t.TempDir()

	inputFile, outputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	defer inputFile.Close()
	defer outputFile.Close()

	flags := map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace :=	"replacement"

	fileReplacer := NewFileReplacer(inputFile, outputFile, flags)
	pattern := fileReplacer.CompilePattern(search)

	confirm := input.ConfirmAnswer('n')
	err := fileReplacer.ReplaceInFile(pattern, replace, nil, &confirm)

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(path.Join(tempDir, "output"))

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace")
	}

	want := "this is some text\nthis is some other text\n"

	if string(result) != want {
		t.Errorf(`ReplaceInFile(%s, %s) = "%q", want "%q"`, search, replace, result, want)
	}
}

func TestReplaceInFile_ConfirmQuit(t *testing.T) {
	tempDir := t.TempDir()

	inputFile, outputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	defer inputFile.Close()
	defer outputFile.Close()

	flags := map[string]bool{"insensitive": false, "confirm": true, "literal": false}

	search := "text"
	replace :=	"replacement"

	fileReplacer := NewFileReplacer(inputFile, outputFile, flags)
	pattern := fileReplacer.CompilePattern(search)

	confirm := input.ConfirmAnswer('q')
	err := fileReplacer.ReplaceInFile(pattern, replace, nil, &confirm)

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(path.Join(tempDir, "output"))

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace")
	}

	want := "this is some text\nthis is some other text\n"

	if string(result) != want {
		t.Errorf(`ReplaceInFile(%s, %s) = "%q", want "%q"`, search, replace, result, want)
	}
}

func TestResolveInputFileSymlink(t *testing.T) {
	tempDir := t.TempDir()

	os.Create(path.Join(tempDir, "file"))
	os.Symlink(path.Join(tempDir, "file"), path.Join(tempDir, "symlink"))

	resolvedFile, _ := resolveInputFile(path.Join(tempDir, "symlink"))
	stat, _ := resolvedFile.Stat()

	if stat.Mode() == os.ModeSymlink {
		t.Errorf("resolveInputFile() resolved a symlink instead of file")
	}
}
