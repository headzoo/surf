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


// Follow a link on the page where the link text is "new".
err = browser.FollowLink(":contains('new')")
if err != nil { panic(err) }
fmt.Println(browser.Title())
// Outputs: "newest submissions: reddit.com"


// Login to the site via their login form.
fm := browser.Form("[class='login-form']")
fm.Input("user", "JoeRedditor")
fm.Input("passwd", "d234rlkasd")
err = fm.Submit()
if err != nil { panic(err) }


// Now that we're logged in, follow the link to our profile.
err = browser.FollowLink(":contains('JoeRedditor'))
if err != nil { panic(err) }
fmt.Println(browser.Title())
// Outputs: "overview for JoeRedditor"

// Finally print the page body.
fmt.Println(browser.Body())
```


### Credits
GoSurf uses the awesome [goquery](https://github.com/PuerkitoBio/goquery) by Martin Angers, and was written using [Intellij](http://www.jetbrains.com/idea/) and the [golang plugin](http://plugins.jetbrains.com/plugin/5047).


