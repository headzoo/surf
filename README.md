GoSurf
======
Stateful programmatic web browsing in Go, modeled after John J. Lee's Python library [mechanize](https://github.com/jjlee/mechanize).


### Installation
Download the library using go.  
`go get github.com/headzoo/gosurf`

Import the library into your project.  
`import "github.com/headzoo/gosurf"`


### Usage
```go
// Start by creating a new browser.
browser := gosurf.NewBrowser()

// Requesting a page.
err := browser.Get("http://www.reddit.com")
if err != nil { panic(err) }
fmt.Println(browser.Title())
// Outputs: "reddit: the front page of the internet"


// Follow a link on the page where the link text is "new". GoSurf uses the selector
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
See the [API documentation](https://github.com/headzoo/gosurf/tree/master/docs) for more information.


### Settings
```go
browser := gosurf.NewBrowser()

// Override the default user agent.
browser.UserAgent = "MyBrowser"

// Attributes control how the browser behaves.
browser.SetAttribute(gosurf.AttributeSendReferer, false)
browser.SetAttribute(gosurf.AttributeHandleRefresh, false)
browser.SetAttribute(gosurf.AttributeRollowRedirects, false)

// Override the build in cookie jar.
jar, err := cookiejar.New(nil)
if err != nil { panic(err) }
browser.Jar = jar
```
See the [API documentation](https://github.com/headzoo/gosurf/tree/master/docs) for more information.


### Credits
GoSurf uses the awesome [goquery](https://github.com/PuerkitoBio/goquery) by Martin Angers, and was written using [Intellij](http://www.jetbrains.com/idea/) and the [golang plugin](http://plugins.jetbrains.com/plugin/5047). API documentation was created using [godocdown](https://github.com/robertkrimen/godocdown) by Robert Krimen.


### TODO
* Add user authentication.
* Write more tests. 
* File uploading in forms.
* Handle checkboxes correctly.
