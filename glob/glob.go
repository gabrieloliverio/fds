package glob

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func GetFilesInDir(root string) ([]string, error) {
	fileSystem := os.DirFS(root)
	var filepaths []string

	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
		    log.Fatal(err)
		}

        if !d.IsDir() {
            filepaths = append(filepaths, filepath.Join(root, path))
        }

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filepaths, nil
}
