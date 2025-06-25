package fds

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func ReplaceInFile(replacer FileReplacer, stdin io.Reader, stdout io.Writer, confirmAnswer *ConfirmAnswer) error {
	var err error

	file := replacer.inputFilePath
	search := replacer.search
	replace := replacer.replace

	inputStat, _ := os.Stat(file)
	originalModTime := inputStat.ModTime()

	if replacer.HasFlag("verbose") {
		log.Printf("Replacing %s for %s in file %s", search, replace, file)
	}

	tmpFile, err := replacer.Replace(stdin, stdout, confirmAnswer)

	if err != nil {
		return err
	}

	if tmpFile == nil {
		if replacer.HasFlag("verbose") {
			log.Printf("Nothing replaced in file %s", file)
		}

		return nil
	}

	inputStat, _ = os.Stat(file)
	inputFileChangedSinceRead := inputStat.ModTime().After(originalModTime)
	renameFile := true

	if replacer.HasFlag("verbose") {
		log.Printf("Replace in temp file completed")
		log.Printf("Original timestamp of file %s: %s", file, originalModTime)
	}

	if inputFileChangedSinceRead {
		if replacer.HasFlag("verbose") {
			log.Printf("File %s has been modified since %s", file, originalModTime)
		}

		confirmText := fmt.Sprintf("File %s was modified after initial read. Overwrite anyway? [y]es [n]o", file)
		answer, _ := Confirm(stdin, confirmText, []rune{'y', 'n'})

		if answer == 'n' {
			renameFile = false

			if replacer.HasFlag("verbose") && inputFileChangedSinceRead {
				log.Printf("File %s will not be overwritten", file)

				return NewAbortedOperationError()
			}
		}
	}

	if renameFile {
		if replacer.HasFlag("verbose") && inputFileChangedSinceRead {
			log.Printf("Overwriting file %s with contents from temp file", file)
		}

		err = os.Rename(tmpFile.Name(), file)

		if err != nil {
			return NewRenameFileError(file)
		}

		if replacer.HasFlag("verbose") {
			log.Printf("Renamed temp file %s to %s", tmpFile.Name(), file)
		}
	}

	return err
}

func worker(id int, args Args, wg *sync.WaitGroup, stdin io.Reader, stdout io.Writer, config Config, jobs <-chan string, errors chan<- string) {
	defer wg.Done()

	if config.Flags["verbose"] {
		log.Printf("Worker %d initialized", id)
	}

	for file := range jobs {
		replacer := NewFileReplacer(file, args.Search, args.Replace, config)

		err := ReplaceInFile(replacer, stdin, stdout, nil)

		if err != nil {
			errors <- err.Error()
		}
	}
}

func replaceInFilesConcurrently(files []string, stdin io.Reader, stdout io.Writer, args Args, config Config) error {
	jobs := make(chan string, len(files))
	errors := make(chan string)
	var wg sync.WaitGroup

	for i := range config.Workers {
		wg.Add(1)

		go worker(i, args, &wg, stdin, stdout, config, jobs, errors)
	}

	for _, file := range files {
		jobs <- file
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(errors)
	}()

	for error := range errors {
		fmt.Println(error)
	}

	return nil
}

func ReplaceInFiles(files []string, stdin io.Reader, stdout io.Writer, args Args, config Config, confirmAnswer *ConfirmAnswer) error {
	useWorkers := config.Flags["verbose"] && len(files) > 1

	if !config.Flags["confirm"] {
		if useWorkers {
			log.Printf("Number of workers set for operation: %d", config.Workers)
		}

		return replaceInFilesConcurrently(files, stdin, stdout, args, config)
	}

	if config.Flags["verbose"] {
		log.Printf("Find/replace won't be performed concurrently as flag confirm was supplied")
	}

	for _, file := range files {
		replacer := NewFileReplacer(file, args.Search, args.Replace, config)

		err := ReplaceInFile(replacer, stdin, stdout, confirmAnswer)

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
			return NewDirectoryReadError(d.Name())
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
