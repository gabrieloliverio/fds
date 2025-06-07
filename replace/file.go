package replace

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gabrieloliverio/fds/input"
	"github.com/gabrieloliverio/fds/match"
)

const (
	confirmText = "[y]es [n]o [a]ll q[uit]: "

	ConfirmYes  = 'y'
	ConfirmNo   = 'n'
	ConfirmAll  = 'a'
	ConfirmQuit = 'q'
)

type FileReplacer struct {
	LineReplacer

	inputFilePath string
}

func NewFileReplacer(inputFilePath, search, replace string, flags map[string]bool) FileReplacer {
	replacer := FileReplacer{LineReplacer: LineReplacer{flags: flags, replace: replace}, inputFilePath: inputFilePath}
	replacer.search = replacer.compilePattern(search)

	return replacer
}

/*
 * ReplaceInFile replaces a given pattern when found in `inputFile`. Lines are written into `outputFile`
 */
func (r FileReplacer) Replace(stdin *os.File, confirmAnswer *input.ConfirmAnswer) (outputFile *os.File, fileChanged bool, err error) {
	if r.flags["confirm"] {
		outputFile, fileChanged, err = r.confirmAndReplace(stdin, confirmAnswer)

		return
	}

	return r.replaceAll()
}

func (r FileReplacer) replaceAll() (tmpFile *os.File, fileChanged bool, err error) {
	inputFile, err := openInputFile(r.inputFilePath)
	tmpFile, _ = os.CreateTemp("", filepath.Base(inputFile.Name()))
	writer := bufio.NewWriter(tmpFile)

	if err != nil {
		return
	}

	if inputFileStat, _ := inputFile.Stat(); inputFileStat.Size() == 0 {
		return
	}

	reader := bufio.NewReader(inputFile)

	for {
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			return nil, false, fmt.Errorf("Error while reading file: %s", err)
		}

		replacedLine, lineChanged := r.LineReplacer.Replace(line)

		if lineChanged {
			fileChanged = true
		}

		// TODO: Only create temp file when there is a change
		_, errWrite := writer.WriteString(replacedLine)

		if errWrite != nil {
			return nil, false, fmt.Errorf("Error while writing temporary file: %s", err)
		}

		if err != nil && err == io.EOF {
			break
		}
	}

	writer.Flush()

	return tmpFile, fileChanged, err
}

func (r FileReplacer) confirmAndReplace(stdin *os.File, confirmAnswer *input.ConfirmAnswer) (outputFile *os.File, fileChanged bool, err error) {
	var (
		lineNumber                  int
		confirmedAll, confirmedQuit bool
		lineChanged                 bool
	)

	// feedback confirmedAll that propagated all the way back to the main function
	confirmedAll = *confirmAnswer == ConfirmAll

	inputFile, err := openInputFile(r.inputFilePath)
	tmpFile, _ := os.CreateTemp("", filepath.Base(inputFile.Name()))
	writer := bufio.NewWriter(tmpFile)

	if err != nil {
		return
	}

	if inputFileStat, _ := inputFile.Stat(); inputFileStat.Size() == 0 {
		return
	}

	reader := bufio.NewReader(inputFile)

	for {
		line, err := reader.ReadString('\n')
		lineNumber++

		if err != nil && err != io.EOF {
			return nil, false, fmt.Errorf("Error while reading file: %s", err)
		}

		if confirmedAll {
			line, lineChanged = r.LineReplacer.Replace(line)
		}

		if !confirmedAll && !confirmedQuit {
			matches := match.FindStringOrPattern(r.search, r.replace, line, 50)

			line, lineChanged = r.confirmMatches(matches, line, lineNumber, stdin, confirmAnswer)
		}

		if lineChanged {
			fileChanged = true
		}

		_, errWrite := writer.WriteString(line)

		if errWrite != nil {
			return nil, false, fmt.Errorf("Error while writing temporary file: %s", err)
		}

		if err != nil && err == io.EOF {
			break
		}
	}

	writer.Flush()

	return tmpFile, fileChanged, nil
}

func (r FileReplacer) confirmMatches(matches []match.MatchString, line string, lineNumber int, stdin *os.File, confirmAnswer *input.ConfirmAnswer) (replacedLine string, lineChanged bool) {
	var answer rune
	var err error

	confirmedQuit := rune(*confirmAnswer) == ConfirmQuit
	confirmedAll := rune(*confirmAnswer) == ConfirmAll

	replacedLine = line

	for i, thisMatch := range matches {
		if confirmedQuit {
			continue
		}

		stringRange := [2]int{0, thisMatch.IndexEnd}
		if i > 0 {
			// Gets the previous one
			stringRange = [2]int{matches[i-1].IndexEnd, thisMatch.IndexEnd}
		}

		if confirmedAll {
			answer = ConfirmYes
		} else {
			answer, err = match.ConfirmMatch(thisMatch, r.inputFilePath, lineNumber, stdin)
			*confirmAnswer = input.ConfirmAnswer(answer)
		}

		if err != nil {
			fmt.Println(err)
		}

		switch answer {
		case ConfirmYes:
			replacedLine = r.LineReplacer.ReplaceStringRange(replacedLine, stringRange)
		case ConfirmNo:
			// Nothing to do
		case ConfirmAll:
			replacedLine = r.LineReplacer.ReplaceStringRange(replacedLine, stringRange)
			confirmedAll = true
		default:
			confirmedQuit = true
		}
	}

	lineChanged = replacedLine != line

	return replacedLine, lineChanged
}

func openInputFile(path string) (*os.File, error) {
	fileStat, _ := os.Lstat(path)
	inputFilePath := path

	if fileStat.Mode().Type() == os.ModeSymlink.Type() {
		inputFilePath, _ = filepath.EvalSymlinks(path)
		inputFilePath, _ = filepath.Abs(inputFilePath)
	}

	return os.OpenFile(inputFilePath, os.O_RDONLY, fileStat.Mode())
}
