package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type fileArg struct {
	path string
}

type args struct {
	subject string
	search  string
	replace string

	file fileArg
}

const usage = `Usage:
	fds [--literal] subject [ search_pattern ] [ replace ]
	echo subject | fds [ options ] [ search_pattern ] [ replace ]
	fds [--literal] ./file [ search_pattern ] [ replace ]

Options:

	-l, --literal 		Treat pattern as a regular string instead of as Regular Expression
	-i, --insensitive 	Ignore case on search
`

func validate(a args, literalMode, insensitiveMode bool) error {
	usageErr := errors.New(usage)

	_, err := regexp.Compile(a.search)

	if !literalMode && err != nil {
		return fmt.Errorf("%s\n%s\n", usage, "[ search_pattern ] is not a valid Regular Expression")
	}

	if literalMode && insensitiveMode {
		return fmt.Errorf("%s\n%s\n", usage, "[ -l, --literal ] cannot be used along with [ -i, --insensitive ]")
	}

	if a.replace == "" || a.subject == "" {
		return usageErr
	}

	return nil
}

func readArgs() args {
	if subject := readFile(os.Stdin); subject != "" {
		return args{subject: subject, search: flag.Arg(0), replace: flag.Arg(1)}
	}

	args := args{
		subject: flag.Arg(0),
		search:  flag.Arg(1),
		replace: flag.Arg(2),
	}

	fileStat, err := os.Stat(args.subject)

	if err == nil && !fileStat.IsDir() {
		args.file = fileArg{path: args.subject}
	}

	return args
}

func main() {
	var (
		literalFlag, insensitiveFlag bool
	)

	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }

	flag.BoolVar(&literalFlag, "l", false, "Treat pattern as a regular string instead of as Regular Expression")
	flag.BoolVar(&literalFlag, "literal", false, "Treat pattern as a regular string instead of as Regular Expression")

	flag.BoolVar(&insensitiveFlag, "i", false, "Insensitive case on search")
	flag.BoolVar(&insensitiveFlag, "insensitive", false, "Insensitive case on search")

	flag.Parse()

	var inputFile, tmpFile *os.File
	var inputFilePath string
	var err error

	args := readArgs()

	if args.file.path == "" {
		err = validate(args, literalFlag, insensitiveFlag)

		check(err)

		fmt.Println(replaceStringOrPattern(args.search, args.replace, args.subject, literalFlag, insensitiveFlag))

		os.Exit(0)
	}

	inputFilePath = args.file.path
	fileStat, _ := os.Lstat(args.file.path)

	if fileStat.Mode().Type() == os.ModeSymlink.Type() {
		inputFilePath, _ = filepath.EvalSymlinks(args.file.path)
		inputFilePath, _ = filepath.Abs(inputFilePath)
	}

	inputFile, err = os.OpenFile(inputFilePath, os.O_RDONLY, fileStat.Mode())

	check(err)

	defer inputFile.Close()

	tmpFile, err = os.CreateTemp("", filepath.Base(inputFilePath))

	check(err)

	defer tmpFile.Close()

	err = replaceInFile(inputFile, tmpFile, args, literalFlag, insensitiveFlag)

	err = os.Rename(tmpFile.Name(), inputFile.Name())

	check(err)
}

func replaceStringOrPattern(search, replace, subject string, literalFlag, insensitiveFlag bool) string {
	if literalFlag {
		return strings.ReplaceAll(subject, search, replace)
	}

	if insensitiveFlag {
		search = "(?i)" + search
	}

	return regexp.MustCompile(search).ReplaceAllString(subject, replace)
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readFile(file *os.File) string {
	stat, _ := file.Stat()

	if stat.Size() == 0 {
		return ""
	}

	scanner := bufio.NewScanner(file)
	var content string

	for scanner.Scan() {
		content += scanner.Text()
	}

	return content
}

func replaceInFile(inputFile, outputFile *os.File, args args, literalFlag, insensitiveFlag bool) error {
	inputFileStat, _ := inputFile.Stat()
	var err error

	if inputFileStat.Size() == 0 {
		return nil
	}

	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)

	for {
		line, err := reader.ReadString('\n')

		if err != nil && err == io.EOF {
			break
		}

		replaced := replaceStringOrPattern(args.search, args.replace, line, literalFlag, insensitiveFlag)

		_, err = writer.WriteString(replaced)

		if err != nil {
			return err
		}
	}

	writer.Flush()

	return err
}
