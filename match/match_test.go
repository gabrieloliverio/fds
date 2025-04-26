package match

import (
	"reflect"
	"testing"
)

func TestFindStringOrPattern(t *testing.T) {
	var tests = []struct{
		name string
		subject string
		search string
		replace string
		flags map[string]bool
		want []MatchString
	}{
		{
			name: "no match",
			search: "text",
			subject: "foo",
			replace: "replacement",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},

			want: []MatchString{},
		},
		{
			name: "single occurrence - no characters left and no characters right",
			search: "text",
			subject: "text",
			replace: "replacement",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},

			want: []MatchString{
			    {
				    Search: "text",
				    Replace: "replacement",
				    Before: "",
				    After: "",
				    IndexStart: 0,
				    IndexEnd: 4,
				},
			},
		},
		{
			name: "single occurrence - few characters left and few characters right",
			search: "text",
			subject: "this is some text and that's it",
			replace: "replacement",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},

			want: []MatchString{
			    {
				    Search: "text",
				    Replace: "replacement",
				    Before: "this is some ",
				    After: " and that's it",
				    IndexStart: 13,
				    IndexEnd: 17,
			    },
			},
		},
		{
			name: "single occurrence - lots of characters left and lots of characters right",
			search: "text",
			subject: "this is some unnecessarily long text, and now I ran out of criativity",
			replace: "replacement",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},

			want: []MatchString{
			    {
				    Search: "text",
				    Replace: "replacement",
				    Before: " unnecessarily long ",
				    After: ", and now I ran out ",
				    IndexStart: 32,
				    IndexEnd: 36,
				},
			},
		},
		{
			name: "single occurrence - exactly 20 characters left and right",
			search: "text",
			subject: "some random 20-char text some random 20-char",
			replace: "replacement",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},

			want: []MatchString{
			    {
				    Search: "text",
				    Replace: "replacement",
				    Before: "some random 20-char ",
				    After: " some random 20-char",
				    IndexStart: 20,
				    IndexEnd: 24,
				},
			},
		},
		{
			name: "single occurrence - exactly 20 characters left and right with literal flag",
			search: "text",
			subject: "some random 20-char text some random 20-char",
			replace: "$1",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": true},

			want: []MatchString{
			    {
				    Search: "text",
				    Replace: "$1",
				    Before: "some random 20-char ",
				    After: " some random 20-char",
				    IndexStart: 20,
				    IndexEnd: 24,
				},
			},
		},
		{
			name: "single occurrence - exactly 20 characters left and right matching group",
			search: "(text)",
			subject: "some random 20-char text some random 20-char",
			replace: "other $1",
			flags: map[string]bool{"insensitive": true, "confirm": false, "literal": false},

			want: []MatchString{
			    {
				    Search: "text",
				    Replace: "other $1",
				    Before: "some random 20-char ",
				    After: " some random 20-char",
				    IndexStart: 20,
				    IndexEnd: 24,
				},
			},
		},
		{
			name: "single occurrence - exactly 20 characters left and right with insensitive flag",
			search: "TEXT",
			subject: "some random 20-char text some random 20-char",
			replace: "replacement",
			flags: map[string]bool{"insensitive": true, "confirm": false, "literal": false},

			want: []MatchString{
			    {
				    Search: "text",
				    Replace: "replacement",
				    Before: "some random 20-char ",
				    After: " some random 20-char",
				    IndexStart: 20,
				    IndexEnd: 24,
				},
			},
		},
		{
			name: "single occurrence - with regular expression",
			search: ".ext",
			subject: "some random 20-char text some random 20-char",
			replace: "replacement",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},

			want: []MatchString{
			    {
				    Search: "text",
				    Replace: "replacement",
				    Before: "some random 20-char ",
				    After: " some random 20-char",
				    IndexStart: 20,
				    IndexEnd: 24,
				},
			},
		},
		{
			name: "multiple matches",
			search: "text",
			subject: "this is one text, this is another text, this is yet another text",
			replace: "replacement",
			flags: map[string]bool{"insensitive": false, "confirm": false, "literal": false},

			want: []MatchString{
			    {
				    Search: "text",
				    Replace: "replacement",
				    Before: "this is one ",
				    After: ", this is another te",
				    IndexStart: 12,
				    IndexEnd: 16,
				},
			    {
				    Search: "text",
				    Replace: "replacement",
				    Before: "xt, this is another ",
				    After: ", this is yet anothe",
				    IndexStart: 34,
				    IndexEnd: 38,
				},
			    {
				    Search: "text",
				    Replace: "replacement",
				    Before: "this is yet another ",
				    After: "",
				    IndexStart: 60,
				    IndexEnd: 64,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FindStringOrPattern(tc.search, tc.replace, tc.subject, tc.flags, 20)
			
			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("FindStringOrPattern() = %+v, want %+v", result, tc.want)
			}
		})
	}
}


