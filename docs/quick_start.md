# Quick Start
Start by creating a new \*browser.Browser and making a GET request to golang.org.
```go
bow := surf.NewBrowser()
err := bow.Open("http://golang.org")
if err != nil {
	panic(err)
}

// Outputs: "The Go Programming Language"
fmt.Println(bow.Title())
```

If you need to you can add additional request headers.
```go
bow := surf.NewBrowser()
bow.AddRequestHeader("Accept", "text/html")
bow.AddRequestHeader("Accept-Charset", "utf8")

err := bow.Open("http://golang.org")
if err != nil {
	panic(err)
}

fmt.Println(bow.Title())
```

It's important to note that `Browser.Open()` does not return any kind of response object. Rather, the "state" of
the browser changes to reflect the current page. Calling `Open()` is analogous to typing an URL into
your web browser address bar. The "state" of the browser changes after requesting the new page.

When we open a new page, the state changes to reflection the current page.
```go
err := bow.Open("http://reddit.com")
if err != nil {
	panic(err)
}

// Outputs: "reddit: the front page of the internet"
fmt.Println(bow.Title())
```

Just like a real web browser, Surf maintains a history that you can move back through. You can also
bookmark pages and come back to them later.
```go
// Bookmark the page so we can come back to it later.
err = bow.Bookmark("reddit")
if err != nil {
	panic(err)
}

// Now move back to the golang.org site.
bow.Back()

// And then back to reddit using our bookmark.
bow.OpenBookmark("reddit")
```

By default the bookmarks are kept in memory, and will disappear when your \*browser.Browser instance
is destroyed. See the [settings](settings/#storage-jars) for information on saving your bookmarks to a file.

