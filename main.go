package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gabrieloliverio/fds/glob"
	"github.com/gabrieloliverio/fds/input"
	"github.com/gabrieloliverio/fds/replace"
)

func main() {
	var (
		literalFlag, insensitiveFlag, confirmFlag bool
		err                                       error
	)

	flag.Usage = func() { fmt.Fprint(os.Stderr, input.Usage) }

	flag.BoolVar(&literalFlag, "l", false, "Treat pattern as a regular string instead of as Regular Expression")
	flag.BoolVar(&literalFlag, "literal", false, "Treat pattern as a regular string instead of as Regular Expression")

	flag.BoolVar(&insensitiveFlag, "i", false, "Insensitive case on search")
	flag.BoolVar(&insensitiveFlag, "insensitive", false, "Insensitive case on search")

	flag.BoolVar(&confirmFlag, "c", false, "Confirm each substitution")
	flag.BoolVar(&confirmFlag, "confirm", false, "Confirm each substitution")

	flag.Parse()

	flags := map[string]bool{"confirm": confirmFlag, "insensitive": insensitiveFlag, "literal": literalFlag}

	args, err := input.ReadArgs(os.Stdin, flag.Args())
	check(err)

	err = input.Validate(args, flags)
	check(err)

	if args.File.Path == "" {
		replacer := replace.NewReplacer(flags)

		fmt.Print(replacer.Replace(args.Subject, args.Search, args.Replace))

		return
	}

	if !args.File.IsDir {
		err := replaceInFile(args.File.Path, args, flags, nil)
		check(err)

		return
	}

	filepaths, err := glob.GetFilesInDir(args.File.Path)
	check(err)

	var defaultAnswer = input.ConfirmAnswer('n')
	var confirmAnswer *input.ConfirmAnswer = &defaultAnswer

	for _, file := range filepaths {
		err = replaceInFile(file, args, flags, confirmAnswer)

		check(err)

		if rune(*confirmAnswer) == input.ConfirmQuit {
			break
		}
	}
}

func replaceInFile(path string, args input.Args, flags map[string]bool, confirmAnswer *input.ConfirmAnswer) error {
	inputFile, tmpFile, err := replace.OpenInputAndTempFile(path)

	check(err)

	defer inputFile.Close()
	defer tmpFile.Close()

	replacer := replace.NewFileReplacer(inputFile, tmpFile, flags)

	err = replacer.ReplaceInFile(args.Search, args.Replace, os.Stdin, confirmAnswer)

	check(err)

	err = os.Rename(tmpFile.Name(), inputFile.Name())
	check(err)

	return err
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
