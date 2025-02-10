package match

import (
	"fmt"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/gabrieloliverio/fds/input"
)

type MatchString struct {
	Search string
	Replace string
	Before string
	After string

	IndexStart int
	IndexEnd int
}

func ConfirmMatch(match MatchString, filename string, stdin *os.File) (rune, error) {
    red := color.New(color.FgHiRed, color.Bold, color.Italic)
    green := color.New(color.FgHiGreen, color.Bold)

    fmt.Printf("File %s\n", filename)
    fmt.Printf("%s%s%s%s\n", match.Before, red.Sprint(match.Search), green.Sprint(match.Replace), match.After)

	ret, err := input.Confirm(stdin)

	if err != nil {
		return 0, err
	}

	return ret, nil
}

func FindStringOrPattern(search, replace, subject string, flags map[string]bool, bytesInDiff int) []MatchString {
	searchWithModifiers := search 

	if flags["literal"] {
		searchWithModifiers = regexp.QuoteMeta(search)
	}

	if flags["insensitive"] {
		searchWithModifiers = "(?i)" + search
	}

	pattern := regexp.MustCompile(searchWithModifiers)
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
	    	Search: string(matchString),
	    	Replace: replace,
	    	Before: string(subjectSlice[leftmostIndex:indexes[0]]),
	    	After: string(subjectSlice[indexes[1]:rightmostIndex]),
	    	IndexStart: indexes[0],
	    	IndexEnd: indexes[1],
	    })
    }

	return matches
}
