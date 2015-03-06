# Overview

```go
// Start by creating a new bow.
bow := surf.NewBrowser()

// Add additional request headers.
bow.AddRequestHeader("Accept", "text/html")
bow.AddRequestHeader("Accept-Charset", "utf8")

// Requesting a page.
err := bow.Open("http://www.reddit.com")
if err != nil { panic(err) }
fmt.Println(bow.Title())
// Outputs: "reddit: the front page of the internet"

// Follow a link on the page where the link text is "new". Surf uses the selector
// engine from goquery, which has a similar syntax to jQuery.
err = bow.Click("a:contains('new')")
if err != nil { panic(err) }
fmt.Println(bow.Title())
// Outputs: "newest submissions: reddit.com"

// Bookmark the page so we can come back to it later.
err = bow.Bookmark("reddit-new")
if err != nil { panic(err) }

// Login to the site via their login form. Again, we're using the goquery selector
// syntax.
fm, err := bow.Form("form.login-form")
if err != nil { panic(err) }
fm.Input("user", "JoeRedditor")
fm.Input("passwd", "d234rlkasd")
err = fm.Submit()
if err != nil { panic(err) }

// Now that we're logged in, follow the link to our profile.
err = bow.Click("a:contains('JoeRedditor')")
if err != nil { panic(err) }
fmt.Println(bow.Title())
// Outputs: "overview for JoeRedditor"

// The underlying goquery.Document is exposed via the Dom() method, which
// can be used to parse values from the body. See the goquery documentation
// for more information on selecting page elements.
// Load our previously saved bookmark, and print the titles for each submission
// on the reddit home page.
err = bow.GetBookmark("reddit-new")
if err != nil { panic(err) }
bow.Dom().Find("a.title").Each(func(_ int, s *goquery.Selection) {
    fmt.Println(s.Text())
})

// The most common Dom() methods can be called directly from the browser.
// The need to find elements on the page is common enough that the above could
// be written like this.
bow.Find("a.title").Each(func(_ int, s *goquery.Selection) {
    fmt.Println(s.Text())
})

// Last, but not least, write the document to a file using the Download()
// method. The Download() method accepts any io.Writer.
file, err := os.Create("reddit.html")
if err != nil { panic(err) }
defer file.Close()
bow.Download(file)
```
