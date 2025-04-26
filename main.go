package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
)

func main() {
	flag.Parse()

	var subject, searchArg, replace string
	subject = readSubject(os.Stdin)

	if subject != "" {
		searchArg = flag.Arg(0)
		replace = flag.Arg(1)
	} else {
		subject = flag.Arg(0)
		searchArg = flag.Arg(1)
		replace = flag.Arg(2)
	}

	searchPattern, err := validateArgs(searchArg, replace, subject)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	result := findReplace(searchPattern, replace, subject)

	fmt.Println(result)
}

func readSubject(stdin *os.File) string {
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
	return string(searchPattern.ReplaceAll([]byte(subject), []byte(replace)))
}
