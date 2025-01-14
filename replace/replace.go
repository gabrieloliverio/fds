package replace

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/gabrieloliverio/fds/input"
)

func ReplaceStringOrPattern(search, replace, subject string, literalFlag, insensitiveFlag bool) string {
	if literalFlag {
		return strings.ReplaceAll(subject, search, replace)
	}

	if insensitiveFlag {
		search = "(?i)" + search
	}

	return regexp.MustCompile(search).ReplaceAllString(subject, replace)
}

func ReplaceInFile(inputFile, outputFile *os.File, args input.Args, literalFlag, insensitiveFlag bool) error {
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

		replaced := ReplaceStringOrPattern(args.Search, args.Replace, line, literalFlag, insensitiveFlag)

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

