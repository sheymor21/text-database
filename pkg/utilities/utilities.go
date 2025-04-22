package utilities

import (
	"errors"
	"os"
)

func IsFileExist(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func RemoveEmptyIndex(data []string) []string {

	filteredData := make([]string, 0)
	for _, d := range data {
		if d != "" && d != "\r\n" {
			filteredData = append(filteredData, d)
		}
	}
	return filteredData
}
