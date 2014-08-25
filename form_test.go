package gosurf

import (
	"github.com/headzoo/gosurf/unittest"
	"testing"
	"fmt"
	"net/http/httptest"
	"net/http"
)

func TestFormAttribs(t *testing.T) {
	unittest.Run(t)
	unittest.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, htmlForm)
	}))
	defer ts.Close()

	b := NewBrowser()
	b.Get(ts.URL)
	f, err := b.Form("[name='default']")
	unittest.AssertEquals(nil, err)
	unittest.AssertEquals("POST", f.Method())
	unittest.AssertEquals(ts.URL+"/", f.Action())

}
