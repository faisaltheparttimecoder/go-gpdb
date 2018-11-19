package main

import (
	"os"
	"strings"
)

// Function that checks if the string is available on a array.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Check is the value is empty
func IsValueEmpty(v string) bool {
	if len(strings.TrimSpace(v)) == 0 {
		return true
	}
	return false
}

// exists returns whether the given file or directory exists or not
func doesFileOrDirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// Create directory
func CreateDir(path string) {
	// Check if the path or directory exists
	exists, err := doesFileOrDirExists(path)
	if err != nil {
		Fatalf("Failed to check the directory status, the error: %v", err)
	}
	// If not exists then create one
	if !exists {
		Warnf("Directory \"%s\" does not exists, creating one", path)
		err := os.MkdirAll(path, 0755)
		if err != nil {
			Fatalf("Failed to create the directory, the error: %v", err)
		}
	}
}