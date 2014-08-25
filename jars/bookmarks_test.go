package jars

import (
	ut "github.com/headzoo/surf/unittest"
	"testing"
)

func TestMemoryBookmarks(t *testing.T) {
	ut.Run(t)

	b := NewMemoryBookmarks()

	err := b.Save("test1", "http://localhost")
	ut.AssertNil(err)
	err = b.Save("test2", "http://127.0.0.1")
	ut.AssertNil(err)
	err = b.Save("test1", "http://localhost")
	ut.AssertNotNil(err)

	url, err := b.Read("test1")
	ut.AssertNil(err)
	ut.AssertEquals("http://localhost", url)
	url, err = b.Read("test2")
	ut.AssertEquals("http://127.0.0.1", url)
	url, err = b.Read("test3")
	ut.AssertNotNil(err)

	r := b.Remove("test2")
	ut.AssertTrue(r)
	r = b.Remove("test3")
	ut.AssertFalse(r)

	r = b.Has("test1")
	ut.AssertTrue(r)
	r = b.Has("test4")
	ut.AssertFalse(r)
}
