package core

import (
	"os"
)

// function to check if file exists
func DoesFileExist(fileName string) bool {
	file, error := os.Stat(fileName)
	// check if error is "file not exists"
	if os.IsNotExist(error) {
		return false
	}
	if file.IsDir() {
		return false
	}
	return true
}
