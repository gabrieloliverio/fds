package input

import (
	"io"
	"os"
	"regexp"
	"strings"
)

const Usage = `Usage:
	echo subject | fds [ options ] search_pattern replace
	fds [ options ] search_pattern replace ./file
	fds [ options ] search_pattern replace ~/directory
	fds [ options ] search_pattern replace ~/directory/**/somepattern*

Options:

	-l, --literal 		Treat pattern as a regular string instead of as Regular Expression
	-i, --insensitive 	Ignore case on search
	-c, --confirm 		Confirm each substitution
`

type pathArg struct {
	Value string
	fileInfo os.FileInfo
}

func (p pathArg) IsDir() bool {
	return p.fileInfo.IsDir()
}

func (p pathArg) IsFile() bool {
	return !p.fileInfo.IsDir()
}

type Args struct {
	Subject string
	Search  string
	Replace string

	Path pathArg
}

func Validate(args Args, flags map[string] bool) error {
	_, err := regexp.Compile(args.Search)

	if !flags["literal"] && err != nil {
		return NewInvalidRegExpError()
	}

	if flags["literal"] && flags["insensitive"] {
		return NewLiteralInsensitiveError()
	}

	if flags["confirm"] && args.Path.Value == "" {
		return NewConfirmNotOnFileError()
	}

	if strings.TrimSpace(args.Replace) == "" || strings.TrimSpace(args.Subject) == "" || strings.TrimSpace(args.Search) == "" {
		return NewInvalidArgumentsError()
	}

	return nil
}

func ReadArgs(stdin *os.File, inputArgs []string) (Args, error) {
	stdinStat, _ := stdin.Stat()

	if stdinStat.Size() > 0 {
		stdin, _ := io.ReadAll(stdin)

		if len(inputArgs) < 2 {
			return Args{}, NewInvalidArgumentsError()
		}

		return Args{Subject: string(stdin), Search: inputArgs[0], Replace: inputArgs[1]}, nil
	}

	if len(inputArgs) < 3 {
		return Args{}, NewInvalidArgumentsError()
	}

	args := Args{
		Search:  inputArgs[0],
		Replace: inputArgs[1],
		Subject: inputArgs[2],
	}

	fileStat, err := os.Stat(args.Subject)

	if err != nil {
		return Args{}, NewInvalidArgumentsErrorFileNotFound(args.Subject)
	}

	args.Path = pathArg{Value: args.Subject, fileInfo: fileStat}

	return args, nil
}
