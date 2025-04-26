package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"regexp"
)

func main() {
	flag.Parse()

	var subjectArg, searchArg, replace, content string
	var subjectIsFile bool
	var inputFile, outputFile *os.File
	var err error
	var fileStat fs.FileInfo

	content = readFile(os.Stdin)

	if content != "" {
		searchArg = flag.Arg(0)
		replace = flag.Arg(1)
	} else {
		subjectArg = flag.Arg(0)
		searchArg = flag.Arg(1)
		replace = flag.Arg(2)

		fileStat, err = os.Stat(subjectArg)

		if err == nil && !fileStat.IsDir() {
			subjectIsFile = true

			inputFile, err = os.OpenFile(subjectArg, os.O_RDONLY, fileStat.Mode().Perm())
			check(err)

			defer inputFile.Close()

			content = readFile(inputFile)

			outputFile, err = os.OpenFile(subjectArg, os.O_WRONLY|os.O_TRUNC, fileStat.Mode().Perm())
			check(err)

			defer outputFile.Close()
		} else {
			content = subjectArg
		}
	}

	searchPattern, err := validateArgs(searchArg, replace, content)

	check(err)

	result := findReplace(searchPattern, replace, content)

	if subjectIsFile {
		_, err = outputFile.WriteString(result)

		check(err)
	} else {
		fmt.Println(result)
	}
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

func validateArgs(search, replace, subject string) (*regexp.Regexp, error) {
	searchPattern, err := regexp.Compile(search)
	if err != nil {
		err = fmt.Errorf("search is not a valid pattern")
	}

	if replace == "" {
		err = fmt.Errorf("replace is empty")
	}

	if subject == "" {
		err = fmt.Errorf("subject is empty")
	}

	return searchPattern, err
}

func findReplace(searchPattern *regexp.Regexp, replace, subject string) string {
	return string(searchPattern.ReplaceAllString(subject, replace))
}
