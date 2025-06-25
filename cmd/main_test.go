package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/gabrieloliverio/fds"
)

func TestExecuteHelp(t *testing.T) {
	tempDir := t.TempDir()

	config := fds.NewConfig()
	config.Flags = map[string]bool{"help": true}

	var stdout bytes.Buffer
	var stdin, _ = os.Create(filepath.Join(tempDir, "stdin"))
	stdin.WriteString("lorem ipsum")

	err := execute([]string{"cmd"}, config, stdin, &stdout)

	if err != nil {
		t.Errorf("execute() was not supposed to return error, but err %s was returned", err)
	}

	want := fds.Usage
	output, err := io.ReadAll(&stdout)

	if want != string(output) {
		t.Errorf("execute() printed %q, wanted %q", output, want)
	}
}

func TestExecuteInvalidArgumentError(t *testing.T) {
	tempDir := t.TempDir()

	config := fds.NewConfig()
	config.Flags = map[string]bool{}

	var stdout bytes.Buffer
	var stdin, _ = os.Create(filepath.Join(tempDir, "stdin"))
	stdin.WriteString("lorem ipsum")

	err := execute([]string{"cmd"}, config, stdin, &stdout)

	if err == nil {
		t.Error("execute() was supposed to return error, but none was returned")
	}

	if caughtErr, ok := err.(fds.Error); ok {
		if caughtErr.Code != 43 {
			t.Errorf("execute() was supposed to return an Invalid Arguments error (code 43). Error %d was returned", caughtErr.Code)
		}
	}
}

func TestExecuteLiteralAndInsensitiveError(t *testing.T) {
	tempDir := t.TempDir()

	config := fds.NewConfig()
	config.Flags = map[string]bool{"literal": true, "insensitive": true}

	args := []string{"foo", "bar"}

	var stdout bytes.Buffer
	var stdin, _ = os.Create(filepath.Join(tempDir, "stdin"))
	stdin.WriteString("lorem ipsum")

	err := execute(args, config, stdin, &stdout)

	if err == nil {
		t.Error("execute() was supposed to return error, but none was returned")
	}

	if caughtErr, ok := err.(fds.Error); ok {
		if caughtErr.Code != 45 {
			t.Errorf("execute() was supposed to return an Literal Insensitive error (code 45). Error %d was returned", caughtErr.Code)
		}
	}
}

func TestExecuteFileNotFoundError(t *testing.T) {
	tempDir := t.TempDir()

	config := fds.NewConfig()
	config.Flags = map[string]bool{}

	args := []string{"foo", "bar", "baz"}

	var stdin, _ = os.Create(filepath.Join(tempDir, "stdin"))
	var stdout bytes.Buffer

	err := execute(args, config, stdin, &stdout)

	if err == nil {
		t.Error("execute() was supposed to return error, but none was returned")
	}

	if caughtErr, ok := err.(fds.Error); ok {
		if caughtErr.Code != 44 {
			t.Errorf("execute() was supposed to return an File not Found error (code 44). Error %d was returned", caughtErr.Code)
		}
	}
}

func TestExecuteInvalidRegexError(t *testing.T) {
	tempDir := t.TempDir()

	config := fds.NewConfig()
	config.Flags = map[string]bool{}

	args := []string{"(lorem", "bar"}

	var stdout bytes.Buffer
	var stdin, _ = os.Create(filepath.Join(tempDir, "stdin"))
	stdin.WriteString("lorem ipsum")
	stdin.Seek(0, io.SeekStart)

	err := execute(args, config, stdin, &stdout)

	if err == nil {
		t.Error("execute() was supposed to return error, but none was returned")
	}

	if caughtErr, ok := err.(fds.Error); ok {
		if caughtErr.Code != 42 {
			t.Errorf("execute() was supposed to return an Invalid RegExp error (code 42). Error %d was returned", caughtErr.Code)
		}
	}
}

func TestExecuteWithStdinSuccess(t *testing.T) {
	tempDir := t.TempDir()

	config := fds.NewConfig()
	config.Flags = map[string]bool{}

	args := []string{"lorem", "bar"}

	var stdout bytes.Buffer
	var stdin, _ = os.Create(filepath.Join(tempDir, "stdin"))
	stdin.WriteString("lorem ipsum")
	stdin.Seek(0, io.SeekStart)

	err := execute(args, config, stdin, &stdout)

	if err != nil {
		t.Errorf("execute() was not supposed to return error, but %q was returned", err)
	}

	want := "bar ipsum"
	result, _ := io.ReadAll(&stdout)

	if string(result) != want {
		t.Errorf("execute() was supposed to print %q. %q printed instead", want, string(result))
	}
}

func TestExecuteWithFileSuccess(t *testing.T) {
	tempDir := t.TempDir()

	config := fds.NewConfig()
	config.Flags = map[string]bool{}

	path := filepath.Join(tempDir, "input")
	file, err := os.Create(path)

	file.WriteString("lorem ipsum")
	file.Seek(0, io.SeekStart)

	args := []string{"lorem", "bar", path}

	var stdin, _ = os.Create(filepath.Join(tempDir, "stdin"))
	var stdout bytes.Buffer

	err = execute(args, config, stdin, &stdout)

	if err != nil {
		t.Errorf("execute() was not supposed to return error, but %q was returned", err)
	}

	want := "bar ipsum"
	result, _ := os.ReadFile(path)

	if string(result) != want {
		t.Errorf("execute() result is %q, want %q", string(result), want)
	}
}

func TestExecuteWithFileAndInsensitiveFlagSuccess(t *testing.T) {
	tempDir := t.TempDir()

	path := filepath.Join(tempDir, "input")
	file, err := os.Create(path)

	file.WriteString("Lorem ipsum")
	file.Seek(0, io.SeekStart)

	args := []string{"lorem", "bar", path}

	config := fds.NewConfig()
	config.Flags = map[string]bool{"insensitive": true}

	var stdin, _ = os.Create(filepath.Join(tempDir, "stdin"))
	var stdout bytes.Buffer

	err = execute(args, config, stdin, &stdout)

	if err != nil {
		t.Errorf("execute() was not supposed to return error, but %q was returned", err)
	}

	want := "bar ipsum"
	result, _ := os.ReadFile(path)

	if string(result) != want {
		t.Errorf("execute() result is %q, want %q", string(result), want)
	}
}

func TestExecuteWithDirectory(t *testing.T) {
	tempDir := t.TempDir()

	path1 := filepath.Join(tempDir, "input1")
	path2 := filepath.Join(tempDir, "input2")

	file1, err := os.Create(path1)
	file2, err := os.Create(path2)

	file1.WriteString("lorem ipsum")
	file1.Seek(0, io.SeekStart)
	file2.WriteString("dolor sit amet")
	file2.Seek(0, io.SeekStart)

	args := []string{"lorem", "bar", tempDir}

	config := fds.NewConfig()
	config.Flags = map[string]bool{}

	var stdin, _ = os.Create(filepath.Join(tempDir, "stdin"))
	var stdout bytes.Buffer

	err = execute(args, config, stdin, &stdout)

	if err != nil {
		t.Errorf("execute() was not supposed to return error, but %q was returned", err)
	}

	want1 := "bar ipsum"
	resultFile1, _ := os.ReadFile(path1)

	want2 := "dolor sit amet"
	resultFile2, _ := os.ReadFile(path2)

	if string(resultFile1) != want1 {
		t.Errorf("execute() result is %q, want %q", string(resultFile1), want1)
	}

	if string(resultFile2) != want2 {
		t.Errorf("execute() result is %q, want %q", string(resultFile2), want2)
	}
}
