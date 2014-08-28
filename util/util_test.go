package util

import (
	ut "github.com/headzoo/surf/unittest"
	"testing"
	"bytes"
)

func TestFileExists(t *testing.T) {
	ut.Run(t)

	ex := FileExists("./util_test.go")
	ut.AssertTrue(ex)

	ex = FileExists("./util.txt")
	ut.AssertFalse(ex)
}

func TestDownload(t *testing.T) {
	ut.Run(t)

	buff := &bytes.Buffer{}
	l, err := Download("http://i.imgur.com/HW4bJtY.jpg", buff)
	ut.AssertNil(err)
	ut.AssertGreaterThan(0, int(l))
	ut.AssertEquals(int(l), buff.Len())
}
