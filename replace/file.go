package replace

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gabrieloliverio/fds/input"
	"github.com/gabrieloliverio/fds/match"
)

const (
    confirmText = "[y]es [n]o [a]ll q[uit]: "

    ConfirmYes = 'y'
    ConfirmNo = 'n'
    ConfirmAll = 'a'
    ConfirmQuit = 'q'
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

/*
 * ReplaceInFile replaces a given pattern when found in `inputFile`. Lines are written into `outputFile`
 */
func (r FileReplacer) ReplaceInFile(pattern *regexp.Regexp, replace string, stdin *os.File, confirmAnswer *input.ConfirmAnswer) error {
	if r.flags["confirm"] {
		err := r.confirmAndReplace(pattern, replace, stdin, confirmAnswer)

		return err
	}

	return r.replaceAll(pattern, replace)
}

func (r FileReplacer) replaceAll(pattern *regexp.Regexp, replace string) error {
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

		replaced := r.LineReplacer.Replace(pattern, line, replace)

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

func (r FileReplacer) confirmAndReplace(pattern *regexp.Regexp, replace string, stdin *os.File, confirmAnswer *input.ConfirmAnswer) error {
	var replacer = NewReplacer(r.flags)
	var lineNumber int
	var confirmedAll, confirmedQuit bool

	// feedback confirmedAll that propagated all the way back to the main function
	confirmedAll = *confirmAnswer == ConfirmAll

	if inputFileStat, _ := r.inputFile.Stat(); inputFileStat.Size() == 0 {
		return nil
	}

	reader := bufio.NewReader(r.inputFile)
	writer := bufio.NewWriter(r.outputFile)

	for {
		line, err := reader.ReadString('\n')
		lineNumber++

		if err != nil && err != io.EOF {
			return fmt.Errorf("Error while reading file: %s", err)
		}

		if confirmedAll {
			line = replacer.Replace(pattern, line, replace)
		}

		if !confirmedAll && !confirmedQuit  {
			matches := match.FindStringOrPattern(pattern, replace, line, 50)

			line = r.confirmMatches(matches, line, replace, pattern, lineNumber, stdin, confirmAnswer)
		}

		_, errWrite := writer.WriteString(line)

		if errWrite != nil  {
			return fmt.Errorf("Error while writing temporary file: %s", err)
		}

		if err != nil && err == io.EOF {
			break
		}
	}

	writer.Flush()

	return nil
}

func (r FileReplacer) confirmMatches(matches []match.MatchString, line, replace string, pattern *regexp.Regexp, lineNumber int, stdin *os.File, confirmAnswer *input.ConfirmAnswer) string {
	var answer rune
	var err error

	confirmedQuit := rune(*confirmAnswer) == ConfirmQuit
	confirmedAll := rune(*confirmAnswer) == ConfirmAll
	
	replacedLine := line

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
			answer = ConfirmYes
		} else {
			answer, err = match.ConfirmMatch(thisMatch, r.inputFile.Name(), lineNumber, stdin)
			*confirmAnswer = input.ConfirmAnswer(answer)
		}

		if err != nil {
			fmt.Println(err)
		}

		switch answer {
		case ConfirmYes:
			replacedLine = r.LineReplacer.ReplaceStringRange(pattern, replacedLine, replace, stringRange)
		case ConfirmNo:
			// Nothing to do
		case ConfirmAll:
			replacedLine = r.LineReplacer.ReplaceStringRange(pattern, replacedLine, replace, stringRange)
			confirmedAll = true
		default:
			confirmedQuit = true
		}
	}

	return replacedLine
}

func resolveInputFile(path string) (*os.File, error) {
	fileStat, _ := os.Lstat(path)
	inputFilePath := path

	if fileStat.Mode().Type() == os.ModeSymlink.Type() {
		inputFilePath, _ = filepath.EvalSymlinks(path)
		inputFilePath, _ = filepath.Abs(inputFilePath)
	}

	return os.OpenFile(inputFilePath, os.O_RDONLY, fileStat.Mode())
}

func OpenInputAndTempFile(inputPath string) (inputFile, tmpFile *os.File, err error) {
	inputFile, err = resolveInputFile(inputPath)

	if err != nil {
		return nil, nil, err
	}

	tmpFile, err = os.CreateTemp("", filepath.Base(inputPath))

	return
}
