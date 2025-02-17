package replace

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/gabrieloliverio/fds/input"
	"github.com/gabrieloliverio/fds/match"
)

type FileReplacer struct {
	LineReplacer

	inputFile *os.File
	outputFile *os.File
}

func NewFileReplacer(inputFile, outputFile *os.File, flags map[string]bool) FileReplacer {
	return FileReplacer{
		inputFile: inputFile,
		outputFile: outputFile,
		LineReplacer: LineReplacer{flags: flags},
	}
}

/**
 * ReplaceInFile replaces a given pattern when found in `inputFile`. Lines are written into `outputFile`
 */
func (r FileReplacer) ReplaceInFile(search, replace string, stdin *os.File) error {
	if r.flags["confirm"] {
		return r.confirmAndReplace(search, replace, stdin)
	}

	return r.replaceAll(search, replace)
}

func (r FileReplacer) replaceAll(search, replace string) error {
	var err error

	if inputFileStat, _ := r.inputFile.Stat(); inputFileStat.Size() == 0 {
		return nil
	}

	reader := bufio.NewReader(r.inputFile)
	writer := bufio.NewWriter(r.outputFile)

	for {
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			return fmt.Errorf("Error while reading file: %s", err)
		}

		replaced := r.LineReplacer.Replace(line, search, replace)

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

func (r FileReplacer) confirmAndReplace(search, replace string, stdin *os.File) error {
	var err error
	var confirmedAll, confirmedQuit bool
	var replacer = NewReplacer(r.flags)

	if inputFileStat, _ := r.inputFile.Stat(); inputFileStat.Size() == 0 {
		return nil
	}

	reader := bufio.NewReader(r.inputFile)
	writer := bufio.NewWriter(r.outputFile)

	for {
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			return fmt.Errorf("Error while reading file: %s", err)
		}

		if err != nil && err == io.EOF {
			break
		}

		if confirmedAll {
			line = replacer.Replace(line, search, replace)
		}

		if !confirmedAll && !confirmedQuit  {
			matches := match.FindStringOrPattern(search, replace, line, r.flags, 50)

			line, confirmedAll, confirmedQuit = r.confirmMatches(matches, line, search, replace, stdin)
		}

		_, errWrite := writer.WriteString(line)

		if errWrite != nil  {
			return fmt.Errorf("Error while writing temporary file: %s", err)
		}
	}

	writer.Flush()

	return err
}

func (r FileReplacer) confirmMatches(matches []match.MatchString, line, search, replace string, stdin *os.File) (replacedLine string, confirmedAll, confirmedQuit bool) {
	var answer rune
	var err error
	
	replacedLine = line

	for i, thisMatch := range matches {
		if confirmedQuit {
			continue
		}

		stringRange := [2]int{0, thisMatch.IndexEnd}
		if i > 0 {
			// Gets the previous one
			stringRange = [2]int{matches[i - 1].IndexEnd, thisMatch.IndexEnd}
		}

		if confirmedAll {
			answer = input.ConfirmYes
		} else {
			answer, err = match.ConfirmMatch(thisMatch, r.inputFile.Name(), stdin)
		}

		if err != nil {
			fmt.Println(err)
		}

		switch answer {
		case input.ConfirmYes:
			replacedLine = r.LineReplacer.ReplaceStringRange(line, search, replace, stringRange)
		case input.ConfirmNo:
			// Nothing to do
		case input.ConfirmAll:
			replacedLine = r.LineReplacer.ReplaceStringRange(line, search, replace, stringRange)
			confirmedAll = true
		default:
			// ConfirmedQuit
			// TODO: copy the remaining lines to tmp file and flush. Maybe next iteration?
			confirmedQuit = true
		}
	}

	return
}

