package main

import (
	"regexp"
	"testing"
)

func TestFindReplace(t *testing.T) {
	searchPattern := regexp.MustCompile("text")
	replace := "replacement"
	subject := "this is some text, this is some other text"
	result := findReplace(searchPattern, replace, subject)

	want := regexp.MustCompile("this is some replacement, this is some other replacement")

	if !want.MatchString(result) {
		t.Fatalf(`search() = %q, want %#q`, result, want)
	}

}
func TestFindReplaceNotMatch(t *testing.T) {
	searchPattern := regexp.MustCompile("<fooo>")
	replace := "replacement"
	subject := "this is some text, this is some other text"
	result := findReplace(searchPattern, replace, subject)

	want := regexp.MustCompile("this is some text, this is some other text")

	if !want.MatchString(result) {
		t.Fatalf(`search() = %q, want %#q`, result, want)
	}

}
