package replace

import (
	"io"
	"os"
	"path"
	"regexp"
	"testing"
)

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

func TestLineReplacer_Replace(t *testing.T) {
	var tests = []struct{
		name string
		subject string
		search string
		replace string
		flags map[string]bool
		want *regexp.Regexp
	}{
		{
			name: "replace with literal string",
			search: "text",
			replace: "replacement",
			subject: "this is some text, this is some other text",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": true},
			want: regexp.MustCompile("this is some replacement, this is some other replacement"),
		},
		{
			name: "replace with regular expression",
			search: "t.xt",
			replace: "replacement",
			subject: "this is some text",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},
			want: regexp.MustCompile("this is some replacement"),
		},
		{
			name: "replace with regular expression ignoring case",
			search: "Text",
			replace: "replacement",
			subject: "this is some text",
			flags: map[string]bool{"insensitive": true, "confirm": false, "literal": false},
			want: regexp.MustCompile("this is some replacement"),
		},
		{
			name: "regular expression not match",
			search: "(text)",
			replace: "other $1",
			subject: "this is some text",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},
			want: regexp.MustCompile("this is some other text"),
		},
		{
			name: "regular expression capturing group",
			search: "<fooo>",
			replace: "replacement",
			subject: "this is some text, this is some other text",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": true},
			want: regexp.MustCompile("this is some text, this is some other text"),
		},
		{
			name: "literal not match",
			search: "<fooo>",
			replace: "replacement",
			subject: "this is some text, this is some other text",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": true},
			want: regexp.MustCompile("this is some text, this is some other text"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			replacer := NewReplacer(tc.flags)
			pattern := replacer.CompilePattern(tc.search)
			result := replacer.Replace(pattern, tc.subject, tc.replace)

			if !tc.want.MatchString(result) {
				t.Errorf(`Replacer.Replace() = "%q", want "%#v"`, result, tc.want)
			}
		})
	}
}


func TestLineReplacer_ReplaceStringRange(t *testing.T) {
	var tests = []struct{
		name string
		subject string
		search string
		replace string
		stringRange [2]int
		flags map[string]bool
		want string
	}{
		{
			name: "string range starting in 0",
			subject: "this is some text, this is the rest of the text",
			search: "text",
			replace: "replacement",
			stringRange: [2]int{0, 17},
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},
			want: "this is some replacement, this is the rest of the text",
		},
		{
			name: "string range starting and ending in the middle",
			subject: "this is some text, this is the rest of the text",
			search: "this",
			replace: "that",
			stringRange: [2]int{17, 35},
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},
			want: "this is some text, that is the rest of the text",
		},
		{
			name: "string range starting in the middle and ending in the end",
			subject: "this is some text, this is the rest of the text",
			search: "text",
			replace: "replacement",
			stringRange: [2]int{42, 47},
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},
			want: "this is some text, this is the rest of the replacement",
		},
		{
			name: "no match",
			subject: "this is some text, this is the rest of the text",
			search: "banana",
			replace: "that",
			stringRange: [2]int{42, 47},
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},
			want: "this is some text, this is the rest of the text",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			replacer := NewReplacer(tc.flags)
			pattern := replacer.CompilePattern(tc.search)
			result := replacer.ReplaceStringRange(pattern, tc.subject, tc.replace, tc.stringRange)

			if result != tc.want {
				t.Errorf("ReplaceMatch() = %s, want %s", result, tc.want)
			}
		})
	}
}
