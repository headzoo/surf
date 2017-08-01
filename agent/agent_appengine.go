// +build appengine

package agent

import (
	"runtime"
)

// osName returns the name of the OS.
func osName() string {
	return runtime.GOOS
}

// osVersion returns the OS version.
func osVersion() string {
	return "0.0"
}
