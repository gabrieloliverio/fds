package input

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
        args Args
        usage string
        literalFlag bool
        insensitiveFlag bool
        confirmFlag bool
    }

    type test []struct{
        name string
        input validationInput
        expectError bool
    }

    tests := test{
        // Success scenarios
        {
            name: "Valid subject, search, replace. Literal and insensitive flag false",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalFlag: false,
                insensitiveFlag: false,
                confirmFlag: false,
            },
            expectError: false,
        },
        {
            name: "Valid subject, search, replace. Literal flag true",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalFlag: true,
                insensitiveFlag: false,
                confirmFlag: false,
            },
            expectError: false,
        },
        {
            name: "Valid subject, search, replace. Insensitive flag true",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalFlag: false,
                insensitiveFlag: true,
                confirmFlag: false,
            },
            expectError: false,
        },
        {
            name: "Valid subject (file content), search, replace. Confirm flag true",
            input: validationInput{
                args: Args{File: fileArg{Path: "./foo"}, Subject: "Foo", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalFlag: false,
                insensitiveFlag: false,
                confirmFlag: true,
            },
            expectError: false,
        },
        // Error scenarios
        {
            name: "Insensitive and literal flag true",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalFlag: true,
                insensitiveFlag: true,
            },
            expectError: true,
        },
        {
            name: "No Subject",
            input: validationInput{
                args: Args{Subject: "", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalFlag: false,
                insensitiveFlag: false,
            },
            expectError: true,
        },
        {
            name: "No search",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "", Replace: "Baz"},
                usage: "",
                literalFlag: false,
                insensitiveFlag: false,
            },
            expectError: true,
        },
        {
            name: "No replace",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: ""},
                usage: "",
                literalFlag: false,
                insensitiveFlag: false,
            },
            expectError: true,
        },
        {
            name: "Confirm flag without file",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: ""},
                usage: "",
                literalFlag: false,
                insensitiveFlag: true,
            },
            expectError: true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            err := Validate(
                tc.input.args,
                tc.input.literalFlag,
                tc.input.insensitiveFlag,
                tc.input.confirmFlag,
            )

            if err != nil && !tc.expectError {
                t.Errorf("TestValidate(%+v) did not expected error, got error %s", tc.input, err);
            }

            if err == nil && tc.expectError {
                t.Errorf("TestValidate(%+v) expected error, did not get error", tc.input);
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

func TestReadArgs_Stdin_No_Parameters_Return_Error(t *testing.T) {
    stdin := createTempFile(os.TempDir(), "my subject", t)

	_, err := ReadArgs(stdin, []string{})

	if err == nil {
		t.Errorf(`ReadArgs() expected error, did not get error"`)
	}
}

func TestReadArgs_PositionalArguments(t *testing.T) {
    stdin := createTempFile(os.TempDir(), "", t)

	want := Args{Subject: "my subject", Search: "search", Replace: "replace"}
	result, _ := ReadArgs(stdin, []string{"my subject", "search", "replace"})

	if result != want {
		t.Errorf(`ReadArgs() = "%+v", want "%+v"`, result, want)
	}
}

func TestReadArgs_PositionalArguements_No_Parameters_Return_Error(t *testing.T) {
    stdin := createTempFile(os.TempDir(), "", t)

	_, err := ReadArgs(stdin, []string{})

	if err == nil {
		t.Errorf(`ReadArgs() expected error, did not get error"`)
	}
}
