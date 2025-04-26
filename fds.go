package fds

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/gabrieloliverio/fds/input"
	"github.com/gabrieloliverio/fds/replace"
)

func ReplaceInFile(path string, args input.Args, flags map[string]bool, confirmAnswer *input.ConfirmAnswer) error {
	inputFile, tmpFile, err := replace.OpenInputAndTempFile(path)

	CheckError(err)

	defer inputFile.Close()
	defer tmpFile.Close()

	replacer := replace.NewFileReplacer(inputFile, tmpFile, flags)

	err = replacer.ReplaceInFile(args.Search, args.Replace, os.Stdin, confirmAnswer)

	CheckError(err)

	err = os.Rename(tmpFile.Name(), inputFile.Name())
	CheckError(err)

	return err
}

func ReplaceInFiles(files []string, args input.Args, flags map[string]bool, confirmAnswer *input.ConfirmAnswer) error {
	for _, file := range files {
		err := ReplaceInFile(file, args, flags, confirmAnswer)

		if err != nil {
			return err
		}

		if rune(*confirmAnswer) == input.ConfirmQuit {
			break
		}
	}

	return nil
}

func GetFilesInDir(root string, ignoreGlobs input.IgnoreGlobs, verbose bool) ([]string, error) {
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

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
