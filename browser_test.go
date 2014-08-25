package surf

import (
	"fmt"
	"github.com/headzoo/surf/unittest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBrowserGet(t *testing.T) {
	unittest.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, html)
	}))
	defer ts.Close()

	b := NewBrowser()
	err := b.Get(ts.URL)
	unittest.AssertEquals(nil, err)
	unittest.AssertEquals("Surf", b.Title())
	unittest.AssertContains("<p>Hello, Surf!</p>", b.Body())

}

func TestBrowseFollowLink(t *testing.T) {
	unittest.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprint(w, htmlLinks)
		} else if r.URL.Path == "/page2" {
			fmt.Fprint(w, html)
		}
	}))
	defer ts.Close()

	b := NewBrowser()
	err := b.Get(ts.URL)
	unittest.AssertEquals(nil, err)

	err = b.FollowLink(":contains('click')")
	unittest.AssertEquals(nil, err)
	unittest.AssertContains("<p>Hello, Surf!</p>", b.Body())
}

func TestBrowseForm(t *testing.T) {
	unittest.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprint(w, htmlForm)
		} else {
			r.ParseForm()
			fmt.Fprint(w, r.Form.Encode())
		}
	}))
	defer ts.Close()

	b := NewBrowser()
	err := b.Get(ts.URL)
	unittest.AssertEquals(nil, err)

	f, err := b.Form("[name='default']")
	unittest.AssertEquals(nil, err)

	f.Input("age", "55")
	f.Input("gender", "male")
	err = f.Click("submit2")
	unittest.AssertEquals(nil, err)
	unittest.AssertContains("age=55", b.Body())
	unittest.AssertContains("gender=male", b.Body())
	unittest.AssertContains("submit2=submitted2", b.Body())
}

var html = `<!doctype html>
<html>
	<head>
		<title>Surf</title>
	</head>
	<body>
		<p>Hello, Surf!</p>
	</body>
</html>
`

var htmlLinks = `<!doctype html>
<html>
	<head>
		<title>Surf</title>
	</head>
	<body>
		<p>Click the link below.</p>
		<a href="/page2">click</a>
	</body>
</html>
`

var htmlForm = `<!doctype html>
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
