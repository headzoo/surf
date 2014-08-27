package util

import "os"

// FileExists returns if the given file exists.
func FileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
