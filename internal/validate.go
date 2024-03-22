package internal

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"os"
	"path/filepath"
	"regexp"
)

type DocRegex struct {
	web      *regexp.Regexp
	file     *regexp.Regexp
	internal *regexp.Regexp
	image    *regexp.Regexp
}

func ValidateLine(line string, lineNum int, filePath string, regexs DocRegex, onlyErrors bool) error {
	// Supported links can only have characters or numbers in the name of the link

	linksError := validateInternalLinks(os.Stdout, regexs.file.FindAllStringSubmatch(line, -1), filePath, lineNum)
	imgError := validateImages(os.Stdout, regexs.image.FindAllStringSubmatch(line, -1), filePath, lineNum)
	webError := validateWebUrls(os.Stdout, regexs.web.FindAllStringSubmatch(line, -1), filePath, lineNum, onlyErrors)
	internalError := validateInternalReferenceLinks(os.Stdout, regexs.internal.FindAllStringSubmatch(line, -1), filePath, lineNum)

	if linksError != 0 || imgError != 0 || webError != 0 || internalError != 0 {
		return fmt.Errorf("\u001b[31m# error validating line in file %s:%d\u001b[0m", filePath, lineNum)
	}
	return nil
}

func validateInternalLinks(w io.Writer, links [][]string, filePath string, lineNum int) int {
	for _, link := range links {
		if check_length(link) {
			continue
		}
		url := link[2]
		absPath, err := filepath.Abs(filepath.Dir(filePath))
		if err != nil {
			err = fmt.Errorf("\u001b[31m# error getting absolute path for file %s:%v\u001b[0m", filePath, err)
			fmt.Fprintln(w, err) // Handle the error appropriately
			continue
		}
		targetPath := filepath.Join(absPath, url)
		if _, err := os.Stat(targetPath); err != nil {
			err = fmt.Errorf("\u001b[31m# broken file link in file %s:%d issue: %s\u001b[0m", filePath, lineNum, url)
			fmt.Fprintln(w, err) // Handle the error appropriately
			return 1
		}

	}
	return 0
}
func validateInternalReferenceLinks(w io.Writer, links [][]string, filePath string, lineNum int) int {
	for _, link := range links {
		if check_length(link) {
			continue
		}
		url := link[2]
		parts := strings.Split(url, "#")
		var header string
		var targetPath string
		var fileName string

		_, fileName = filepath.Split(filePath)
		// If there is a # in the link, split the link into the path and the header

		// Get the root from the file path
		absPath, err := filepath.Abs(filepath.Dir(filePath))

		if err != nil {
			err = fmt.Errorf("\u001b[31m# error getting absolute path for file %s:%v\u001b[0m", filePath, err)
			fmt.Fprintln(w, err) // Handle the error appropriately
			continue
		}
		if len(parts) > 1 {
			if parts[0] == "" {
				fileName = filepath.Base(filePath)
			} else {
				fileName = parts[0]
			}
			header = parts[1]
		}
		targetPath = filepath.Join(absPath, fileName)

		if _, err := os.Stat(targetPath); err != nil {
			err = fmt.Errorf("\u001b[31m# broken reference link in file %s:%d issue: %s\u001b[0m", filePath, lineNum, url)
			fmt.Fprintln(w, err) // Handle the error appropriately
			return 1
		}
		headers, err := findHeaders(targetPath)
		if err != nil {
			// check if header exists in headers
			err = fmt.Errorf("\u001b[31m# error getting headers for file %s:%v\u001b[0m", filePath, err)
			fmt.Fprintln(w, err) // Handle the error appropriately
			continue
		}
		headerExists := false
		for _, h := range headers {
			if h == header {
				headerExists = true
				break
			}
		}
		if !headerExists {
			err = fmt.Errorf("\u001b[31m# broken header link in file %s:%d issue: %s\u001b[0m", filePath, lineNum, url)
			fmt.Fprintln(w, err) // Handle the error appropriately
			return 1
		}

	}
	return 0
}
func findHeaders(absPath string) ([]string, error) {
	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	headerRegex := regexp.MustCompile(`^#{1,6} (.*)$`)
	var headers []string

	for scanner.Scan() {
		line := scanner.Text()
		matches := headerRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			headers = append(headers, convertHeader(matches[1]))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return headers, nil
}

func convertHeader(header string) string {
	return strings.ReplaceAll(removeSpecialCharacters(strings.ToLower(header)), " ", "-")
}

func removeSpecialCharacters(input string) string {
	// Regular expression to match non-alphanumeric characters
	regex := regexp.MustCompile("[^a-zA-Z0-9 -]+")
	// Replace non-alphanumeric characters with an empty string
	cleaned := regex.ReplaceAllString(input, "")
	return cleaned
}

func validateImages(w io.Writer, images [][]string, filePath string, lineNum int) int {
	for _, link := range images {
		if check_length(link) {
			continue
		}
		url := link[2]
		absPath, err := filepath.Abs(filepath.Dir(filePath))
		if err != nil {
			err = fmt.Errorf("\u001b[31m# error getting absolute path for image file %s:%v\u001b[0m", filePath, err)
			fmt.Fprintln(w, err) // Handle the error appropriately
			continue
		}
		targetPath := filepath.Join(absPath, url)
		if _, err := os.Stat(targetPath); err != nil {
			err = fmt.Errorf("\u001b[31m# broken image file link in file %s:%d issue: %s\u001b[0m", filePath, lineNum, url)
			fmt.Fprintln(w, err) // Handle the error appropriately
			return 1
		}
	}
	return 0
}

func validateWebUrls(w io.Writer, urls [][]string, filePath string, lineNum int, onlyErrors bool) int {
	for _, link := range urls {
		if check_length(link) {
			continue
		}
		url := link[2]

		if !onlyErrors {
			fmt.Fprintf(w, "open %s # filepath: %s:%d\n", url, filePath, lineNum)
		}
	}
	return 0
}

func ValidateLinks(filePath string, extension string, onlyErrors bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var validateError error = nil

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++
		validateError = ValidateLine(line, lineNum, filePath, ExtDocRegex(extension), onlyErrors)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return validateError
}

// default on markdown
func ExtDocRegex(extension string) DocRegex {
	switch extension {
	case ".rst":
		return DocRegex{
			file:     regexp.MustCompile(""), // not supported
			web:      regexp.MustCompile("`(.*) <(https?://[-%()_.!~*'#;/?:@&=+$,A-Za-z0-9]+)>`_"),
			image:    regexp.MustCompile(`(::image )(.*.[png|svg|gif])`),
			internal: regexp.MustCompile(`\[([a-zA-Z0-9 ]+)\]\((#[^)]+|.*\.md#[^)]+)\)`),
		}
	default:
		return DocRegex{
			file:     regexp.MustCompile(`\[([a-zA-Z0-9 ]+)\]\(([^)]+.md)\)`),
			web:      regexp.MustCompile(`\[([a-zA-Z0-9 ]+)\]\((https?://[-%()_.!~*'#;/?:@&=+$,A-Za-z0-9]+)\)`),
			image:    regexp.MustCompile(`!\[(.*)\]\(([^)]+.[png|svg|gif])\)`),
			internal: regexp.MustCompile(`\[([a-zA-Z0-9 ]+)\]\((#[^)]+|[^)]+\.md\#[^)]+)\)`),
		}
	}
}

func check_length(arr []string) bool {
	return len(arr) != 3
}
