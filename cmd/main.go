package main

import (
	"fmt"
	"io"
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

	config := fds.NewConfig()
	config.Flags = map[string]bool{"confirm": confirm, "insensitive": insensitive, "literal": literal, "verbose": verbose}

	if err := execute(os.Args[1:], config, os.Stdin, os.Stdout); err != nil {
		if thrownErr, ok := err.(fds.Error); ok {
			fmt.Fprintln(os.Stderr, thrownErr.Error())
			os.Exit(thrownErr.Code)
		}

		fmt.Println(err)
		os.Exit(1)
	}
}

func execute(inputArgs []string, config fds.Config, stdin io.Reader, stdout io.Writer) (err error) {
	if config.Flags["help"] {
		fmt.Fprint(stdout, fds.Usage)

		return
	}

	args, err := fds.ReadArgs(stdin, inputArgs)

	if err != nil {
		return
	}

	err = fds.Validate(args, config.Flags)

	if err != nil {
		return
	}

	if args.Path.Value == "" {
		replacer := fds.NewLineReplacer(args.Search, args.Replace, config.Flags)
		result, _ := replacer.Replace(args.Subject)

		fmt.Fprint(stdout, result)

		return
	}

	if args.Path.IsFile() {
		replacer := fds.NewFileReplacer(args.Path.Value, args.Search, args.Replace, config)

		err = fds.ReplaceInFile(args, replacer, stdin, stdout, confirmAnswer)

		return
	}

	files, err := fds.GetFilesInDir(args.Path.Value, ignoreGlobs, verbose)

	if err != nil {
		return
	}

	err = fds.ReplaceInFiles(files, stdin, stdout, args, config, confirmAnswer)

	return
}
