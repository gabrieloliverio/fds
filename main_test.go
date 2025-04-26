package main

import (
	"regexp"
	"testing"
)

func TestReplaceStringOrPattern_LiteralString(t *testing.T) {
	search := "text"
	replace := "replacement"
	subject := "this is some text, this is some other text"
	result := replaceStringOrPattern(search, replace, subject, true, false)

	want := regexp.MustCompile("this is some replacement, this is some other replacement")

	if !want.MatchString(result) {
		t.Fatalf(`replaceStringOrPattern() = %q, want %#q`, result, want)
	}
}

func TestReplaceStringOrPattern_RegEx(t *testing.T) {
	searchPattern := "t.xt"
	replace := "replacement"
	subject := "this is some text"
	result := replaceStringOrPattern(searchPattern, replace, subject, false, false)

	want := regexp.MustCompile("this is some replacement")

	if !want.MatchString(result) {
		t.Fatalf(`replaceStringOrPattern() = %q, want %#q`, result, want)
	}
}

func TestReplaceStringOrPattern_RegExIgnoringCase(t *testing.T) {
	searchPattern := "Text"
	replace := "replacement"
	subject := "this is some text"
	result := replaceStringOrPattern(searchPattern, replace, subject, false, true)

	want := regexp.MustCompile("this is some replacement")

	if !want.MatchString(result) {
		t.Fatalf(`replaceStringOrPattern() = %q, want %#q`, result, want)
	}
}

func TestReplaceStringOrPattern_RegExCapturingGroup(t *testing.T) {
	searchPattern := "(text)"
	replace := "other $1"
	subject := "this is some text"
	result := replaceStringOrPattern(searchPattern, replace, subject, false, false)

	want := regexp.MustCompile("this is some other text")

	if !want.MatchString(result) {
		t.Fatalf(`replaceStringOrPattern() = %q, want %#q`, result, want)
	}
}

func TestReplaceStringOrPattern_RegExNotMatch(t *testing.T) {
	searchPattern := "<fooo>"
	replace := "replacement"
	subject := "this is some text, this is some other text"
	result := replaceStringOrPattern(searchPattern, replace, subject, false, false)

	want := regexp.MustCompile("this is some text, this is some other text")

	if !want.MatchString(result) {
		t.Fatalf(`replaceStringOrPattern() = %q, want %#q`, result, want)
	}
}

func TestReplaceStringOrPattern_LiteralNotMatch(t *testing.T) {
	searchPattern := "<fooo>"
	replace := "replacement"
	subject := "this is some text, this is some other text"
	result := replaceStringOrPattern(searchPattern, replace, subject, true, false)

	want := regexp.MustCompile("this is some text, this is some other text")

	if !want.MatchString(result) {
		t.Fatalf(`replaceStringOrPattern() = %q, want %#q`, result, want)
	}
}
