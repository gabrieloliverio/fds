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

    var searchArg = flag.Arg(0)
    var replace = flag.Arg(1)
    var subject = readSubject(os.Stdin)

    searchPattern, err := validateArgs(searchArg, replace, subject)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    
    result := findReplace(searchPattern, replace, subject)

    fmt.Println(result)
}

func readSubject(stdin *os.File) string {
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
