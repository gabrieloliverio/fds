package fds

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (r FileReplacer) confirmAndReplace(stdin *os.File, confirmAnswer *ConfirmAnswer) (outputFile *os.File, err error) {
	var (
		lineNumber                  int
		tmpFile                     *os.File
		confirmedAll, confirmedQuit bool
		lineChanged                 bool
		fileChanged                 bool
	)

	// feedback confirmedAll that propagated all the way back to the main function
	confirmedAll = *confirmAnswer == ConfirmAll

	inputFile, err := openInputFile(r.inputFilePath)

	if err != nil {
		return
	}

	if inputFileStat, _ := inputFile.Stat(); inputFileStat.Size() == 0 {
		return
	}

	reader := bufio.NewReader(inputFile)

	buffer := &bytes.Buffer{}
	writer := bufio.NewWriter(buffer)

	for {
		line, err := reader.ReadString('\n')
		lineNumber++

		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("Error while reading file: %s", err)
		}

		if confirmedAll {
			line, lineChanged = r.LineReplacer.Replace(line)
		}

		if !confirmedAll && !confirmedQuit {
			matches := FindStringOrPattern(r.search, r.replace, line, 50)

			line, lineChanged = r.confirmMatches(matches, line, lineNumber, stdin, confirmAnswer)
		}

		if lineChanged {
			fileChanged = true
		}

		_, errWrite := writer.WriteString(line)

		if errWrite != nil {
			return nil, fmt.Errorf("Error while writing temporary file: %s", err)
		}

		if err != nil && err == io.EOF {
			break
		}
	}

	writer.Flush()

	if fileChanged {
		tmpFile, _ = os.CreateTemp("", filepath.Base(inputFile.Name()))
		io.Copy(tmpFile, buffer)
	}

	return tmpFile, nil
}

func (r FileReplacer) confirmMatches(matches []MatchString, line string, lineNumber int, stdin *os.File, confirmAnswer *ConfirmAnswer) (replacedLine string, lineChanged bool) {
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
			answer, err = ConfirmMatch(thisMatch, r.inputFilePath, lineNumber, stdin)
			*confirmAnswer = ConfirmAnswer(answer)
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
