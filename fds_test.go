package fds

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/gabrieloliverio/fds/config"
	"github.com/gabrieloliverio/fds/input"
	"github.com/gabrieloliverio/fds/replace"
)

func TestReplaceInFile_RenameTmpFileToOriginalFileWhenNotNil(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := path.Join(tempDir, "input")

	createTestFile(tempDir, "input", "Lorem ipsum dolor sit amet", t)

	args := input.Args{Path: input.PathArg{Value: inputPath}, Search: "Lorem", Replace: "mamãe"}
	config := config.NewConfig()
	config.Flags = map[string]bool{}
	var confirmAnswer *input.ConfirmAnswer
	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	var replacer = replace.NewFileReplacer(inputPath, args.Search, args.Replace, config)

	err := ReplaceInFile(args, replacer, stdin, confirmAnswer)

	if err != nil {
		t.Errorf("ReplaceInFile() returned an expected error '%s'\n", err)
	}

	result, err := os.ReadFile(inputPath)
	want := "mamãe ipsum dolor sit amet"

	if result := string(result); result != want {
		t.Errorf(`ReplaceInFile() = %q, want %q`, result, want)
	}
}

func TestReplaceInFile_LeavesFileUntouchedWhenNothingWasReplaced(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := path.Join(tempDir, "input")

	createTestFile(tempDir, "input", "Lorem ipsum dolor sit amet", t)
	originalStat, err := os.Stat(inputPath)

	args := input.Args{Path: input.PathArg{Value: inputPath}, Search: "no existe", Replace: "bar"}
	config := config.NewConfig()
	config.Flags = map[string]bool{}
	var confirmAnswer *input.ConfirmAnswer
	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	var replacer = replace.NewFileReplacer(inputPath, args.Search, args.Replace, config)

	err = ReplaceInFile(args, replacer, stdin, confirmAnswer)

	if err != nil {
		t.Errorf("ReplaceInFile() returned an expected error '%s'\n", err)
	}

	statAfterReplace, err := os.Stat(inputPath)

	result, err := os.ReadFile(inputPath)
	want := "Lorem ipsum dolor sit amet"

	if result := string(result); result != want {
		t.Errorf(`ReplaceInFile() = %q, want %q`, result, want)
	}

	if originalStat.ModTime().Unix() < statAfterReplace.ModTime().Unix() {
		t.Errorf(`ReplaceInFile() replaced input file when nothing was replaced in its content`)
	}
}

func TestReplaceInFiles(t *testing.T) {
	tempDir := t.TempDir()

	defaultAnswer := input.ConfirmAnswer('n')

	inputPath1 := path.Join(tempDir, "input1")
	createTestFile(tempDir, "input1", "Lorem ipsum dolor sit amet", t)

	inputPath2 := path.Join(tempDir, "input2")
	createTestFile(tempDir, "input2", "Lorem ipsum dolor sit amet", t)

	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	args := input.Args{Path: input.PathArg{Value: inputPath1}, Search: "Lorem", Replace: "mamãe"}

	config := config.NewConfig()
	config.Flags = map[string]bool{}

	err := ReplaceInFiles([]string{inputPath1, inputPath2}, stdin, args, config, &defaultAnswer)

	if err != nil {
		t.Errorf("ReplaceInFile() returned an expected error '%s'\n", err)
	}

	result1, err := os.ReadFile(inputPath1)
	result2, err := os.ReadFile(inputPath1)

	want := "mamãe ipsum dolor sit amet"

	if result := string(result1); result != want {
		t.Errorf(`ReplaceInFile() = %q, want %q`, result, want)
	}

	if result := string(result2); result != want {
		t.Errorf(`ReplaceInFile() = %q, want %q`, result, want)
	}
}

func TestGetFilesInDir_NoIgnoreGlobs_FindAllFiles(t *testing.T) {
	tempDir := t.TempDir()

	createTreeStructure(tempDir)

	result, err := GetFilesInDir(tempDir, input.IgnoreGlobs{}, false)

	if err != nil {
		t.Errorf("GetFilesInDir() returned expected error")
	}

	want := []string{
		filepath.Join(tempDir, "dir1", "file11"),
		filepath.Join(tempDir, "dir1", "file12"),

		filepath.Join(tempDir, "dir2", "file21"),
		filepath.Join(tempDir, "dir2", "file22"),

		filepath.Join(tempDir, "file1"),
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("GetFilesInDir() = %q, want %q", result, want)
	}
}

func TestGetFilesInDir_IgnoreGlobs(t *testing.T) {
	tempDir := t.TempDir()

	createTreeStructure(tempDir)

	result, err := GetFilesInDir(
		tempDir,
		input.IgnoreGlobs{filepath.Join(tempDir, "dir2/**")},
		false,
	)

	if err != nil {
		t.Errorf("GetFilesInDir() returned expected error")
	}

	want := []string{
		filepath.Join(tempDir, "dir1", "file11"),
		filepath.Join(tempDir, "dir1", "file12"),

		filepath.Join(tempDir, "file1"),
	}

	if !reflect.DeepEqual(result, want) {
		t.Errorf("GetFilesInDir() = %q, want %q", result, want)
	}
}

func TestGetFilesInDir_IgnoreGlobs_ReturnsNoFiles(t *testing.T) {
	tempDir := t.TempDir()

	createTreeStructure(tempDir)

	result, err := GetFilesInDir(
		tempDir,
		input.IgnoreGlobs{filepath.Join(tempDir, "**")},
		false,
	)

	if err != nil {
		t.Errorf("GetFilesInDir() returned expected error")
	}

	if count := len(result); count > 0 {
		t.Errorf("GetFilesInDir() expects no files, returned %d files", count)
	}
}

func createTreeStructure(tempDir string) {
	os.Create(path.Join(tempDir, "file1"))

	os.Mkdir(filepath.Join(tempDir, "dir1"), 0755)
	os.Create(path.Join(tempDir, "dir1", "file11"))
	os.Create(path.Join(tempDir, "dir1", "file12"))

	os.Mkdir(filepath.Join(tempDir, "dir2"), 0755)
	os.Create(path.Join(tempDir, "dir2", "file21"))
	os.Create(path.Join(tempDir, "dir2", "file22"))
}

func createTestFile(tempDir, fileName, inputContent string, t *testing.T) *os.File {
	var inputFile *os.File

	inputFile, err := os.Create(path.Join(tempDir, fileName))

	if err != nil {
		t.Fatalf("Failed to open input file")
	}

	_, err = inputFile.WriteString(inputContent)

	inputFile.Seek(0, io.SeekStart)

	return inputFile
}
