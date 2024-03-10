package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func main() {

	directory := "."
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".md" {
			fmt.Printf("Validating %s \n", path)
			if err := validateLinks(path); err != nil {
				fmt.Printf("Error validating links in file %s: %v\n", path, err)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", directory, err)
	}
}

func validateLinks(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Supported links
	httpregex := regexp.MustCompile(`\[(.*)\]\((http.*)\)`)
	fileregex := regexp.MustCompile(`\[(.*)\]\((.*.md)\)`)
	imgregex := regexp.MustCompile(`!\[(.+)\]\((.*.[png|svg|gif])\)`)

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		links := fileregex.FindAllStringSubmatch(line, -1)
		urls := httpregex.FindAllStringSubmatch(line, -1)
		images := imgregex.FindAllStringSubmatch(line, -1)

		// Validate internal links
		for _, link := range links {
			if check_length(link) {
				continue
			}
			url := link[2]
			absPath, err := filepath.Abs(filepath.Dir(filePath))
			if err != nil {
				fmt.Printf("Error getting absolute path for file %s: %v\n", filePath, err)
				continue
			}
			targetPath := filepath.Join(absPath, url)
			if _, err := os.Stat(targetPath); err != nil {
				fmt.Printf("Broken file link in file %s at line %d: %s\n", filePath, lineNum, url)
			}
		}

		// Validate images
		for _, link := range images {
			if check_length(link) {
				continue
			}
			url := link[2]
			absPath, err := filepath.Abs(filepath.Dir(filePath))
			if err != nil {
				fmt.Printf("Error getting absolute path for image file %s: %v\n", filePath, err)
				continue
			}
			targetPath := filepath.Join(absPath, url)
			if _, err := os.Stat(targetPath); err != nil {
				fmt.Printf("Broken image file link in file %s at line %d: %s\n", filePath, lineNum, url)
			}
		}

		// Validate web urls
		for _, link := range urls {
			if check_length(link) {
				continue
			}
			url := link[2]
			if _, err := http.Get(url); err != nil {
				fmt.Printf("Broken web link in file %s at line %d: %s\n", filePath, lineNum, url)
			}
		}

	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func check_length(arr []string) bool {
	return len(arr) != 3
}
