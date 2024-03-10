package app

import (
	"fmt"
	"os"
	"path/filepath"
)

func Run(dir string, debug bool) error {

	directory := dir
	f := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".md" {
			if debug {
				fmt.Printf("Validating %s \n", path)
			}

			if err := ValidateLinks(path); err != nil {
				fmt.Printf("Error validating links in file %s: %v\n", path, err)
			}
		}
		return nil
	}
	err := filepath.Walk(directory, f)

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", directory, err)
		return err
	}

	return nil
}
