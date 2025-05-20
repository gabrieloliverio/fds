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

	inputStat, _ := inputFile.Stat()
	originalModTime := inputStat.ModTime()

	replacer := replace.NewFileReplacer(inputFile, tmpFile, flags)

	if flags["verbose"] {
		log.Printf("Replacing %s for %s in file %s", args.Search, args.Replace, path)
	}

	pattern := replacer.CompilePattern(args.Search)

	err = replacer.ReplaceInFile(pattern, args.Replace, os.Stdin, confirmAnswer)

	CheckError(err)

	inputStat, _ = inputFile.Stat()
	inputFileHasChanged := inputStat.ModTime().After(originalModTime)
	renameFile := true

	if flags["verbose"] {
		log.Printf("Replace in temp file completed")
		log.Printf("Original timestamp of file %s: %s", path, originalModTime)
	}

	if inputFileHasChanged {
		if flags["verbose"] {
			log.Printf("File %s has been modified since %s", path, originalModTime)
		}

		confirmText := fmt.Sprintf("File %s was modified after initial read. Overwrite anyway? [y]es [n]o", path)
		answer, _ := input.Confirm(os.Stdin, confirmText, []rune{'y', 'n'})

		if answer == 'n' {
			renameFile = false

			if flags["verbose"] && inputFileHasChanged {
				log.Printf("File %s will not be overwritten", path)
			}
		}
	}

	if renameFile {
		if flags["verbose"] && inputFileHasChanged {
			log.Printf("Overwriting file %s with contents from temp file", path)
		}

		err = os.Rename(tmpFile.Name(), inputFile.Name())
		CheckError(err)

		if flags["verbose"] {
			log.Printf("Renamed temp file %s to %s", tmpFile.Name(), inputFile.Name())
		}
	}

	return err
}

func ReplaceInFiles(files []string, args input.Args, flags map[string]bool, confirmAnswer *input.ConfirmAnswer) error {
	for _, file := range files {
		err := ReplaceInFile(file, args, flags, confirmAnswer)

		if err != nil {
			return err
		}

		if rune(*confirmAnswer) == replace.ConfirmQuit {
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
