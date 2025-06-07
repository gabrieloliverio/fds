package replace

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/gabrieloliverio/fds/input"
)

type MatchString struct {
	Search     string
	Replace    string
	Before     string
	After      string
	LineNumber int

	IndexStart int
	IndexEnd   int
}

func ConfirmMatch(match MatchString, filename string, lineNumber int, stdin *os.File) (rune, error) {
	red := color.New(color.FgHiRed, color.Bold, color.Italic)
	green := color.New(color.FgHiGreen, color.Bold)

	fmt.Printf("File\t%s\n", filename)
	fmt.Printf("%d\t%s%s%s%s\n", lineNumber, match.Before, red.Sprint(match.Search), green.Sprint(match.Replace), match.After)

	confirmText := "[y]es [n]o [a]ll q[uit]"
	valid := []rune{'y', 'n', 'a', 'q'}
	ret, err := input.Confirm(stdin, confirmText, valid)

	if err != nil {
		return 0, err
	}

	fmt.Println()

	return ret, nil
}

func FindStringOrPattern(pattern *regexp.Regexp, replace, subject string, bytesInDiff int) []MatchString {
	allIndexes := pattern.FindAllStringIndex(subject, -1)

	matches := make([]MatchString, 0)

	if allIndexes == nil {
		return matches
	}

	for _, indexes := range allIndexes {
		matchString := []byte(subject[indexes[0]:indexes[1]])

		leftmostIndex := indexes[0] - bytesInDiff
		rightmostIndex := indexes[1] + bytesInDiff

		if leftmostIndex < 0 {
			leftmostIndex = 0
		}

		if rightmostIndex > len(subject) {
			rightmostIndex = len(subject)
		}

		subjectSlice := []byte(subject)

		matches = append(matches, MatchString{
			Search:     string(matchString),
			Replace:    replace,
			Before:     string(subjectSlice[leftmostIndex:indexes[0]]),
			After:      string(subjectSlice[indexes[1]:rightmostIndex]),
			IndexStart: indexes[0],
			IndexEnd:   indexes[1],
		})
	}

	return matches
}
