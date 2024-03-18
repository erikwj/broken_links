package internal

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestValidateConversion(t *testing.T) {
	header := "This header's title has lots of words"
	expected := "this-headers-title-has-lots-of-words"

	res := convertHeader(header)

	// Assert the expected result
	if res != expected {
		t.Errorf("Expected low case header with dashes but got %s", res)
	}
}
func TestValidateConversion2(t *testing.T) {
	header := "Load-balancing algorithms"
	expected := "load-balancing-algorithms"

	res := convertHeader(header)

	// Assert the expected result
	if res != expected {
		t.Errorf("Expected low case header with dashes but got %s", res)
	}
}

func TestValidateRegexFileMd(t *testing.T) {
	// Prepare the test data
	regexs := ExtDocRegex(".md")
	line := " asdfas df [glossary](../testfiles/glossary.md) or a [correct](../testfiles/correct.md), ... ![image](../img/glossary.png) and [corrupt](../testfiles/corrupt.md)"

	data := regexs.file.FindAllStringSubmatch(line, -1)
	fmt.Println(data)
	// Assert the expected result
	if len(data) != 3 {
		t.Errorf("Expected non empty array of length 3, but got %d", len(data))
	}
}

func TestValidateRegexImageMd(t *testing.T) {
	// Prepare the test data
	regexs := ExtDocRegex(".md")
	line := " asdfas df [glossary](../testfiles/glossary.md) or a [correct](../testfiles/correct.md), ... ![image](../img/glossary.png) and [corrupt](../testfiles/corrupt.md)"

	data := regexs.image.FindAllStringSubmatch(line, -1)
	fmt.Println(data)
	// Assert the expected result
	if len(data) != 1 {
		t.Errorf("Expected non empty array of length 1, but got %d", len(data))
	}
	if data[0][2] != "../img/glossary.png" {
		t.Errorf("Expected link to image but got %s", data[0][2])
	}
}

func TestValidateRegexWebMd(t *testing.T) {
	// Prepare the test data
	regexs := ExtDocRegex(".md")
	line := " asdfas df [glossary](../testfiles/glossary.md) or a [correct](../testfiles/correct.md), ... [glossary](../testfiles/glossary.md) and [corrupt](../testfiles/corrupt.md)"

	data := regexs.web.FindAllStringSubmatch(line, -1)
	fmt.Println(data)
	// Assert the expected result
	if len(data) != 0 {
		t.Errorf("Expected empty array, but got %d", len(data))
	}
}

func TestValidateRegexInternalMd(t *testing.T) {
	// Prepare the test data
	regexs := ExtDocRegex(".md")
	line := " To illustrate, if fixed [retry strategy](../abc/01-03-0002-retry-strategy.md) or a [Retry budget](../../glossary/terms/retry-budget.md). Alternatively, it may implement an [exponential backoff](../d/file.md#Exponential-Backoff), where ... "

	data := regexs.internal.FindAllStringSubmatch(line, -1)
	fmt.Println(data)
	// Assert the expected result
	if len(data) != 1 {
		t.Errorf("Expected non empty array of length 1, but got %d", len(data))
	}
	if data[0][2] != "../d/file.md#Exponential-Backoff" {
		t.Errorf("Expected link to internal reference but got %s", data[0][2])
	}
}

func TestValidateWebUrls(t *testing.T) {
	// Prepare the test data
	var buf bytes.Buffer
	urls := [][]string{
		{"", "", "http://example.com"},
		{"", "", "http://google.com"},
		{"", "", "http://github.com"},
	}
	filePath := "/path/to/file.md"
	lineNum := 1

	// Call the function being tested
	result := validateWebUrls(&buf, urls, filePath, lineNum, false)

	// Assert the expected result
	if result != 0 {
		t.Errorf("Expected result to be 0, but got %d", result)
	}

	// Assert the expected output
	expectedOutput := "open http://example.com # filepath: /path/to/file.md linenumber: 1\n" +
		"open http://google.com # filepath: /path/to/file.md linenumber: 1\n" +
		"open http://github.com # filepath: /path/to/file.md linenumber: 1\n"
	if buf.String() != expectedOutput {
		t.Errorf("Unexpected output.\nExpected:\n%s\nGot:\n%s", expectedOutput, buf.String())
	}
}
func TestValidateSilentWebUrls(t *testing.T) {
	// Prepare the test data
	var buf bytes.Buffer
	urls := [][]string{
		{"", "", "http://example.com"},
		{"", "", "http://google.com"},
		{"", "", "http://github.com"},
	}
	filePath := "/path/to/file.md"
	lineNum := 1

	// Call the function being tested
	result := validateWebUrls(&buf, urls, filePath, lineNum, true)

	// Assert the expected result
	if result != 0 {
		t.Errorf("Expected result to be 0, but got %d", result)
	}

	// Assert the expected output
	expectedOutput := ""
	if buf.String() != expectedOutput {
		t.Errorf("Unexpected output.\nExpected:\n%s\nGot:\n%s", expectedOutput, buf.String())
	}
}

func TestValidateInternalLinks(t *testing.T) {
	// Prepare the test data
	var buf bytes.Buffer
	links := [][]string{
		{"link1", "description1", "../testfiles/correct.md"},
		{"link1", "description1", "../testfiles/glossary.md"},
		{"link1", "description1", "../testfiles/corrupt.md"},
	}
	filePath := "./"
	lineNum := 10

	// Call the function being tested
	result := validateInternalLinks(&buf, links, filePath, lineNum)

	// Assert the expected result
	if result != 0 {
		t.Errorf("Expected validateInternalLinks to return 0, but got %d", result)
	}

	// Assert the output written to the writer
	expectedOutput := ""
	if buf.String() != expectedOutput {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, buf.String())
	}
}

