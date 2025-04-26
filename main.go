package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gabrieloliverio/fds/input"
	"github.com/gabrieloliverio/fds/replace"
)

func main() {
	var (
		literalFlag, insensitiveFlag, confirmFlag bool
		inputFile, tmpFile                        *os.File
		inputFilePath                             string
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

	inputFile, err = getInputFile(args)
	check(err)

	defer inputFile.Close()

	tmpFile, err = os.CreateTemp("", filepath.Base(inputFilePath))
	check(err)

	defer tmpFile.Close()

	replacer := replace.NewFileReplacer(inputFile, tmpFile, flags)

	err = replacer.ReplaceInFile(args.Search, args.Replace, os.Stdin)

	check(err)

	err = os.Rename(tmpFile.Name(), inputFile.Name())
	check(err)
}

func getInputFile(args input.Args) (*os.File, error) {
	fileStat, _ := os.Lstat(args.File.Path)
	inputFilePath := args.File.Path

	if fileStat.Mode().Type() == os.ModeSymlink.Type() {
		inputFilePath, _ = filepath.EvalSymlinks(args.File.Path)
		inputFilePath, _ = filepath.Abs(inputFilePath)
	}

	return os.OpenFile(inputFilePath, os.O_RDONLY, fileStat.Mode())
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
