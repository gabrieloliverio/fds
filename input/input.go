package input

import (
	"io"
	"os"
	"regexp"
	"strings"
)

const Usage = `Usage:
	fds [ options ] subject search_pattern  replace
	echo subject | fds [ options ] search_pattern replace
	fds [ options ] ./file search_pattern replace

Options:

	-l, --literal 		Treat pattern as a regular string instead of as Regular Expression
	-i, --insensitive 	Ignore case on search
`

type fileArg struct {
	Path string
}

type Args struct {
	Subject string
	Search  string
	Replace string

	File fileArg
}

func Validate(a Args, literalMode, insensitiveMode bool) error {
	_, err := regexp.Compile(a.Search)

	if !literalMode && err != nil {
		return NewInvalidRegExpError()
	}

	if literalMode && insensitiveMode {
		return NewLiteralInsensitiveError()
	}

	if strings.Trim(a.Replace, " ") == "" || strings.Trim(a.Subject, " ") == "" || strings.Trim(a.Search, " ") == "" {
		return NewInvalidArgumentsError()
	}

	return nil
}

func ReadArgs(stdin *os.File, flagArgs []string) Args {
	stdinStat, _ := stdin.Stat()

	if stdinStat.Size() > 0 {
		stdin, _ := io.ReadAll(stdin)

		return Args{Subject: string(stdin), Search: flagArgs[0], Replace: flagArgs[1]}
	}

	args := Args{
		Subject: flagArgs[0],
		Search:  flagArgs[1],
		Replace: flagArgs[2],
	}

	fileStat, err := os.Stat(args.Subject)

	if err == nil && !fileStat.IsDir() {
		args.File = fileArg{Path: args.Subject}
	}

	return args
}
