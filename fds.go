package fds

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func ReplaceInFile(args Args, replacer FileReplacer, stdin io.Reader, stdout io.Writer, confirmAnswer *ConfirmAnswer) error {
	var err error
	inputStat, _ := os.Stat(args.Path.Value)
	originalModTime := inputStat.ModTime()

	if replacer.HasFlag("verbose") {
		log.Printf("Replacing %s for %s in file %s", args.Search, args.Replace, args.Path.Value)
	}

	tmpFile, err := replacer.Replace(stdin, stdout, confirmAnswer)

	CheckError(err)

	if tmpFile == nil {
		if replacer.HasFlag("verbose") {
			log.Printf("Nothing replaced in file %s", args.Path.Value)
		}

		return nil
	}

	inputStat, _ = os.Stat(args.Path.Value)
	inputFileChangedSinceRead := inputStat.ModTime().After(originalModTime)
	renameFile := true

	if replacer.HasFlag("verbose") {
		log.Printf("Replace in temp file completed")
		log.Printf("Original timestamp of file %s: %s", args.Path.Value, originalModTime)
	}

	if inputFileChangedSinceRead {
		if replacer.HasFlag("verbose") {
			log.Printf("File %s has been modified since %s", args.Path.Value, originalModTime)
		}

		confirmText := fmt.Sprintf("File %s was modified after initial read. Overwrite anyway? [y]es [n]o", args.Path.Value)
		answer, _ := Confirm(stdin, confirmText, []rune{'y', 'n'})

		if answer == 'n' {
			renameFile = false

			if replacer.HasFlag("verbose") && inputFileChangedSinceRead {
				log.Printf("File %s will not be overwritten", args.Path.Value)
				err = fmt.Errorf("File %s was overritten since it was read. Operation aborted", args.Path.Value)
			}
		}
	}

	if renameFile {
		if replacer.HasFlag("verbose") && inputFileChangedSinceRead {
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

func ReplaceInFiles(files []string, stdin io.Reader, stdout io.Writer, args Args, config Config, confirmAnswer *ConfirmAnswer) error {
	for _, file := range files {
		args.Path.Value = file

		replacer := NewFileReplacer(args.Path.Value, args.Search, args.Replace, config)

		err := ReplaceInFile(args, replacer, stdin, stdout, confirmAnswer)

		if err != nil {
			return err
		}

		if rune(*confirmAnswer) == ConfirmQuit {
			break
		}
	}

	return nil
}

func GetFilesInDir(root string, ignoreGlobs IgnoreGlobs, verbose bool) ([]string, error) {
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
		log.Printf("Found %d files in %s", len(filepaths), root)
	}

	return filepaths, nil
}

func CheckError(err error) {
	if err == nil {
		return
	}

	if thrownErr, ok := err.(Error); ok {
		fmt.Fprintln(os.Stderr, thrownErr.Error())
		os.Exit(thrownErr.Code)
	}

	fmt.Println(err)
	os.Exit(1)
}
