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

func ReplaceInFile(args input.Args, replacer replace.FileReplacer, stdin *os.File, confirmAnswer *input.ConfirmAnswer) error {
	inputStat, _ := os.Stat(args.Path.Value)
	originalModTime := inputStat.ModTime()

	if replacer.HasFlag("verbose") {
		log.Printf("Replacing %s for %s in file %s", args.Search, args.Replace, args.Path.Value)
	}

	tmpFile, fileChanged, err := replacer.Replace(os.Stdin, confirmAnswer)

	CheckError(err)

	if !fileChanged {
		if replacer.HasFlag("verbose") {
			log.Printf("Nothing replaced in file %s. Removing temp file", args.Path.Value)
		}

		os.Remove(tmpFile.Name())

		return nil
	}

	inputStat, _ = os.Stat(args.Path.Value)
	inputFileHasChanged := inputStat.ModTime().After(originalModTime)
	renameFile := true

	if replacer.HasFlag("verbose") {
		log.Printf("Replace in temp file completed")
		log.Printf("Original timestamp of file %s: %s", args.Path.Value, originalModTime)
	}

	if inputFileHasChanged {
		if replacer.HasFlag("verbose") {
			log.Printf("File %s has been modified since %s", args.Path.Value, originalModTime)
		}

		confirmText := fmt.Sprintf("File %s was modified after initial read. Overwrite anyway? [y]es [n]o", args.Path.Value)
		answer, _ := input.Confirm(stdin, confirmText, []rune{'y', 'n'})

		if answer == 'n' {
			renameFile = false

			if replacer.HasFlag("verbose") && inputFileHasChanged {
				log.Printf("File %s will not be overwritten", args.Path.Value)
			}
		}
	}

	if renameFile {
		if replacer.HasFlag("verbose") && inputFileHasChanged {
			log.Printf("Overwriting file %s with contents from temp file", args.Path.Value)
		}

		err = os.Rename(tmpFile.Name(), args.Path.Value)
		CheckError(err)

		if replacer.HasFlag("verbose") {
			log.Printf("Renamed temp file %s to %s", tmpFile.Name(), args.Path.Value)
		}
	}

	return err
}

func ReplaceInFiles(files []string, stdin *os.File, args input.Args, flags map[string]bool, confirmAnswer *input.ConfirmAnswer) error {
	for _, file := range files {
		args.Path.Value = file

		replacer := replace.NewFileReplacer(args.Path.Value, args.Search, args.Replace, flags)

		err := ReplaceInFile(args, replacer, stdin, confirmAnswer)

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
