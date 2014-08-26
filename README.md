Surf
====
Surf is a Go (golang) library that implements a virtual browser that you can control pragmatically. Just like a real browser you can open pages, follow links, bookmark pages, submit forms, and many other things. Surf is modeled after John J. Lee's Python library [mechanize](https://github.com/jjlee/mechanize).


* [Installation](#installation)
* [Usage](#usage)
* [Settings](#settings)
* [User Agents](#user-agents)
* [Credits](#credits)
* [Use Cases](#use-cases)
* [TODO](#todo)

### Installation
Download the library using go.  
`go get github.com/headzoo/surf`

Import the library into your project.  
`import "github.com/headzoo/surf"`


### Usage
```go
// Start by creating a new browser.
browser, err := surf.NewBrowser()
if err != nil { panic(err) }

// Set additional request headers.
browser.RequestHeaders.Add("Accept", "text/html")
browser.RequestHeaders.Add("Accept-Charset", "utf8")


// Requesting a page.
err = browser.Get("http://www.reddit.com")
if err != nil { panic(err) }
fmt.Println(browser.Title())
// Outputs: "reddit: the front page of the internet"


// Follow a link on the page where the link text is "new". Surf uses the selector
// engine from goquery, which has a similar syntax to jQuery. With the FollowLink()
// method the "a" is explicit. The selector below is actually "a:contains('new')".
err = browser.FollowLink(":contains('new')")
if err != nil { panic(err) }
fmt.Println(browser.Title())
// Outputs: "newest submissions: reddit.com"


// Bookmark the page so we can come back to it later.
err = browser.BookmarkPage("reddit-new")
if err != nil { panic(err) }


// Login to the site via their login form. Again, we're using the goquery selector
// syntax. The "form" is explicit. The selector below is actually "form.login-form".
fm, err := browser.Form(".login-form")
if err != nil { panic(err) }
fm.Input("user", "JoeRedditor")
fm.Input("passwd", "d234rlkasd")
err = fm.Submit()
if err != nil { panic(err) }


// Now that we're logged in, follow the link to our profile.
err = browser.FollowLink(":contains('JoeRedditor')")
if err != nil { panic(err) }
fmt.Println(browser.Title())
// Outputs: "overview for JoeRedditor"

// Move back to the home page, and print the page body.
err = browser.Back()
if err != nil { panic(err) }
fmt.Println(browser.Body())


// The underlying goquery.Selection is exposed and can be used to parse
// values from the body. Load our previously saved bookmark, and print
// the titles for each submission on the reddit home page.
err = browser.GetBookmark("reddit-new")
if err != nil { panic(err) }
browser.Query().Find("a.title").Each(func(_ int, s *goquery.Selection) {
    fmt.Println(s.Text())
})
```
See the [API documentation](https://github.com/headzoo/surf/tree/master/docs) for more information.


### Settings
```go
browser, err := surf.NewBrowser()
if err != nil { panic(err) }

// Override the default user agent.
browser.UserAgent = "MyBrowser"

// Set the user agent globally. Each Browser instance you create will use this.
surf.DefaultUserAgent = "MyBrowser"


// Attributes control how the browser behaves.
browser.SetAttribute(surf.SendRefererAttribute, false)
browser.SetAttribute(surf.MetaRefreshHandlingAttribute, false)
browser.SetAttribute(surf.FollowRedirectsAttribute, false)

// The attributes may also be set globally.
surf.DefaultSendRefererAttribute = false
surf.DefaultMetaRefreshHandlingAttribute = false
surf.DefaultFollowRedirectsAttribute = false


// Override the build in cookie jar.
cookies, err := cookiejar.New(nil)
if err != nil { panic(err) }
browser.Cookies = cookies

// Override the build in bookmarks container.
bookmarks, err := jars.NewMemoryBookmarks()
if err != nil { panic(err) }
browser.Bookmarks = bookmarks
```
See the [API documentation](https://github.com/headzoo/surf/tree/master/docs) for more information.


### User Agents
The agents package contains a number of methods for creating user agent strings for popular browsers and crawlers, and for generating your own user agents.
```go
browser, err := surf.NewBrowser()
if err != nil { panic(err) }

// Use the Google Chrome user agent. The Chrome() method returns:
// "Mozilla/5.0 (Windows NT 6.3; x64) Chrome/37.0.2049.0 Safari/537.36".
browser.UserAgent = agents.Chrome()

// The Firefox() method returns:
// "Mozilla/5.0 (Windows NT 6.3; x64; rv:31.0) Gecko/20100101 Firefox/31.0".
browser.UserAgent = agents.Firefox()

// The Safari() method returns:
// "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_8) AppleWebKit/536.26 (KHTML, like Gecko) Version/6.0 Safari/8536.25".
browser.UserAgent = agents.Safari()

// There are methods for a number of browsers and crawlers. For example
// Opera(), MSIE(), AOL(), GoogleBot(), and many more. You can even choose
// the browser version. This will create:
// "Mozilla/5.0 (Windows NT 6.3; x64) Chrome/35 Safari/537.36".
browser.UserAgent = agents.CreateVersion("chrome", "35")

// Creating your own custom user agent is just as easy. The following code
// generates the user agent:
// "MyBrowser/1.0 (Windows NT 6.1; WOW64; x64)".
agents.Name = "MyBrowser"
agents.Version = "1.0"
agents.OSName = "Windows NT"
agents.OSVersion = "6.1"
agents.Comments = []string{"WOW64", "x64"}
browser.UserAgent = agents.Create()
```
The agents package has an internal database for many different versions of many different browsers. See the [API documentation](https://github.com/headzoo/surf/tree/master/docs) for more information.


### Credits
Surf uses the awesome [goquery](https://github.com/PuerkitoBio/goquery) by Martin Angers, and was written using [Intellij](http://www.jetbrains.com/idea/) and the [golang plugin](http://plugins.jetbrains.com/plugin/5047). API documentation was created using [godocdown](https://github.com/robertkrimen/godocdown) by Robert Krimen.


### Use Cases
* Interacting with sites that do not have public APIs.
* Testing/Stressing your sites.
* Creating a web crawler.


### TODO
* Add user authentication.
* Run JavaScript found in the page?
* Add AttributeDownloadAssets so the browser downloads the images, scripts, stylesheets, etc.
* Write more tests. 
* File uploading in forms.
* Handle checkboxes correctly.
