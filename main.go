package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
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

	var inputFile *os.File
	var err error

	args := readArgs()

	if args.file.path == "" {
		err = args.validate()

		check(err)

		result := replacePattern(regexp.MustCompile(args.search), args.replace, args.subject)

		fmt.Println(result)

		os.Exit(0)
	}

	fileStat, _ := os.Stat(args.file.path)
	inputFile, err = os.OpenFile(args.file.path, os.O_RDWR, fileStat.Mode())

	check(err)

	defer inputFile.Close()

	err = replaceInFileChunks(inputFile, args)

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

func replaceInFileChunks(file *os.File, args args) error {
	inputFileStat, _ := file.Stat()
	searchPattern := regexp.MustCompile(args.search)
	var err error

	if inputFileStat.Size() == 0 {
		return nil
	}

	scanner := bufio.NewScanner(file)
	var bs []byte
	buf := bytes.NewBuffer(bs)

	var text string
	for scanner.Scan() {
		text = scanner.Text()

		replaced := replacePattern(searchPattern, args.replace, text)

		_, err = buf.WriteString(replaced + "\n")

		if err != nil {
			return err
		}
	}

	err = file.Truncate(0)

	if err != nil {
		return err
	}

	_, err = file.Seek(0, 0)

	if err != nil {
		return err
	}

	_, err = buf.WriteTo(file)

	return err
}
