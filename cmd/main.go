package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gabrieloliverio/fds"
	"github.com/gabrieloliverio/fds/input"
	"github.com/gabrieloliverio/fds/replace"
)

func main() {
	var (
		literalFlag, insensitiveFlag, confirmFlag, verboseFlag bool
		ignoreGlobs                                            input.IgnoreGlobs
		err                                                    error
		defaultAnswer                                          = input.ConfirmAnswer('n')
		confirmAnswer                                          = &defaultAnswer
	)

	flag.Usage = func() { fmt.Fprint(os.Stderr, input.Usage) }

	flag.BoolVar(&literalFlag, "l", false, "Treat pattern as a regular string instead of as Regular Expression")
	flag.BoolVar(&literalFlag, "literal", false, "Treat pattern as a regular string instead of as Regular Expression")

	flag.BoolVar(&insensitiveFlag, "i", false, "Insensitive case on search")
	flag.BoolVar(&insensitiveFlag, "insensitive", false, "Insensitive case on search")

	flag.BoolVar(&confirmFlag, "c", false, "Confirm each substitution")
	flag.BoolVar(&confirmFlag, "confirm", false, "Confirm each substitution")

	flag.BoolVar(&verboseFlag, "v", false, "Print debug information")
	flag.BoolVar(&verboseFlag, "verbose", false, "Print debug information")

	flag.Var(&ignoreGlobs, "ignore", "Confirm each substitution")

	flag.Parse()

	flags := map[string]bool{"confirm": confirmFlag, "insensitive": insensitiveFlag, "literal": literalFlag}

	args, err := input.ReadArgs(os.Stdin, flag.Args())
	fds.CheckError(err)

	err = input.Validate(args, flags)
	fds.CheckError(err)

	if args.Path.Value == "" {
		replacer := replace.NewReplacer(flags)

		fmt.Print(replacer.Replace(args.Subject, args.Search, args.Replace))

		return
	}

	if args.Path.IsFile() {
		err := fds.ReplaceInFile(args.Path.Value, args, flags, confirmAnswer)
		fds.CheckError(err)

		return
	}

	files, err := fds.GetFilesInDir(args.Path.Value, ignoreGlobs, verboseFlag)
	fds.CheckError(err)

	err = fds.ReplaceInFiles(files, args, flags, confirmAnswer)
	fds.CheckError(err)
}

