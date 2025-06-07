package replace

import (
	"reflect"
	"regexp"
	"testing"
)

func TestFindStringOrPattern(t *testing.T) {
	var tests = []struct {
		name    string
		subject string
		pattern *regexp.Regexp
		replace string
		want    []MatchString
	}{
		{
			name:    "no match",
			pattern: regexp.MustCompile("text"),
			subject: "foo",
			replace: "replacement",

			want: []MatchString{},
		},
		{
			name:    "single occurrence - no characters left and no characters right",
			pattern: regexp.MustCompile("text"),
			subject: "text",
			replace: "replacement",

			want: []MatchString{
				{
					Search:     "text",
					Replace:    "replacement",
					Before:     "",
					After:      "",
					IndexStart: 0,
					IndexEnd:   4,
				},
			},
		},
		{
			name:    "single occurrence - few characters left and few characters right",
			pattern: regexp.MustCompile("text"),
			subject: "this is some text and that's it",
			replace: "replacement",

			want: []MatchString{
				{
					Search:     "text",
					Replace:    "replacement",
					Before:     "this is some ",
					After:      " and that's it",
					IndexStart: 13,
					IndexEnd:   17,
				},
			},
		},
		{
			name:    "single occurrence - lots of characters left and lots of characters right",
			pattern: regexp.MustCompile("text"),
			subject: "this is some unnecessarily long text, and now I ran out of criativity",
			replace: "replacement",

			want: []MatchString{
				{
					Search:     "text",
					Replace:    "replacement",
					Before:     " unnecessarily long ",
					After:      ", and now I ran out ",
					IndexStart: 32,
					IndexEnd:   36,
				},
			},
		},
		{
			name:    "single occurrence - exactly 20 characters left and right",
			pattern: regexp.MustCompile("text"),
			subject: "some random 20-char text some random 20-char",
			replace: "replacement",

			want: []MatchString{
				{
					Search:     "text",
					Replace:    "replacement",
					Before:     "some random 20-char ",
					After:      " some random 20-char",
					IndexStart: 20,
					IndexEnd:   24,
				},
			},
		},
		{
			name:    "single occurrence - exactly 20 characters left and right literal",
			pattern: regexp.MustCompile(regexp.QuoteMeta("text")),
			subject: "some random 20-char text some random 20-char",
			replace: "$1",

			want: []MatchString{
				{
					Search:     "text",
					Replace:    "$1",
					Before:     "some random 20-char ",
					After:      " some random 20-char",
					IndexStart: 20,
					IndexEnd:   24,
				},
			},
		},
		{
			name:    "single occurrence - exactly 20 characters left and right matching group",
			pattern: regexp.MustCompile("(?i)(text)"),
			subject: "some random 20-char text some random 20-char",
			replace: "other $1",

			want: []MatchString{
				{
					Search:     "text",
					Replace:    "other $1",
					Before:     "some random 20-char ",
					After:      " some random 20-char",
					IndexStart: 20,
					IndexEnd:   24,
				},
			},
		},
		{
			name:    "single occurrence - exactly 20 characters left and right insensitive",
			pattern: regexp.MustCompile("(?i)TEXT"),
			subject: "some random 20-char text some random 20-char",
			replace: "replacement",

			want: []MatchString{
				{
					Search:     "text",
					Replace:    "replacement",
					Before:     "some random 20-char ",
					After:      " some random 20-char",
					IndexStart: 20,
					IndexEnd:   24,
				},
			},
		},
		{
			name:    "single occurrence - with regular expression",
			pattern: regexp.MustCompile(".ext"),
			subject: "some random 20-char text some random 20-char",
			replace: "replacement",

			want: []MatchString{
				{
					Search:     "text",
					Replace:    "replacement",
					Before:     "some random 20-char ",
					After:      " some random 20-char",
					IndexStart: 20,
					IndexEnd:   24,
				},
			},
		},
		{
			name:    "multiple matches",
			pattern: regexp.MustCompile("text"),
			subject: "this is one text, this is another text, this is yet another text",
			replace: "replacement",

			want: []MatchString{
				{
					Search:     "text",
					Replace:    "replacement",
					Before:     "this is one ",
					After:      ", this is another te",
					IndexStart: 12,
					IndexEnd:   16,
				},
				{
					Search:     "text",
					Replace:    "replacement",
					Before:     "xt, this is another ",
					After:      ", this is yet anothe",
					IndexStart: 34,
					IndexEnd:   38,
				},
				{
					Search:     "text",
					Replace:    "replacement",
					Before:     "this is yet another ",
					After:      "",
					IndexStart: 60,
					IndexEnd:   64,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FindStringOrPattern(tc.pattern, tc.replace, tc.subject, 20)

			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("FindStringOrPattern() = %+v, want %+v", result, tc.want)
			}
		})
	}
}
