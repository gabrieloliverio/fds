package fds

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

const (
	LiteralUsage     = "Treat pattern as a regular string instead of as Regular Expression"
	ConfirmUsage     = "Confirm each substitution"
	InsensitiveUsage = "Ignore case on search"
	VerboseUsage     = "Print debug information"
	IgnoreUsage      = "Ignore glob patterns, comma-separated. Ex. --ignore-globs \"vendor/**,node_modules/lib/**.js\""
	HelpUsage        = "Print out help"
)

var Usage = fmt.Sprintf(`fds is modern and opinionated find/replace CLI program

Usage:
	echo subject | fds [ options ] search_pattern replace
	fds [ options ] search_pattern replace ./file
	fds [ options ] search_pattern replace ~/directory
	fds [ options ] search_pattern replace ~/directory/**/somepattern*

Options:

	-l, --literal        %s
	-i, --insensitive    %s
	-c, --confirm        %s
	-v, --verbose        %s
	--ignore-globs       %s
	-h, --help           %s
`, LiteralUsage, InsensitiveUsage, ConfirmUsage, VerboseUsage, IgnoreUsage, HelpUsage)

type PathArg struct {
	Value    string
	fileInfo os.FileInfo
}

func (p PathArg) IsDir() bool {
	return p.fileInfo.IsDir()
}

func (p PathArg) IsFile() bool {
	return !p.fileInfo.IsDir()
}

type Args struct {
	Subject string
	Search  string
	Replace string

	Path PathArg
}

func Validate(args Args, flags map[string]bool) error {
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

func ReadArgs(stdin io.Reader, inputArgs []string) (Args, error) {
	stdInput, err := io.ReadAll(stdin)

	if len(stdInput) > 0 {
		if len(inputArgs) < 2 {
			return Args{}, NewInvalidArgumentsError()
		}

		return Args{Subject: string(stdInput), Search: inputArgs[0], Replace: inputArgs[1]}, nil
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

	args.Path = PathArg{Value: args.Subject, fileInfo: fileStat}

	return args, nil
}

type IgnoreGlobs []string

func (i *IgnoreGlobs) String() string {
	return strings.Join(*i, ",")
}

func (i *IgnoreGlobs) Type() string {
	return "stringSlice"
}

func (i *IgnoreGlobs) Get() []string {
	return []string(*i)
}

func (i *IgnoreGlobs) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i IgnoreGlobs) MatchAny(filePath string) bool {
	for _, pattern := range []string(i) {
		matches, _ := doublestar.PathMatch(pattern, filePath)

		if matches {
			return true
		}
	}

	return false
}
