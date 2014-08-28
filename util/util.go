package util

import (
	"os"
	"io"
	"net/http"
)

// FileExists returns if the given file exists.
func FileExists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// Download copies a remote file to the given writer.
//
// Returns the number of bytes written and any errors.
func Download(url string, out io.Writer) (int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return io.Copy(out, resp.Body)
}
