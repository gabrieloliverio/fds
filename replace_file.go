package fds

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

	config        Config
	inputFilePath string
}

func NewFileReplacer(inputFilePath, search, replace string, config Config) FileReplacer {
	replacer := FileReplacer{
		LineReplacer:  LineReplacer{flags: config.Flags, replace: replace, search: search},
		inputFilePath: inputFilePath,
		config:        config,
	}
	replacer.searchRegexp = replacer.compilePattern(search)

	return replacer
}

func (r FileReplacer) Replace(stdin io.Reader, stdout io.Writer, confirmAnswer *ConfirmAnswer) (outputFile *os.File, err error) {
	if r.flags["confirm"] {
		outputFile, err = r.replaceInteractive(stdin, stdout, confirmAnswer)

		return
	}

	return r.replaceAll()
}

func (r FileReplacer) replaceAll() (tmpFile *os.File, err error) {
	var fileChanged bool

	inputFile, err := openInputFile(r.inputFilePath)

	if err != nil {
		return nil, NewFileReadError(r.inputFilePath)
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
			return nil, NewFileReadError(r.inputFilePath)
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
		tmpFile, err = os.CreateTemp("", filepath.Base(inputFile.Name()))

		if err != nil {
			return nil, NewTempFileWriteError(filepath.Base(inputFile.Name()))
		}

		_, err = io.Copy(tmpFile, buffer)

		if err != nil {
			return nil, NewFileWriteError(inputFile.Name())
		}
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
