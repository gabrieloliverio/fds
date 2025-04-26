package main

import (
	"bufio"
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
	usage := "usage:  fds subject search_pattern replace\n" +
		"\techo subject search_pattern replace\n" +
		"\tfds ./file search_pattern replace\n"

	if _, err := regexp.Compile(a.search); err != nil {
		return fmt.Errorf(usage)
	}

	if a.replace == "" {
		return fmt.Errorf(usage)
	}

	if a.subject == "" {
		return fmt.Errorf(usage)
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

	var inputFile, outputFile *os.File
	var err error

	args := readArgs()

	if args.file.path == "" {
		err = args.validate()

		check(err)

		result := replacePattern(regexp.MustCompile(args.search), args.replace, args.subject)

		fmt.Println(result)

		os.Exit(0)
	}

	inputFile, err = os.Open(args.file.path)
	check(err)

	defer inputFile.Close()

	fileStat, _ := inputFile.Stat()
	subject := readFile(inputFile)

	outputFile, err = os.OpenFile(args.file.path, os.O_WRONLY|os.O_TRUNC, fileStat.Mode().Perm())
	check(err)

	defer outputFile.Close()

	err = replaceInFile(subject, outputFile, args)

	check(err)
}

func replacePattern(searchPattern *regexp.Regexp, replace, subject string) string {
	return searchPattern.ReplaceAllString(subject, replace)
}

func replaceInFile(fileContent string, outputFile *os.File, args args) error {
	result := replacePattern(regexp.MustCompile(args.search), args.replace, fileContent)

	_, err := outputFile.WriteString(result)

	return err
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readFile(stdin *os.File) string {
	stat, _ := stdin.Stat()

	if stat.Size() == 0 {
		return ""
	}

	scanner := bufio.NewScanner(stdin)
	subject := ""

	for scanner.Scan() {
		subject += scanner.Text()
	}

	return subject
}
