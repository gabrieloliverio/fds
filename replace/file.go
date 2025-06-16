package replace

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gabrieloliverio/fds/config"
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

	config config.Config
	inputFilePath string
}

func NewFileReplacer(inputFilePath, search, replace string, config config.Config) FileReplacer {
	replacer := FileReplacer{
		LineReplacer: LineReplacer{flags: config.Flags, replace: replace},
		inputFilePath: inputFilePath,
		config: config,
	}
	replacer.search = replacer.compilePattern(search)

	return replacer
}

/*
 * ReplaceInFile replaces a given pattern when found in `inputFile`. Lines are written into `outputFile`
 */
func (r FileReplacer) Replace(stdin *os.File, confirmAnswer *input.ConfirmAnswer) (outputFile *os.File, err error) {
	if r.flags["confirm"] {
		outputFile, err = r.confirmAndReplace(stdin, confirmAnswer)

		return
	}

	return r.replaceAll()
}

func (r FileReplacer) replaceAll() (tmpFile *os.File, err error) {
	var fileChanged bool

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

		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("Error while reading file: %s", err)
		}

		replacedLine, lineChanged := r.LineReplacer.Replace(line)

		if lineChanged {
			fileChanged = true
		}

		_, errWrite := writer.WriteString(replacedLine)

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

	return tmpFile, err
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
