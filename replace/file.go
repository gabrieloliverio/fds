package replace

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gabrieloliverio/fds/input"
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

func openInputFile(path string) (*os.File, error) {
	fileStat, _ := os.Lstat(path)
	inputFilePath := path

	if fileStat.Mode().Type() == os.ModeSymlink.Type() {
		inputFilePath, _ = filepath.EvalSymlinks(path)
		inputFilePath, _ = filepath.Abs(inputFilePath)
	}

	return os.OpenFile(inputFilePath, os.O_RDONLY, fileStat.Mode())
}
