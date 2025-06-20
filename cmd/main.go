package main

import (
	"fmt"
	"os"

	"github.com/gabrieloliverio/fds"
	"github.com/spf13/pflag"
)

var (
	literal, insensitive, confirm, verbose, help bool
	ignoreGlobs                            fds.IgnoreGlobs
	err                                    error
	defaultAnswer                          = fds.ConfirmAnswer('n')
	confirmAnswer                          = &defaultAnswer
)

func main() {
	pflag.Usage = func() { fmt.Fprint(os.Stderr, fds.Usage) }

	pflag.BoolVarP(&literal, "literal", "l", false, fds.LiteralUsage)
	pflag.BoolVarP(&insensitive, "insensitive", "i", false, fds.InsensitiveUsage)
	pflag.BoolVarP(&confirm, "confirm", "c", false, fds.ConfirmUsage)
	pflag.BoolVarP(&verbose, "verbose", "v", false, fds.VerboseUsage)
	pflag.BoolVarP(&help, "help", "h", false, fds.HelpUsage)
	pflag.Var(&ignoreGlobs, "ignore-globs", fds.IgnoreUsage)

	pflag.Parse()

	if help {
		fmt.Println(fds.Usage)

		os.Exit(0)
	}

	config := fds.NewConfig()
	config.Flags = map[string]bool{"confirm": confirm, "insensitive": insensitive, "literal": literal, "verbose": verbose}

	args, err := fds.ReadArgs(os.Stdin, pflag.Args())
	fds.CheckError(err)

	err = fds.Validate(args, config.Flags)
	fds.CheckError(err)

	if args.Path.Value == "" {
		replacer := fds.NewLineReplacer(args.Search, args.Replace, config.Flags)

		fmt.Print(replacer.Replace(args.Subject))

		return
	}

	if args.Path.IsFile() {
		replacer := fds.NewFileReplacer(args.Path.Value, args.Search, args.Replace, config)

		err = fds.ReplaceInFile(args, replacer, os.Stdin, confirmAnswer)
		fds.CheckError(err)

		return
	}

	files, err := fds.GetFilesInDir(args.Path.Value, ignoreGlobs, verbose)
	fds.CheckError(err)

	err = fds.ReplaceInFiles(files, os.Stdin, args, config, confirmAnswer)
	fds.CheckError(err)
}
