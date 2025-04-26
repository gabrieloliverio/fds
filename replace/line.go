package replace

import (
	"regexp"
	"slices"
)

type LineReplacer struct {
	flags map[string]bool
}

func NewReplacer(flags map[string]bool) LineReplacer {
	return LineReplacer{flags: flags}
}

func (s LineReplacer) Replace(subject, search, replace string) string {
	searchWithModifiers := search 

	if s.flags["literal"] {
		searchWithModifiers = regexp.QuoteMeta(search)
	}

	if s.flags["insensitive"] {
		searchWithModifiers = "(?i)" + search
	}

	pattern := regexp.MustCompile(searchWithModifiers)

	return pattern.ReplaceAllString(subject, replace)
}

/**
 * ReplaceStringRange replaces a given string or pattern when found in a range defined in `stringRange`
 * All other matches found out of the supplied range are ignored and therefore, not replaced 
 */
func (r LineReplacer) ReplaceStringRange(subject, search, replace string, stringRange [2]int) string {
	var prepend, append []byte

	searchWithModifiers := search

	if r.flags["literal"] {
		searchWithModifiers = regexp.QuoteMeta(search)
	}

	if r.flags["insensitive"] {
		searchWithModifiers = "(?i)" + search
	}

	pattern := regexp.MustCompile(searchWithModifiers)
	subjectSubstring := []byte(subject)[stringRange[0]:stringRange[1]]
	replaced := pattern.ReplaceAll(subjectSubstring, []byte(replace))

	prepend = []byte(subject)[0:stringRange[0]]
	append = []byte(subject)[stringRange[1]:]

	return string(slices.Concat(prepend, replaced, append))
}

