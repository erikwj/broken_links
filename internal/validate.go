package internal

import (
	"bufio"
	"fmt"

	// "net/http"
	"os"
	"path/filepath"
	"regexp"
)

type DocRegex struct {
	web   *regexp.Regexp
	file  *regexp.Regexp
	image *regexp.Regexp
}

func ValidateLine(line string, lineNum int, filePath string, regexs DocRegex) error {
	// Supported links can only have characters or numbers in the name of the link
	// httpregex := regexp.MustCompile(`\[([a-zA-Z0-9 ]+)\]\((https?://[-%()_.!~*';/?:@&=+$,A-Za-z0-9]+)\)`)
	//fileregex := regexp.MustCompile(`\[(.*)\]\((.*.md)\)`)
	// imgregex := regexp.MustCompile(`!\[(.*)\]\((.*.[png|svg|gif])\)`)

	linksError := validateInternalLinks(regexs.file.FindAllStringSubmatch(line, -1), filePath, lineNum)
	imgError := validateImages(regexs.image.FindAllStringSubmatch(line, -1), filePath, lineNum)
	webError := validateWebUrls(regexs.web.FindAllStringSubmatch(line, -1), filePath, lineNum)

	if linksError != 0 || imgError != 0 || webError != 0 {
		return fmt.Errorf("\u001b[31m# error validating line in file %s:%d\u001b[0m", filePath, lineNum)
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
			err = fmt.Errorf("\u001b[31m# error getting absolute path for file %s:%v\u001b[0m", filePath, err)
			fmt.Println(err) // Handle the error appropriately
			continue
		}
		targetPath := filepath.Join(absPath, url)
		if _, err := os.Stat(targetPath); err != nil {
			err = fmt.Errorf("\u001b[31m# broken file link in file %s:%d issue: %s\u001b[0m", filePath, lineNum, url)
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
			err = fmt.Errorf("\u001b[31m# error getting absolute path for image file %s:%v\u001b[0m", filePath, err)
			fmt.Println(err) // Handle the error appropriately
			continue
		}
		targetPath := filepath.Join(absPath, url)
		if _, err := os.Stat(targetPath); err != nil {
			err = fmt.Errorf("\u001b[31m# broken image file link in file %s:%d issue: %s\u001b[0m ", filePath, lineNum, url)
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
		// Below code doesn't work since lots of pages don't return 404 on broken links
		// Therefore only print the urls to be opened via commandline
		// if resp, err := http.Get(url); err != nil || resp.StatusCode != 200 {
		// 	err = fmt.Errorf("broken web link in file %s:%d issue: %s", filePath, lineNum, url)
		// 	fmt.Print(err) // Handle the error appropriately

		// 	return 1
		// }
	}
	return 0
}

func ValidateLinks(filePath string, extension string) error {
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
		ValidateLine(line, lineNum, filePath, ExtDocRegex(extension))
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// default on markdown
func ExtDocRegex(extension string) DocRegex {
	switch extension {
	case ".rst":
		return DocRegex{
			// file:  regexp.MustCompile("(:ref:)`([^`]*)`"), // doesn't work due to some macro stuff or so
			file:  regexp.MustCompile(""),
			web:   regexp.MustCompile("`(.*) <(https?://[-%()_.!~*'#;/?:@&=+$,A-Za-z0-9]+)>`_"),
			image: regexp.MustCompile(`(::image )(.*.[png|svg|gif])`),
		}
	default:
		return DocRegex{
			file:  regexp.MustCompile(`\[(.*)\]\((.*.md)\)`),
			web:   regexp.MustCompile(`\[([a-zA-Z0-9 ]+)\]\((https?://[-%()_.!~*'#;/?:@&=+$,A-Za-z0-9]+)\)`),
			image: regexp.MustCompile(`!\[(.*)\]\((.*.[png|svg|gif])\)`),
		}
	}
}

func check_length(arr []string) bool {
	return len(arr) != 3
}
