package input

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type fileArg struct {
	Path string
}

type Args struct {
	Subject string
	Search  string
	Replace string

	File fileArg
}

func Validate(a Args, usage string, literalMode, insensitiveMode bool) error {
	usageErr := errors.New(usage)

	_, err := regexp.Compile(a.Search)

	if !literalMode && err != nil {
		return fmt.Errorf("%s\n%s\n", usage, "[ search_pattern ] is not a valid Regular Expression")
	}

	if literalMode && insensitiveMode {
		return fmt.Errorf("%s\n%s\n", usage, "[ -l, --literal ] cannot be used along with [ -i, --insensitive ]")
	}

	if strings.Trim(a.Replace, " ") == "" || strings.Trim(a.Subject, " ") == "" || strings.Trim(a.Search, " ") == "" {
		return usageErr
	}

	return nil
}

func ReadArgs() Args {
	if subject := readFile(os.Stdin); subject != "" {
		return Args{Subject: subject, Search: flag.Arg(0), Replace: flag.Arg(1)}
	}

	args := Args{
		Subject: flag.Arg(0),
		Search:  flag.Arg(1),
		Replace: flag.Arg(2),
	}

	fileStat, err := os.Stat(args.Subject)

	if err == nil && !fileStat.IsDir() {
		args.File = fileArg{Path: args.Subject}
	}

	return args
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
