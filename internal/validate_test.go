package internal_test

import (
	"github.com/erikwj/brokenlinks/internal"
	"testing"
)

func TestValidateWebLine(t *testing.T) {
	line := "[GitHub](http://github.com) (and some extra text) [Gitlab](http://gitlab.com) "
	lineNum := 1
	filePath := "/path/to/file.md"

	// Test your validateLine function here
	// with the given line, lineNum, and filePath variables as input
	err := internal.ValidateLine(line, lineNum, filePath)

	// Assert the expected result
	if err != nil {
		t.Errorf("Expected validateWebLine to pass, but it failed")
	}
}

func TestValidateFileLine(t *testing.T) {
	line := "[glossary](../testfiles/glossary.md)"
	lineNum := 1
	filePath := "./"

	// Test your validateLine function here
	// with the given line, lineNum, and filePath variables as input
	err := internal.ValidateLine(line, lineNum, filePath)

	// Assert the expected result
	if err != nil {
		t.Errorf("Expected validateLine to pass, but it failed")
	}
}

func TestValidateImgLine(t *testing.T) {
	line := "![](../testfiles/img/btn.png)"
	lineNum := 1
	filePath := "./"

	// Test your validateLine function here
	// with the given line, lineNum, and filePath variables as input
	err := internal.ValidateLine(line, lineNum, filePath)

	// Assert the expected result
	if err != nil {
		t.Errorf("Expected validateLine to pass, but it failed")
	}
}

// test for failure of ValidateLine not working since no active checking on web urls
// func TestValidateWebLineFail(t *testing.T) {
// 	line := "[GitHub](https://github.c)"
// 	lineNum := 1
// 	filePath := "/path/to/file.md"

// 	// Test your validateLine function here
// 	err := validateLine(line, lineNum, filePath)

// 	// Assert that the function fails
// 	if err == nil {
// 		t.Errorf("Expected validateLine to fail, but it succeeded")
// 	}
// }

// test for failure of ValidateLine with broken file link
func TestValidateFileLineFail(t *testing.T) {
	line := "[GitHub](broken.md)"
	lineNum := 1
	filePath := "/path/to/file.md"

	// Test your validateLine function here
	err := internal.ValidateLine(line, lineNum, filePath)

	// Assert that the function fails
	if err == nil {
		t.Errorf("Expected validateLine to fail, but it succeeded")
	}
}

// test for failure of ValidateLine with broken image link
func TestValidateImageLineFail(t *testing.T) {
	line := "![](broken.png)"
	lineNum := 1
	filePath := "/path/to/file.md"

	// Test your validateLine function here
	err := internal.ValidateLine(line, lineNum, filePath)

	// Assert that the function fails
	if err == nil {
		t.Errorf("Expected validateLine to fail, but it succeeded")
	}
}
