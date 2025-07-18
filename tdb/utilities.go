package tdb

import (
	"errors"
	"os"
)

// isFileExist checks if a file exists at the specified path.
// It takes a filename string parameter and returns true if the file exists,
// false otherwise.
func isFileExist(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// removeEmptyIndex filters out empty strings and newline characters from a slice of strings.
// It takes a slice of strings as input and returns a new slice containing only non-empty strings.
func removeEmptyIndex(data []string) []string {

	filteredData := make([]string, 0)
	for _, d := range data {
		if d != "" && d != "\r\n" {
			filteredData = append(filteredData, d)
		}
	}
	return filteredData
}
