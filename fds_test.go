package fds

import (
	"bytes"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReplaceInFile_RenameTmpFileToOriginalFileWhenNotNil(t *testing.T) {
	tempDir := t.TempDir()
	inputPath := path.Join(tempDir, "input")

	createTestFile(tempDir, "input", "Lorem ipsum dolor sit amet", t)

	args := Args{Path: PathArg{Value: inputPath}, Search: "Lorem", Replace: "mamãe"}
	config := NewConfig()
	config.Flags = map[string]bool{}

	var confirmAnswer *ConfirmAnswer
	var stdout bytes.Buffer

	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	var replacer = NewFileReplacer(inputPath, args.Search, args.Replace, config)

	err := ReplaceInFile(replacer, stdin, &stdout, confirmAnswer)

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

	args := Args{Path: PathArg{Value: inputPath}, Search: "no existe", Replace: "bar"}
	config := NewConfig()
	config.Flags = map[string]bool{}

	var confirmAnswer *ConfirmAnswer
	var stdout bytes.Buffer

	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	var replacer = NewFileReplacer(inputPath, args.Search, args.Replace, config)

	err = ReplaceInFile(replacer, stdin, &stdout, confirmAnswer)

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

	defaultAnswer := ConfirmAnswer('n')

	inputPath1 := path.Join(tempDir, "input1")
	createTestFile(tempDir, "input1", "Lorem ipsum dolor sit amet", t)

	inputPath2 := path.Join(tempDir, "input2")
	createTestFile(tempDir, "input2", "Lorem ipsum dolor sit amet", t)

	var stdout bytes.Buffer
	stdin, _ := os.Create(path.Join(tempDir, "stdin"))

	args := Args{Path: PathArg{Value: inputPath1}, Search: "Lorem", Replace: "mamãe"}

	config := NewConfig()
	config.Flags = map[string]bool{}

	err := ReplaceInFiles([]string{inputPath1, inputPath2}, stdin, &stdout, args, config, &defaultAnswer)

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

	result, err := GetFilesInDir(tempDir, IgnoreGlobs{}, false)

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
		IgnoreGlobs{filepath.Join(tempDir, "dir2/**")},
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
		IgnoreGlobs{filepath.Join(tempDir, "**")},
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
