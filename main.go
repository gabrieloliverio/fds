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

func (a args) validate() error {
	usageErr := errors.New("usage:  fds subject search_pattern replace\n" +
		"\techo subject search_pattern replace\n" +
		"\tfds ./file search_pattern replace\n")

	if _, err := regexp.Compile(a.search); err != nil {
		return usageErr
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
	flag.Parse()

	var inputFile, tmpFile *os.File
	var inputFilePath string
	var err error

	args := readArgs()

	if args.file.path == "" {
		err = args.validate()

		check(err)

		result := replacePattern(regexp.MustCompile(args.search), args.replace, args.subject)

		fmt.Println(result)

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

	err = replaceInFile(inputFile, tmpFile, args)

	err = os.Rename(tmpFile.Name(), inputFile.Name())

	check(err)
}

func replacePattern(searchPattern *regexp.Regexp, replace, subject string) string {
	return searchPattern.ReplaceAllString(subject, replace)
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

func replaceInFile(inputFile, outputFile *os.File, args args) error {
	inputFileStat, _ := inputFile.Stat()
	searchPattern := regexp.MustCompile(args.search)
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

		replaced := replacePattern(searchPattern, args.replace, line)

		_, err = writer.WriteString(replaced)

		if err != nil {
			return err
		}
	}

	writer.Flush()

	return err
}
