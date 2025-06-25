package fds

import (
	"io"
	"os"
	"path"
	"testing"
)

func createTempFile(tempDir, inputContent string, t *testing.T) *os.File {
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

func TestValidate(t *testing.T) {
	type validationInput struct {
		args  Args
		usage string
		flags map[string]bool
	}

	type test []struct {
		name        string
		input       validationInput
		expectError bool
	}

	tests := test{
		// Success scenarios
		{
			name: "Valid subject, search, replace. Literal and insensitive flag false",
			input: validationInput{
				args:  Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
				usage: "",

				flags: map[string]bool{"literal": false, "insensitive": false, "confirm": false},
			},
			expectError: false,
		},
		{
			name: "Valid subject, search, replace. Literal flag true",
			input: validationInput{
				args:  Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
				usage: "",
				flags: map[string]bool{"literal": true, "insensitive": false, "confirm": false},
			},
			expectError: false,
		},
		{
			name: "Valid subject, search, replace. Insensitive flag true",
			input: validationInput{
				args:  Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
				usage: "",
				flags: map[string]bool{"literal": false, "insensitive": true, "confirm": false},
			},
			expectError: false,
		},
		{
			name: "Valid subject (file content), search, replace. Confirm flag true",
			input: validationInput{
				args:  Args{Path: PathArg{Value: "./foo"}, Subject: "Foo", Search: "Foo", Replace: "Baz"},
				usage: "",
				flags: map[string]bool{"literal": false, "insensitive": false, "confirm": true},
			},
			expectError: false,
		},
		// Error scenarios
		{
			name: "Insensitive and literal flag true",
			input: validationInput{
				args:  Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
				usage: "",
				flags: map[string]bool{"literal": true, "insensitive": true, "confirm": false},
			},
			expectError: true,
		},
		{
			name: "No Subject",
			input: validationInput{
				args:  Args{Subject: "", Search: "Foo", Replace: "Baz"},
				usage: "",
				flags: map[string]bool{"literal": false, "insensitive": false, "confirm": false},
			},
			expectError: true,
		},
		{
			name: "No search",
			input: validationInput{
				args:  Args{Subject: "Foo Bar", Search: "", Replace: "Baz"},
				usage: "",
				flags: map[string]bool{"literal": false, "insensitive": false, "confirm": false},
			},
			expectError: true,
		},
		{
			name: "No replace",
			input: validationInput{
				args:  Args{Subject: "Foo Bar", Search: "Foo", Replace: ""},
				usage: "",
				flags: map[string]bool{"literal": false, "insensitive": false, "confirm": false},
			},
			expectError: true,
		},
		{
			name: "Confirm flag without file",
			input: validationInput{
				args:  Args{Subject: "Foo Bar", Search: "Foo", Replace: ""},
				usage: "",
				flags: map[string]bool{"literal": false, "insensitive": false, "confirm": true},
			},
			expectError: true,
		},
		{
			name: "Invalid regexp",
			input: validationInput{
				args:  Args{Subject: "Foo Bar", Search: "((no bueno)", Replace: ""},
				usage: "",
				flags: map[string]bool{"literal": false, "insensitive": false, "confirm": false},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(
				tc.input.args,
				tc.input.flags,
			)

			if err != nil && !tc.expectError {
				t.Errorf("TestValidate(%+v) did not expected error, got error %s", tc.input, err)
			}

			if err == nil && tc.expectError {
				t.Errorf("TestValidate(%+v) expected error, did not get error", tc.input)
			}
		})
	}
}

func TestReadArgs_Stdin(t *testing.T) {
	stdin := createTempFile(os.TempDir(), "my subject", t)

	want := Args{Subject: "my subject", Search: "search", Replace: "replace"}
	result, _ := ReadArgs(stdin, []string{"search", "replace"})

	if result != want {
		t.Errorf(`ReadArgs() = "%+v", want "%+v"`, result, want)
	}
}

func TestReadArgs_File(t *testing.T) {
	tempDir := os.TempDir()

	file, _ := os.Create(path.Join(tempDir, "file"))
	file.WriteString("Lorem ipsum")
	file.Seek(0, io.SeekStart)

	stdin := createTempFile(tempDir, "", t)
	result, _ := ReadArgs(stdin, []string{"search", "replace", file.Name()})

	if result.Path.Value != file.Name() || result.Subject != file.Name() {
		t.Errorf(`ReadArgs() did not return file path in args. Subject: %v. Path.Value: %v`, result.Subject, result.Path.Value)
	}
}

func TestReadArgs_Stdin_NoParametersReturnError(t *testing.T) {
	stdin := createTempFile(os.TempDir(), "my subject", t)

	_, err := ReadArgs(stdin, []string{})

	if err == nil {
		t.Errorf(`ReadArgs() expected error, did not get error"`)
	}
}

func TestReadArgs_FileNotFound(t *testing.T) {
	stdin := createTempFile(os.TempDir(), "", t)

	_, err := ReadArgs(stdin, []string{"search", "replace", "./file_not_found"})

	if err == nil {
		t.Errorf(`ReadArgs() expected error, did not get error"`)
	}
}

func TestReadArgs_NoParametersReturnError(t *testing.T) {
	stdin := createTempFile(os.TempDir(), "", t)

	_, err := ReadArgs(stdin, []string{})

	if err == nil {
		t.Errorf(`ReadArgs() expected error, did not get error"`)
	}
}

func TestIgnoreGlobs_MatchAny(t *testing.T) {
	type test struct {
		name  string
		input string
		globs []string
		want  bool
	}

	tests := []test{
		{
			name:  "path should match",
			input: "/foo",
			globs: []string{"/foo", "/bar"},
			want:  true,
		},
		{
			name:  "path not should match",
			input: "/foo",
			globs: []string{"/baz", "/bar"},
			want:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ignoreGlob := IgnoreGlobs(tc.globs)

			result := ignoreGlob.MatchAny(tc.input)

			if result != tc.want {
				t.Errorf(`IgnoreGlob.MatchAny() = "%t", want "%t"`, result, tc.want)
			}
		})
	}
}

