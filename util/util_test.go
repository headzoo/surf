package util

import (
	ut "github.com/headzoo/surf/unittest"
	"testing"
)

func TestFileExists(t *testing.T) {
	ut.Run(t)

	ex := FileExists("./util_test.go")
	ut.AssertTrue(ex)

	ex = FileExists("./util.txt")
	ut.AssertFalse(ex)
}
