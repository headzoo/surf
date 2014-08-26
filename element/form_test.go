package element

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
		fmt.Fprint(w, htmlTestForm)
	}))
	defer ts.Close()
	/*
	b, _ := NewBrowser()
	b.Open(ts.URL)
	f, err := b.Form("[name='default']")
	ut.AssertNil(err)
	ut.AssertEquals("POST", f.Method())
	ut.AssertEquals(ts.URL+"/", f.Action())
	*/
}

var htmlTestForm = `<!doctype html>
<html>
	<head>
		<title>Echo Form</title>
	</head>
	<body>
		<form method="post" action="/" name="default">
			<input type="text" name="age" value="" />
			<input type="radio" name="gender" value="male" />
			<input type="radio" name="gender" value="female" />
			<input type="submit" name="submit1" value="submitted1" />
			<input type="submit" name="submit2" value="submitted2" />
		</form>
	</body>
</html>
`