func TestIgnoreGlobs_String(t *testing.T) {
	type test struct {
		globs []string
		want  string
	}

	tests := []test{
		{
			globs: []string{"/foo", "/bar"},
			want:  "/foo,/bar",
		},
		{
			globs: []string{"/baz", "/bar"},
			want:  "/baz,/bar",
		},
	}

	for _, tc := range tests {
		t.Run("IgnoreGlobs.String", func(t *testing.T) {
			ignoreGlob := IgnoreGlobs(tc.globs)

			result := ignoreGlob.String()

			if result != tc.want {
				t.Errorf(`IgnoreGlob.String() = "%s", want "%s"`, result, tc.want)
			}
		})
	}
}

func TestPathArg(t *testing.T) {
	type test struct {
		value    string
		fileInfo os.FileInfo

		isDir  bool
		isFile bool
	}

	tempDir := t.TempDir()
	tempDirStat, _ := os.Stat(tempDir)

	file1, _ := os.CreateTemp(tempDir, "")
	fileStat1, _ := file1.Stat()

	tests := []test{
		{
			value:    file1.Name(),
			fileInfo: fileStat1,
			isFile:   true,
		},
		{
			value:    tempDir,
			fileInfo: tempDirStat,
			isDir:    true,
		},
		{
			value: "/invalid",
		},
	}

	for _, tc := range tests {
		t.Run("PathArg", func(t *testing.T) {
			pathArg := PathArg{Value: tc.value, fileInfo: tc.fileInfo}

			if pathArg.IsDir() != tc.isDir {
				t.Errorf(`PathArg.IsDir() = "%t", want "%t"`, pathArg.IsDir(), tc.isDir)
			}

			if pathArg.IsFile() != tc.isFile {
				t.Errorf(`PathArg.IsFile() = "%t", want "%t"`, pathArg.IsFile(), tc.isFile)
			}
		})
	}
}
