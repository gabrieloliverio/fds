package replace

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"

	"github.com/gabrieloliverio/fds/input"
	"github.com/gabrieloliverio/fds/match"
)

func ReplaceAllStringOrPattern(search, replace, subject string, flags map[string]bool) string {
	searchWithModifiers := search 

	if flags["literal"] {
		searchWithModifiers = regexp.QuoteMeta(search)
	}

	if flags["insensitive"] {
		searchWithModifiers = "(?i)" + search
	}

	pattern := regexp.MustCompile(searchWithModifiers)

	return pattern.ReplaceAllString(subject, replace)
}

/**
 * ReplaceStringRange replaces a given string or pattern when found in a range defined in `stringRange`
 * All other matches found out of the supplied range are ignored and therefore, not replaced 
 */
func ReplaceStringRange(a input.Args, stringRange [2]int, flags map[string]bool) string {
	var prepend, append []byte

	searchWithModifiers := a.Search

	if flags["literal"] {
		searchWithModifiers = regexp.QuoteMeta(a.Search)
	}

	if flags["insensitive"] {
		searchWithModifiers = "(?i)" + a.Search
	}

	pattern := regexp.MustCompile(searchWithModifiers)
	subjectSubstring := []byte(a.Subject)[stringRange[0]:stringRange[1]]
	replaced := pattern.ReplaceAll(subjectSubstring, []byte(a.Replace))

	prepend = []byte(a.Subject)[0:stringRange[0]]
	append = []byte(a.Subject)[stringRange[1]:]

	return string(slices.Concat(prepend, replaced, append))
}

/**
 * ReplaceInFile replaces a given string or pattern when found in `inputFile`. Lines are stored in `outputFile`
 */
func ReplaceInFile(inputFile, outputFile *os.File, args input.Args, flags map[string]bool) error {
	inputFileStat, _ := inputFile.Stat()
	var err error

	if inputFileStat.Size() == 0 {
		return nil
	}

	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)

	for {
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			return fmt.Errorf("Error while reading file: %s", err)
		}

		replaced := ReplaceAllStringOrPattern(args.Search, args.Replace, line, flags)

		_, errWrite := writer.WriteString(replaced)

		if errWrite != nil  {
			return fmt.Errorf("Error while writing temporary file: %s", err)
		}

		if err != nil && err == io.EOF {
			break
		}
	}

	writer.Flush()

	return err
}

func ConfirmAndReplace(inputFile, outputFile, stdin *os.File, args input.Args, flags map[string]bool) error {
	inputFileStat, _ := inputFile.Stat()
	var err error
	var confirmedAll, confirmedQuit bool

	if inputFileStat.Size() == 0 {
		return nil
	}

	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)

	for {
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			return fmt.Errorf("Error while reading file: %s", err)
		}

		if err != nil && err == io.EOF {
			break
		}

		if confirmedAll || !flags["confirm"] {
			line = ReplaceAllStringOrPattern(args.Search, args.Replace, line, flags)
		}

		if flags["confirm"] && !confirmedAll && !confirmedQuit  {
			matches := match.FindStringOrPattern(args.Search, args.Replace, line, flags, 50)

			for i, thisMatch := range matches {
				var answer rune
				var stringRange [2]int

				if confirmedQuit {
					continue
				}

				if i == 0 {
					stringRange = [2]int{0, thisMatch.IndexEnd}
				} else {
					stringRange = [2]int{matches[i - 1].IndexEnd, thisMatch.IndexEnd}
				}

				if confirmedAll {
					answer = input.ConfirmYes
				} else {
					answer, err = match.ConfirmMatch(thisMatch, inputFile.Name(), stdin)
				}

				if err != nil {
					fmt.Println(err)
				}

				args.Subject = line

				switch answer {
				case input.ConfirmYes:
					line = ReplaceStringRange(args, stringRange, flags)
				case input.ConfirmNo:
					// Nothing to do
				case input.ConfirmAll:
					line = ReplaceStringRange(args, stringRange, flags)
					confirmedAll = true
				default:
					// ConfirmedQuit
					// TODO: copy the remaining lines to tmp file and flush. Maybe next iteration?
					confirmedQuit = true
				}
			}
		}

		_, errWrite := writer.WriteString(line)

		if errWrite != nil  {
			return fmt.Errorf("Error while writing temporary file: %s", err)
		}
	}

	writer.Flush()

	return err
}

