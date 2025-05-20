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
		literal, insensitive, confirm, verbose    bool
		ignoreGlobs                               input.IgnoreGlobs
		err                                       error
		defaultAnswer                             = input.ConfirmAnswer('n')
		confirmAnswer                             = &defaultAnswer
	)

	flag.Usage = func() { fmt.Fprint(os.Stderr, input.Usage) }

	flag.BoolVar(&literal, "l", false, input.LiteralUsage)
	flag.BoolVar(&literal, "literal", false, input.LiteralUsage)

	flag.BoolVar(&insensitive, "i", false, input.InsensitiveUsage)
	flag.BoolVar(&insensitive, "insensitive", false, input.InsensitiveUsage)

	flag.BoolVar(&confirm, "c", false, input.ConfirmUsage)
	flag.BoolVar(&confirm, "confirm", false, input.ConfirmUsage)

	flag.BoolVar(&verbose, "v", false, input.VerboseUsage)
	flag.BoolVar(&verbose, "verbose", false, input.VerboseUsage)

	flag.Var(&ignoreGlobs, "ignore-globs", input.IgnoreUsage)

	flag.Parse()

	flags := map[string]bool{"confirm": confirm, "insensitive": insensitive, "literal": literal, "verbose": verbose}

	args, err := input.ReadArgs(os.Stdin, flag.Args())
	fds.CheckError(err)

	err = input.Validate(args, flags)
	fds.CheckError(err)

	if args.Path.Value == "" {
		replacer := replace.NewReplacer(flags)
		pattern := replacer.CompilePattern(args.Search)

		fmt.Print(replacer.Replace(pattern, args.Subject, args.Replace))

		return
	}

	if args.Path.IsFile() {
		err := fds.ReplaceInFile(args.Path.Value, args, flags, confirmAnswer)
		fds.CheckError(err)

		return
	}

	files, err := fds.GetFilesInDir(args.Path.Value, ignoreGlobs, verbose)
	fds.CheckError(err)

	err = fds.ReplaceInFiles(files, args, flags, confirmAnswer)
	fds.CheckError(err)
}

