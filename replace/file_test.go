package replace

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/gabrieloliverio/fds/config"
)

func TestFileReplacer_ReplaceInFile_SingleLine(t *testing.T) {
	tempDir := t.TempDir()

	inputFile := createFiles(tempDir, "this is some text", t)
	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	defer inputFile.Close()
	defer stdin.Close()

	config := config.NewConfig()
	config.Flags = map[string]bool{"insensitive": false, "confirm": false, "literal": false}
	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, config)

	outputFile, err := fileReplacer.Replace(stdin, nil)
	defer os.Remove(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(outputFile.Name())
	fmt.Println(result)

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace: %s", err)
	}

	want:= "this is some replacement"

	if string(result) != want{
		t.Errorf(`ReplaceInFile(%s, %s) = %q, want %qt`, search, replace, result, want)
	}
}

func TestReplaceInFile_Multiline(t *testing.T) {
	tempDir := t.TempDir()

	inputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)
	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	defer inputFile.Close()
	defer stdin.Close()

	config := config.NewConfig()
	config.Flags = map[string]bool{"insensitive": false, "confirm": false, "literal": false}

	search := "text"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, config)

	outputFile, err := fileReplacer.Replace(stdin, nil)
	defer os.Remove(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	result, err := os.ReadFile(outputFile.Name())

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace: %s", err)
	}

	want:= "this is some replacement\nthis is some other replacement\n"

	if string(result) != want{
		t.Errorf(`ReplaceInFile(%s, %s) = %q, want %q`, search, replace, result, want)
	}
}

func TestReplaceInFile_NotFound(t *testing.T) {
	var result []byte
	tempDir := t.TempDir()

	inputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)
	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	defer inputFile.Close()
	defer stdin.Close()

	config := config.NewConfig()
	config.Flags = map[string]bool{"insensitive": false, "confirm": false, "literal": false}

	search := "foo"
	replace := "replacement"

	fileReplacer := NewFileReplacer(inputFile.Name(), search, replace, config)

	outputFile, err := fileReplacer.Replace(stdin, nil)
	if outputFile != nil {
		result, err = os.ReadFile(outputFile.Name())
	}

	if err != nil {
		t.Fatalf("Failed to replace content on file: %q", err)
	}

	if outputFile != nil {
		t.Errorf(`ReplaceInFile(%s, %s) should have returned nil as output file. File with content %q returned`, search, replace, result)
	}
}

func TestOpenInputFile(t *testing.T) {
	tempDir := t.TempDir()

	os.Create(path.Join(tempDir, "file"))
	os.Symlink(path.Join(tempDir, "file"), path.Join(tempDir, "symlink"))

	resolvedFile, _ := openInputFile(path.Join(tempDir, "symlink"))
	stat, _ := resolvedFile.Stat()

	if stat.Mode() == os.ModeSymlink {
		t.Errorf("OpenInputFile() resolved a symlink instead of file")
	}
}
