Surf
====
Surf is a Go (golang) library that implements a virtual browser that you can control pragmatically. Just like a real browser you can open pages, follow links, bookmark pages, submit forms, and many other things. Surf is modeled after Andy Lester's Perl module [WWW::Mechanize](http://search.cpan.org/~ether/WWW-Mechanize-1.73/lib/WWW/Mechanize.pm).

[Complete API documentation is available from the GoDoc website.](http://godoc.org/github.com/headzoo/surf)

_This project is very young, and the API is bound to change often. Use at your own risk. The master branch is the stable branch, while future work is being done on the dev branch._

* [Installation](#installation)
* [Usage](#usage)
* [Settings](#settings)
* [User Agents](#user-agents)
* [Downloading](#downloading)
* [Credits](#credits)
* [Use Cases](#use-cases)
* [TODO](#todo)

### Installation
Download the library using go.  
`go get github.com/headzoo/surf`

You'll need the ut library if you want to run the unit tests.  
`go get github.com/headzoo/ut`  

Import the library into your project.  
`import "github.com/headzoo/surf"`


### Usage
```go
// Start by creating a new bow.
bow, err := surf.NewBrowser()
if err != nil { panic(err) }

// Add additional request headers.
bow.AddHeader("Accept", "text/html")
bow.AddHeader("Accept-Charset", "utf8")


// Requesting a page.
err = bow.Open("http://www.reddit.com")
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
err = bow.BookmarkPage("reddit-new")
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

// Move back to the home page, and print the page body.
err = bow.Back()
if err != nil { panic(err) }
fmt.Println(bow.Body())


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


### Settings
```go
bow, err := surf.NewBrowser()
if err != nil { panic(err) }

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
    browser.SendReferer:         browser.DefaultSendReferer,
    browser.MetaRefreshHandling: browser.DefaultMetaRefreshHandling,
    browser.FollowRedirects:     browser.DefaultFollowRedirects,
})

// The attributes can also be set globally. Now every new browser you create
// will be set with these defaults.
browser.DefaultSendReferer = false
browser.DefaultMetaRefreshHandling = false
browser.DefaultFollowRedirects = false

// Override the build in cookie jar.
// Surf uses cookiejar.Jar by default.
cookies, err := cookiejar.New(nil)
if err != nil { panic(err) }
bow.SetCookieJar(cookies)

// Override the build in bookmarks jar.
// Surf uses jar.MemoryBookmarks by default.
bookmarks, err := jar.NewMemoryBookmarks()
if err != nil { panic(err) }
bow.SetBookmarksJar(bookmarks)

// Use jar.FileBookmarks to read and write your bookmarks to a JSON file.
bookmarks, err = jar.NewFileBookmarks("/home/joe/bookmarks.json")
if err != nil { panic(err) }
bow.SetBookmarksJar(bookmarks)
```


### User Agents
The agent package contains a number of methods for creating user agent strings for popular browsers and crawlers, and for generating your own user agents.
```go
bow, err := surf.NewBrowser()
if err != nil { panic(err) }

// Use the Google Chrome user agent. The Chrome() method returns:
// "Mozilla/5.0 (Windows NT 6.3; x64) Chrome/37.0.2049.0 Safari/537.36".
bow.SetUserAgent(agent.Chrome())

// The Firefox() method returns:
// "Mozilla/5.0 (Windows NT 6.3; x64; rv:31.0) Gecko/20100101 Firefox/31.0".
bow.SetUserAgent(agent.Firefox())

// The Safari() method returns:
// "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_8) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Safari/8536.25".
bow.SetUserAgent(agent.Safari())

// There are methods for a number of bows and crawlers. For example
// Opera(), MSIE(), AOL(), GoogleBot(), and many more. You can even choose
// the bow version. This will create:
// "Mozilla/5.0 (Windows NT 6.3; x64) Chrome/35 Safari/537.36".
ua := agent.CreateVersion("chrome", "35")
bow.SetUserAgent(ua)

// Creating your own custom user agent is just as easy. The following code
// generates the user agent:
// "Mybow/1.0 (Windows NT 6.1; WOW64; x64)".
agent.Name = "Mybow"
agent.Version = "1.0"
agent.OSName = "Windows NT"
agent.OSVersion = "6.1"
agent.Comments = []string{"WOW64", "x64"}
bow.SetUserAgent(agent.Create())
```
The agent package has an internal database for many different versions of many different browsers. See the [agent package API documentation](http://godoc.org/github.com/headzoo/surf/agent) for more information.


### Downloading
Surf makes it easy to download page assets, such as images, stylesheets, and scripts. They can even be downloaded asynchronously.
```go
bow, err := surf.NewBrowser()
if err != nil { panic(err) }
err = bow.Open("http://www.reddit.com")
if err != nil { panic(err) }

// Download the images on the page and write them to files.
for _, image := range bow.Images() {
    fout, err := os.Create("/home/joe/Pictures" + image.URL.Path())
    if err != nil { panic(err) }
    image.Download(fout)
    fout.Close()
}

// Download the images asynchronously and write them to files.
//
// The DownloadAsync() takes an instance of browser.AsyncDownloadChannel, which
// is a channel upon which the method will put an instance of AsyncDownloadResult
// when the download is complete. The AsyncDownloadResult has the downloaded data,
// or an error if there was one.
downloadQueue := 0
results := make(browser.AsyncDownloadChannel, 1)
for _, image := range bow.Images() {
    image.DownloadAsync(results)
    downloadQueue++
}

for {
    select {
    case result := <-results:
    	// The *browser.Image is stored in result.Asset, but that's of type
    	// interface{} (because AsyncDownloadResult handles many types of
    	// downloads). Use type assertion to get the *browser.Image.
    	image := result.Asset.(*browser.Image)
    	if result.Error == nil {
    		// The image data is saved in result.Data, which is type []byte.
    		filename := "/home/joe/Pictures" + image.URL.Path()
    		ioutil.WriteFile(filename, result.Data, 0700)
        } else {
        	fmt.Printf("Error downloading %s.", image.URL.String())
        }
        
        downloadQueue--
        if downloadQueue == 0 {
            goto COMPLETE
        }
    }
}
COMPLETE:
	close(results)
	fmt.Println("Downloads complete!")
```


### Credits
Surf uses the awesome [goquery](https://github.com/PuerkitoBio/goquery) by Martin Angers, and was written using [Intellij](http://www.jetbrains.com/idea/) and the [golang plugin](http://plugins.jetbrains.com/plugin/5047).


### Use Cases
* Interacting with sites that do not have public APIs.
* Testing/Stressing your sites.
* Scraping sites.
* Creating a web crawler.


### TODO
* Add user authentication.
* Run JavaScript found in the page?
* Add AttributeDownloadAssets so the browser downloads the images, scripts, stylesheets, etc.
* Write more tests. 
* File uploading in forms.
* Handle checkboxes correctly.
