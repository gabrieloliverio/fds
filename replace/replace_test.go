package replace

import (
	"io"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/gabrieloliverio/fds/input"
)

func TestReplaceStringOrPattern_LiteralString(t *testing.T) {
	search := "text"
	replace := "replacement"
	subject := "this is some text, this is some other text"
	result := ReplaceStringOrPattern(search, replace, subject, true, false)

	want := regexp.MustCompile("this is some replacement, this is some other replacement")

	if !want.MatchString(result) {
		t.Errorf(`ReplaceStringOrPattern() = "%q", want "%#q"`, result, want)
	}
}

func TestReplaceStringOrPattern_RegEx(t *testing.T) {
	searchPattern := "t.xt"
	replace := "replacement"
	subject := "this is some text"
	result := ReplaceStringOrPattern(searchPattern, replace, subject, false, false)

	want := regexp.MustCompile("this is some replacement")

	if !want.MatchString(result) {
		t.Errorf(`ReplaceStringOrPattern() = "%q", want "%#q"`, result, want)
	}
}

func TestReplaceStringOrPattern_RegExIgnoringCase(t *testing.T) {
	searchPattern := "Text"
	replace := "replacement"
	subject := "this is some text"
	result := ReplaceStringOrPattern(searchPattern, replace, subject, false, true)

	want := regexp.MustCompile("this is some replacement")

	if !want.MatchString(result) {
		t.Errorf(`ReplaceStringOrPattern() = "%q", want "%#q"`, result, want)
	}
}

func TestReplaceStringOrPattern_RegExCapturingGroup(t *testing.T) {
	searchPattern := "(text)"
	replace := "other $1"
	subject := "this is some text"
	result := ReplaceStringOrPattern(searchPattern, replace, subject, false, false)

	want := regexp.MustCompile("this is some other text")

	if !want.MatchString(result) {
		t.Errorf(`ReplaceStringOrPattern() = "%q", want "%#q"`, result, want)
	}
}

func TestReplaceStringOrPattern_RegExNotMatch(t *testing.T) {
	searchPattern := "<fooo>"
	replace := "replacement"
	subject := "this is some text, this is some other text"
	result := ReplaceStringOrPattern(searchPattern, replace, subject, false, false)

	want := regexp.MustCompile("this is some text, this is some other text")

	if !want.MatchString(result) {
		t.Errorf(`ReplaceStringOrPattern() = "%q", want "%#q"`, result, want)
	}
}

func TestReplaceStringOrPattern_LiteralNotMatch(t *testing.T) {
	searchPattern := "<fooo>"
	replace := "replacement"
	subject := "this is some text, this is some other text"
	result := ReplaceStringOrPattern(searchPattern, replace, subject, true, false)

	want := regexp.MustCompile("this is some text, this is some other text")

	if !want.MatchString(result) {
		t.Errorf(`ReplaceStringOrPattern() = "%q", want "%#q"`, result, want)
	}
}

func createFiles(tempDir, inputContent string, t *testing.T) (*os.File, *os.File) {
	var inputFile, outputFile *os.File

	inputFile, err := os.Create(path.Join(tempDir, "input"))

	if err != nil {
		t.Fatalf("Failed to open input file")
	}

	_, err = inputFile.WriteString(inputContent)

	if err != nil {
		t.Fatalf("Failed to write into input file")
	}

	outputFile, _ = os.Create(path.Join(tempDir, "output"))

	if err != nil {
		t.Fatalf("Failed to open output file")
	}

	inputFile.Seek(0, io.SeekStart)

	return inputFile, outputFile
}

func TestReplaceInFile_SingleLine(t *testing.T) {
	tempDir := t.TempDir()

	inputFile, outputFile := createFiles(tempDir, "this is some text", t)

	defer inputFile.Close()
	defer outputFile.Close()

	args := input.Args{ Replace: "replacement", Search: "text" }

	err := ReplaceInFile(inputFile, outputFile, args, false, false)

	if err != nil {
		t.Fatalf("Failed to replace content on file: %s", err)
	}

	result, err := os.ReadFile(path.Join(tempDir, "output"))

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace")
	}

	want := "this is some replacement"

	if string(result) != want {
		t.Errorf(`ReplaceInFile(%+v) = "%s", want "%s"`, args, result, want)
	}
}

func TestReplaceInFile_Multiline(t *testing.T) {
	tempDir := t.TempDir()

	inputFile, outputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	defer inputFile.Close()
	defer outputFile.Close()

	args := input.Args{ Replace: "replacement", Search: "text" }

	err := ReplaceInFile(inputFile, outputFile, args, false, false)

	if err != nil {
		t.Fatalf("Failed to replace content on file: %s", err)
	}

	result, err := os.ReadFile(path.Join(tempDir, "output"))

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace")
	}

	want := "this is some replacement\nthis is some other replacement\n"

	if string(result) != want {
		t.Errorf(`ReplaceInFile(%+v) = "%s", want "%s"`, args, result, want)
	}
}

func TestReplaceInFile_NotFound(t *testing.T) {
	tempDir := t.TempDir()

	inputFile, outputFile := createFiles(tempDir, "this is some text\nthis is some other text\n", t)

	defer inputFile.Close()
	defer outputFile.Close()

	args := input.Args{ Replace: "replacement", Search: "foo" }

	err := ReplaceInFile(inputFile, outputFile, args, false, false)

	if err != nil {
		t.Fatalf("Failed to replace content on file: %s", err)
	}

	result, err := os.ReadFile(path.Join(tempDir, "output"))

	if err != nil {
		t.Fatalf("Failed to read output file after find/replace")
	}

	want := "this is some text\nthis is some other text\n"

	if string(result) != want {
		t.Errorf(`ReplaceInFile(%+v) = "%s", want "%s"`, args, result, want)
	}
}