func TestValidateInternalLinksFailure(t *testing.T) {
	// Prepare the test data
	var buf bytes.Buffer
	url := "../testfiles/broken.md"
	url2 := "../testfiles/lost.md"
	url3 := "../testfiles/gone.md"

	links := [][]string{
		{"link1", "description1", url},
		{"link1", "description1", url2},
		{"link1", "description1", url3},
	}
	filePath := "./"
	lineNum := 1

	// Call the function being tested
	result := validateInternalLinks(&buf, links, filePath, lineNum)

	// Assert the expected result
	if result != 1 {
		t.Errorf("Expected validateInternalLinks to return 1, but got %d", result)
	}

	// Assert the output written to the writer
	// function stops after first broken link
	expectedOutput := fmt.Sprintf("\u001b[31m# broken file link in file %s:%d issue: %s\u001b[0m\n", filePath, lineNum, url)

	if buf.String() != expectedOutput {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, buf.String())
	}
}

func TestValidateImageLinks(t *testing.T) {
	// Prepare the test data
	var buf bytes.Buffer
	links := [][]string{
		{"link1", "img1", "../testfiles/img/btn.gif"},
		{"link2", "img2", "../testfiles/img/btn.png"},
		{"link3", "img3", "../testfiles/img/btn.svg"},
	}
	filePath := "../testfiles/correct.md"
	lineNum := 10

	// Call the function being tested
	result := validateImages(&buf, links, filePath, lineNum)

	// Assert the expected result 0 == succes; 1 == failure
	if result != 0 {
		t.Errorf("Expected validateInternalLinks to return 0, but got %d", result)
	}

	// Assert the output written to the writer
	expectedOutput := ""
	if buf.String() != expectedOutput {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, buf.String())
	}
}

func TestValidateImageLinksFailure(t *testing.T) {
	// Prepare the test data
	var buf bytes.Buffer
	links := [][]string{
		{"link1", "img1", "../testfiles/img/btn.gaf"},
		{"link2", "img2", "../testfiles/img/btn.pnf"},
		{"link3", "img3", "../testfiles/img/btn.svf"},
	}
	filePath := "../testfiles/correct.md"
	lineNum := 10

	// Call the function being tested
	result := validateImages(&buf, links, filePath, lineNum)

	// Assert the expected result 0 == succes; 1 == failure
	if result != 1 {
		t.Errorf("Expected validateInternalLinks to return 1, but got %d", result)
	}

	// Assert the output written to the writer
	expectedOutput := fmt.Sprintf("\u001b[31m# broken image file link in file %s:%d issue: %s\u001b[0m\n", filePath, lineNum, links[0][2])

	if buf.String() != expectedOutput {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, buf.String())
	}
}

func TestValidateInternalReferenceLinks(t *testing.T) {
	// Prepare the test data
	var buf bytes.Buffer
	links := [][]string{
		{"link1", "egg", "./subdir/bla.md#headers-2-with-extra-text"},
		{"link1", "find me", "#find-me"},
	}
	filePath := "../testfiles/correct.md"
	lineNum := 13

	// Call the function being tested
	result := validateInternalReferenceLinks(&buf, links, filePath, lineNum)

	// Assert the expected result
	if result != 0 {
		t.Errorf("Expected validateInternalLinks to return 0, but got %d", result)
	}

	// Assert the output written to the writer
	expectedOutput := ""
	if buf.String() != expectedOutput {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, buf.String())
	}
}

func TestValidateInternalReferenceLinksFailure(t *testing.T) {
	// Prepare the test data
	var buf bytes.Buffer
	links := [][]string{
		{"link1", "egg", "./subdir/bla.md#header-with-extra-text"},
		{"link1", "find me", "#found-me"},
	}
	filePath := "../testfiles/correct.md"
	lineNum := 13

	// Call the function being tested
	result := validateInternalReferenceLinks(&buf, links, filePath, lineNum)

	// Assert the expected result
	if result != 1 {
		t.Errorf("Expected validateInternalLinks to return 1, but got %d", result)
	}

	// Assert the output written to the writer
	log := fmt.Sprintf("\u001b[31m# broken header link in file %s:%d issue: %s\u001b[0m\n", filePath, lineNum, links[0][2])
	expectedOutput := log

	if buf.String() != expectedOutput {
		t.Errorf("Expected output:\n%s\nBut got:\n%s", expectedOutput, buf.String())
	}
}

func TestFindHeaders(t *testing.T) {
	absPath := "../testfiles/subdir/bla.md"

	headers, err := findHeaders(absPath)

	if err != nil {
		t.Errorf("Expected FindHeaders to pass, but it failed with error: %v", err)
	}

	expectedHeaders := []string{"title-of-bla", "headers-2-with-extra-text", "level-6-header"}
	if len(headers) != len(expectedHeaders) {
		t.Errorf("Expected %d headers, but got %d", len(expectedHeaders), len(headers))
	}

	for i, header := range headers {
		if header != expectedHeaders[i] {
			t.Errorf("Expected header '%s', but got '%s'", expectedHeaders[i], header)
		}
	}
}

func TestInternalReference(t *testing.T) {
	absPath := "../testfiles/subdir/bla.md"

	if _, err := os.Stat(absPath); err != nil {
		t.Errorf("Expected header to be found but got error: %v", err)
	}
}
