package main

import (
	"fmt"
	"os"

	"github.com/gabrieloliverio/fds"
	"github.com/gabrieloliverio/fds/config"
	"github.com/gabrieloliverio/fds/input"
	"github.com/gabrieloliverio/fds/replace"
	"github.com/spf13/pflag"
)

var (
	literal, insensitive, confirm, verbose, help bool
	ignoreGlobs                            input.IgnoreGlobs
	err                                    error
	defaultAnswer                          = input.ConfirmAnswer('n')
	confirmAnswer                          = &defaultAnswer
)

func main() {
	pflag.Usage = func() { fmt.Fprint(os.Stderr, input.Usage) }

	pflag.BoolVarP(&literal, "literal", "l", false, input.LiteralUsage)
	pflag.BoolVarP(&insensitive, "insensitive", "i", false, input.InsensitiveUsage)
	pflag.BoolVarP(&confirm, "confirm", "c", false, input.ConfirmUsage)
	pflag.BoolVarP(&verbose, "verbose", "v", false, input.VerboseUsage)
	pflag.BoolVarP(&help, "help", "h", false, input.HelpUsage)
	pflag.Var(&ignoreGlobs, "ignore-globs", input.IgnoreUsage)

	pflag.Parse()

	if help {
		fmt.Println(input.Usage)

		os.Exit(0)
	}

	config := config.NewConfig()
	config.Flags = map[string]bool{"confirm": confirm, "insensitive": insensitive, "literal": literal, "verbose": verbose}

	args, err := input.ReadArgs(os.Stdin, pflag.Args())
	fds.CheckError(err)

	err = input.Validate(args, config.Flags)
	fds.CheckError(err)

	if args.Path.Value == "" {
		replacer := replace.NewLineReplacer(args.Search, args.Replace, config.Flags)

		fmt.Print(replacer.Replace(args.Subject))

		return
	}

	if args.Path.IsFile() {
		replacer := replace.NewFileReplacer(args.Path.Value, args.Search, args.Replace, config)

		err = fds.ReplaceInFile(args, replacer, os.Stdin, confirmAnswer)
		fds.CheckError(err)

		return
	}

	files, err := fds.GetFilesInDir(args.Path.Value, ignoreGlobs, verbose)
	fds.CheckError(err)

	err = fds.ReplaceInFiles(files, os.Stdin, args, config, confirmAnswer)
	fds.CheckError(err)
}
