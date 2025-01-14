package input

import (
	"testing"
)

type validationInput struct {
    args Args
    usage string
    literalMode bool
    insensitiveMode bool
}

type test []struct{
    name string
    input validationInput
    expectError bool
}

func TestValidate(t *testing.T) {
    tests := test{
        // Success scenarios
        {
            name: "Valid subject, search, replace. Literal and insensitive mode false",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalMode: false,
                insensitiveMode: false,
            },
            expectError: false,
        },
        {
            name: "Valid subject, search, replace. Literal mode true",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalMode: true,
                insensitiveMode: false,
            },
            expectError: false,
        },
        {
            name: "Valid subject, search, replace. Insensitive mode true",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalMode: false,
                insensitiveMode: true,
            },
            expectError: false,
        },
        // Error scenarios
        {
            name: "Insensitive and literal mode true",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalMode: true,
                insensitiveMode: true,
            },
            expectError: true,
        },
        {
            name: "No Subject",
            input: validationInput{
                args: Args{Subject: "", Search: "Foo", Replace: "Baz"},
                usage: "",
                literalMode: false,
                insensitiveMode: false,
            },
            expectError: true,
        },
        {
            name: "No search",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "", Replace: "Baz"},
                usage: "",
                literalMode: false,
                insensitiveMode: false,
            },
            expectError: true,
        },
        {
            name: "No replace",
            input: validationInput{
                args: Args{Subject: "Foo Bar", Search: "Foo", Replace: ""}, usage: "", literalMode: false,                insensitiveMode: false,
            },
            expectError: true,
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            result := Validate(
                tc.input.args,
                tc.input.usage,
                tc.input.literalMode,
                tc.input.insensitiveMode,
            )

            if result != nil && !tc.expectError {
                t.Errorf("not expected error, got error. input = %+v", tc.input);
            }

            if result == nil && tc.expectError {
                t.Errorf("expected error, did not get error. input = %+v", tc.input);
            }
        })
    }

}

