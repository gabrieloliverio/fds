package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

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
	check(err)

	err = input.Validate(args, flags)
	check(err)

	if args.Path.Value == "" {
		replacer := replace.NewReplacer(flags)

		fmt.Print(replacer.Replace(args.Subject, args.Search, args.Replace))

		return
	}

	if args.Path.IsFile() {

		err := replaceInFile(args.Path.Value, args, flags, confirmAnswer)
		check(err)

		return
	}

	files, err := getFilesInDir(args.Path.Value, ignoreGlobs, verboseFlag)
	check(err)

	err = replaceInFiles(files, args, flags, confirmAnswer)
	check(err)
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

func replaceInFiles(files []string, args input.Args, flags map[string]bool, confirmAnswer *input.ConfirmAnswer) error {
	for _, file := range files {
		err := replaceInFile(file, args, flags, confirmAnswer)

		if err != nil {
			return err
		}

		if rune(*confirmAnswer) == input.ConfirmQuit {
			break
		}
	}

	return nil
}

func getFilesInDir(root string, ignoreGlobs input.IgnoreGlobs, verbose bool) ([]string, error) {
	fileSystem := os.DirFS(root)
	var filepaths []string

	if verbose {
		log.Printf("Ignoring glob patterns \"%s\"\n", ignoreGlobs.String())
	}

	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		fullpath := filepath.Join(root, path)

		patternMatch := ignoreGlobs.MatchAny(fullpath)

		if verbose && patternMatch {
			log.Printf("Pattern matched path \"%s\"\n", path)
		}

		if !d.IsDir() && !patternMatch {
			filepaths = append(filepaths, fullpath)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if verbose {
		log.Printf("Found %d files: %s\n", len(filepaths), filepaths)
	}

	return filepaths, nil
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
