package fds

import (
	"regexp"
	"slices"
)

type Replacer interface {
	HasFlag(flag string) bool
}

type LineReplacer struct {
	flags map[string]bool

	search  *regexp.Regexp
	replace string
}

func NewLineReplacer(search, replace string, flags map[string]bool) LineReplacer {
	replacer := LineReplacer{flags: flags, replace: replace}
	replacer.search = replacer.compilePattern(search)

	return replacer
}

func (s LineReplacer) compilePattern(search string) *regexp.Regexp {
	searchWithModifiers := search

	if s.flags["literal"] {
		searchWithModifiers = regexp.QuoteMeta(search)
	}

	if s.flags["insensitive"] {
		searchWithModifiers = "(?i)" + search
	}

	return regexp.MustCompile(searchWithModifiers)
}

func (s LineReplacer) Replace(subject string) (result string, replaced bool) {
	result = s.search.ReplaceAllString(subject, s.replace)

	if result != subject {
		replaced = true
	}

	return
}

func (s LineReplacer) HasFlag(flag string) bool {
	return s.flags[flag]
}

/**
 * ReplaceStringRange replaces a given string or pattern when found in a range defined in `stringRange`
 * All other matches found out of the supplied range are ignored and therefore, not replaced
 */
func (r LineReplacer) ReplaceStringRange(subject string, stringRange [2]int) string {
	var prepend, append []byte

	subjectSubstring := []byte(subject)[stringRange[0]:stringRange[1]]
	replaced := r.search.ReplaceAll(subjectSubstring, []byte(r.replace))

	prepend = []byte(subject)[0:stringRange[0]]
	append = []byte(subject)[stringRange[1]:]

	return string(slices.Concat(prepend, replaced, append))
}
