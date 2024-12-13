package main

import (
	"regexp"
	"testing"
)

func TestFindReplace_ReplaceAsSimpleString(t *testing.T) {
	searchPattern := regexp.MustCompile("text")
	replace := "replacement"
	subject := "this is some text, this is some other text"
	result := replacePattern(searchPattern, replace, subject)

	want := regexp.MustCompile("this is some replacement, this is some other replacement")

	if !want.MatchString(result) {
		t.Fatalf(`search() = %q, want %#q`, result, want)
	}
}

func TestFindReplace_CapturingGroup(t *testing.T) {
	searchPattern := regexp.MustCompile("(text)")
	replace := "other $1"
	subject := "this is some text"
	result := replacePattern(searchPattern, replace, subject)

	want := regexp.MustCompile("this is some other text")

	if !want.MatchString(result) {
		t.Fatalf(`search() = %q, want %#q`, result, want)
	}
}

func TestFindReplaceNotMatch(t *testing.T) {
	searchPattern := regexp.MustCompile("<fooo>")
	replace := "replacement"
	subject := "this is some text, this is some other text"
	result := replacePattern(searchPattern, replace, subject)

	want := regexp.MustCompile("this is some text, this is some other text")

	if !want.MatchString(result) {
		t.Fatalf(`search() = %q, want %#q`, result, want)
	}

}
