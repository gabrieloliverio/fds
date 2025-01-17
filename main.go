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
		literalFlag, insensitiveFlag bool
	)

	flag.Usage = func() { fmt.Fprint(os.Stderr, input.Usage) }

	flag.BoolVar(&literalFlag, "l", false, "Treat pattern as a regular string instead of as Regular Expression")
	flag.BoolVar(&literalFlag, "literal", false, "Treat pattern as a regular string instead of as Regular Expression")

	flag.BoolVar(&insensitiveFlag, "i", false, "Insensitive case on search")
	flag.BoolVar(&insensitiveFlag, "insensitive", false, "Insensitive case on search")

	flag.Parse()

	var inputFile, tmpFile *os.File
	var inputFilePath string
	var err error

	args := input.ReadArgs(os.Stdin, flag.Args())

	err = input.Validate(args, literalFlag, insensitiveFlag)
	check(err)

	if args.File.Path == "" {
		fmt.Println(replace.ReplaceStringOrPattern(args.Search, args.Replace, args.Subject, literalFlag, insensitiveFlag))

		os.Exit(0)
	}

	inputFile, err = getInputFile(args)
	check(err)

	defer inputFile.Close()

	tmpFile, err = os.CreateTemp("", filepath.Base(inputFilePath))
	check(err)

	defer tmpFile.Close()

	err = replace.ReplaceInFile(inputFile, tmpFile, args, literalFlag, insensitiveFlag)
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
