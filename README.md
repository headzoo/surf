Surf
====
Stateful programmatic web browsing in Go, modeled after John J. Lee's Python library [mechanize](https://github.com/jjlee/mechanize).


### Installation
Download the library using go.  
`go get github.com/headzoo/surf`

Import the library into your project.  
`import "github.com/headzoo/surf"`


### Usage
```go
// Start by creating a new browser.
browser := surf.NewBrowser()

// Requesting a page.
err := browser.Get("http://www.reddit.com")
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
// values from the body. Lets print the titles for each submission on the
// reddit home page.
browser.Query().Find("a.title").Each(func(_ int, s *goquery.Selection) {
    fmt.Println(s.Text())
})
```
See the [API documentation](https://github.com/headzoo/surf/tree/master/docs) for more information.


### Settings
```go
browser := surf.NewBrowser()

// Override the default user agent.
browser.UserAgent = "MyBrowser"

// Set the user agent globally. Each Browser instance you create will use this.
surf.DefaultUserAgent = "MyBrowser"


// Attributes control how the browser behaves.
browser.SetAttribute(surf.AttributeSendReferer, false)
browser.SetAttribute(surf.AttributeHandleRefresh, false)
browser.SetAttribute(surf.AttributeRollowRedirects, false)

// The attributes may also be set globally.
surf.DefaultAttributeSendReferer = false
surf.DefaultAttributeHandleRefresh = false
surf.DefaultAttributeFollowRedirects = false


// Override the build in cookie jar.
jar, err := cookiejar.New(nil)
if err != nil { panic(err) }
browser.Jar = jar
```
See the [API documentation](https://github.com/headzoo/surf/tree/master/docs) for more information.


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
