package surf

import (
	"bytes"
	"fmt"
	"github.com/headzoo/surf/element"
	ut "github.com/headzoo/surf/unittest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBrowserGet(t *testing.T) {
	ut.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, html)
	}))
	defer ts.Close()

	b, _ := NewBrowser()
	var _ element.Browsable = b

	err := b.Open(ts.URL)
	ut.AssertNil(err)
	ut.AssertEquals("Surf", b.Title())
	ut.AssertContains("<p>Hello, Surf!</p>", b.Body())

	buff := &bytes.Buffer{}
	l, err := b.Download(buff)
	ut.AssertNil(err)
	ut.AssertGreaterThan(0, int(l))
	ut.AssertEquals(int(l), buff.Len())
}

func TestBrowserBookmarks(t *testing.T) {
	ut.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, html)
	}))
	defer ts.Close()

	b, _ := NewBrowser()
	b.Bookmarks.Save("test1", ts.URL)
	b.OpenBookmark("test1")
	ut.AssertEquals("Surf", b.Title())
	ut.AssertContains("<p>Hello, Surf!</p>", b.Body())

	err := b.BookmarkPage("test2")
	ut.AssertNil(err)
	b.OpenBookmark("test2")
	ut.AssertEquals("Surf", b.Title())
}

func TestBrowserClick(t *testing.T) {
	ut.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprint(w, htmlLinks)
		} else if r.URL.Path == "/page2" {
			fmt.Fprint(w, html)
		}
	}))
	defer ts.Close()

	b, _ := NewBrowser()
	err := b.Open(ts.URL)
	ut.AssertNil(err)

	err = b.Click("a:contains('click')")
	ut.AssertNil(err)
	ut.AssertContains("<p>Hello, Surf!</p>", b.Body())
}

func TestBrowserLinks(t *testing.T) {
	ut.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, htmlLinks)
	}))
	defer ts.Close()

	b, _ := NewBrowser()
	err := b.Open(ts.URL)
	ut.AssertNil(err)

	links := b.Links()
	ut.AssertEquals(2, len(links))
	ut.AssertEquals("", links[0].ID)
	ut.AssertEquals(ts.URL+"/page2", links[0].Href)
	ut.AssertEquals("click", links[0].Text)
	ut.AssertEquals("page3", links[1].ID)
	ut.AssertEquals(ts.URL+"/page3", links[1].Href)
	ut.AssertEquals("no clicking", links[1].Text)
}

func TestBrowserImages(t *testing.T) {
	ut.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, html)
	}))
	defer ts.Close()

	b, _ := NewBrowser()
	err := b.Open(ts.URL)
	ut.AssertNil(err)

	images := b.Images()
	ut.AssertEquals(2, len(images))
	ut.AssertEquals("imgur-image", images[0].ID)
	ut.AssertEquals("http://i.imgur.com/HW4bJtY.jpg", images[0].Src)
	ut.AssertEquals("", images[0].Alt)
	ut.AssertEquals("It's a...", images[0].Title)

	ut.AssertEquals("", images[1].ID)
	ut.AssertEquals(ts.URL+"/Cxagv.jpg", images[1].Src)
	ut.AssertEquals("A picture", images[1].Alt)
	ut.AssertEquals("", images[1].Title)

	buff := &bytes.Buffer{}
	l, err := images[0].Download(buff)
	ut.AssertNil(err)
	ut.AssertGreaterThan(0, buff.Len())
	ut.AssertEquals(int(l), buff.Len())
}

func TestBrowserForm(t *testing.T) {
	ut.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprint(w, htmlForm)
		} else {
			r.ParseForm()
			fmt.Fprint(w, r.Form.Encode())
		}
	}))
	defer ts.Close()

	b, _ := NewBrowser()
	err := b.Open(ts.URL)
	ut.AssertNil(err)

	f, err := b.Form("[name='default']")
	ut.AssertNil(err)

	f.Input("age", "55")
	f.Input("gender", "male")
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertContains("age=55", b.Body())
	ut.AssertContains("gender=male", b.Body())
	ut.AssertContains("submit2=submitted2", b.Body())
}

var html = `<!doctype html>
<html>
	<head>
		<title>Surf</title>
	</head>
	<body>
		<img src="http://i.imgur.com/HW4bJtY.jpg" id="imgur-image" title="It's a..." />
		<p>Hello, Surf!</p>
		<img src="/Cxagv.jpg" alt="A picture" />
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
		<a href="/page3" id="page3">no clicking</a>
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
