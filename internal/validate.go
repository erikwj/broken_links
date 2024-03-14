package internal

import (
	"bufio"
	"fmt"

	// "net/http"
	"os"
	"path/filepath"
	"regexp"
)

func ValidateLine(line string, lineNum int, filePath string) error {
	// Supported links can only have characters or numbers in the name of the link
	httpregex := regexp.MustCompile(`\[([a-zA-Z0-9 ]+)\]\((https?://[-%()_.!~*';/?:@&=+$,A-Za-z0-9]+)\)`)
	fileregex := regexp.MustCompile(`\[(.*)\]\((.*.md)\)`)
	imgregex := regexp.MustCompile(`!\[(.*)\]\((.*.[png|svg|gif])\)`)

	linksError := validateInternalLinks(fileregex.FindAllStringSubmatch(line, -1), filePath, lineNum)
	imgError := validateImages(imgregex.FindAllStringSubmatch(line, -1), filePath, lineNum)
	webError := validateWebUrls(httpregex.FindAllStringSubmatch(line, -1), filePath, lineNum)

	if linksError != 0 || imgError != 0 || webError != 0 {
		return fmt.Errorf("error validating line in file %s at line %d", filePath, lineNum)
	}
	return nil
}

func validateInternalLinks(links [][]string, filePath string, lineNum int) int {
	for _, link := range links {
		if check_length(link) {
			continue
		}
		url := link[2]
		absPath, err := filepath.Abs(filepath.Dir(filePath))
		if err != nil {
			err = fmt.Errorf("error getting absolute path for file %s: %v", filePath, err)
			fmt.Println(err) // Handle the error appropriately
			continue
		}
		targetPath := filepath.Join(absPath, url)
		if _, err := os.Stat(targetPath); err != nil {
			err = fmt.Errorf("broken file link in file %s:%d issue: %s", filePath, lineNum, url)
			fmt.Println(err) // Handle the error appropriately
			return 1
		}

	}
	return 0
}

func validateImages(images [][]string, filePath string, lineNum int) int {
	for _, link := range images {
		if check_length(link) {
			continue
		}
		url := link[2]
		absPath, err := filepath.Abs(filepath.Dir(filePath))
		if err != nil {
			err = fmt.Errorf("error getting absolute path for image file %s: %v", filePath, err)
			fmt.Println(err) // Handle the error appropriately
			continue
		}
		targetPath := filepath.Join(absPath, url)
		if _, err := os.Stat(targetPath); err != nil {
			err = fmt.Errorf("broken image file link in file %s:%d issue: %s", filePath, lineNum, url)
			fmt.Println(err) // Handle the error appropriately
			return 1
		}
	}
	return 0
}

func validateWebUrls(urls [][]string, filePath string, lineNum int) int {
	for _, link := range urls {
		if check_length(link) {
			continue
		}
		url := link[2]
		fmt.Println("open", url, "# filepath:", filePath, "linenumber:", lineNum)
		// if resp, err := http.Get(url); err != nil || resp.StatusCode != 200 {
		// 	err = fmt.Errorf("broken web link in file %s:%d issue: %s", filePath, lineNum, url)
		// 	fmt.Print(err) // Handle the error appropriately

		// 	return 1
		// }
	}
	return 0
}

func ValidateLinks(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		ValidateLine(line, lineNum, filePath)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func check_length(arr []string) bool {
	return len(arr) != 3
}
