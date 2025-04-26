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
func (r FileReplacer) ReplaceInFile(search, replace string, stdin *os.File, confirmAnswer *input.ConfirmAnswer) error {
	if r.flags["confirm"] {
		err := r.confirmAndReplace(search, replace, stdin, confirmAnswer)

		return err
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

func (r FileReplacer) confirmAndReplace(search, replace string, stdin *os.File, confirmAnswer *input.ConfirmAnswer) error {
	var replacer = NewReplacer(r.flags)
	var lineNumber int
	var confirmedAll, confirmedQuit bool

	// feedback confirmedAll that propagated all the way back to the main function
	confirmedAll = *confirmAnswer == input.ConfirmAll 

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
			line = replacer.Replace(line, search, replace)
		}

		if !confirmedAll && !confirmedQuit  {
			matches := match.FindStringOrPattern(search, replace, line, r.flags, 50)

			line = r.confirmMatches(matches, line, search, replace, lineNumber, stdin, confirmAnswer)
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

func (r FileReplacer) confirmMatches(matches []match.MatchString, line, search, replace string, lineNumber int, stdin *os.File, confirmAnswer *input.ConfirmAnswer) string {
	var answer rune
	var err error

	confirmedQuit := rune(*confirmAnswer) == input.ConfirmQuit
	confirmedAll := rune(*confirmAnswer) == input.ConfirmAll
	
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
			answer = input.ConfirmYes
		} else {
			answer, err = match.ConfirmMatch(thisMatch, r.inputFile.Name(), lineNumber, stdin)
			*confirmAnswer = input.ConfirmAnswer(answer)
		}

		if err != nil {
			fmt.Println(err)
		}

		switch answer {
		case input.ConfirmYes:
			replacedLine = r.LineReplacer.ReplaceStringRange(replacedLine, search, replace, stringRange)
		case input.ConfirmNo:
			// Nothing to do
		case input.ConfirmAll:
			replacedLine = r.LineReplacer.ReplaceStringRange(replacedLine, search, replace, stringRange)
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
