package surf

import (
	"fmt"
	ut "github.com/headzoo/surf/unittest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFormAttribs(t *testing.T) {
	ut.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, htmlForm)
	}))
	defer ts.Close()

	b, _ := NewBrowser()
	b.Get(ts.URL)
	f, err := b.Form("[name='default']")
	ut.AssertNil(err)
	ut.AssertEquals("POST", f.Method())
	ut.AssertEquals(ts.URL+"/", f.Action())
}
