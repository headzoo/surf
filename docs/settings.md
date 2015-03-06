# Settings

```go
bow := surf.NewBrowser()

// Set the user agent this browser instance will send with each request.
bow.SetUserAgent("SuperCrawler/1.0")

// Or set the user agent globally so every new browser you create uses it.
browser.DefaultUserAgent = "SuperCrawler/1.0"

// Attributes control how the browser behaves. Use the SetAttribute() method
// to set attributes one at a time.
bow.SetAttribute(browser.SendReferer, false)
bow.SetAttribute(browser.MetaRefreshHandling, false)
bow.SetAttribute(browser.FollowRedirects, false)

// Or set the attributes all at once using SetAttributes().
bow.SetAttributes(browser.AttributeMap{
    browser.SendReferer:         surf.DefaultSendReferer,
    browser.MetaRefreshHandling: surf.DefaultMetaRefreshHandling,
    browser.FollowRedirects:     surf.DefaultFollowRedirects,
})

// The attributes can also be set globally. Now every new browser you create
// will be set with these defaults.
surf.DefaultSendReferer = false
surf.DefaultMetaRefreshHandling = false
surf.DefaultFollowRedirects = false

// Override the build in cookie jar.
// Surf uses cookiejar.Jar by default.
bow.SetCookieJar(jar.NewMemoryCookies())

// Override the build in bookmarks jar.
// Surf uses jar.MemoryBookmarks by default.
bow.SetBookmarksJar(jar.NewMemoryBookmarks())

// Use jar.FileBookmarks to read and write your bookmarks to a JSON file.
bookmarks, err = jar.NewFileBookmarks("/home/joe/bookmarks.json")
if err != nil { panic(err) }
bow.SetBookmarksJar(bookmarks)
```
