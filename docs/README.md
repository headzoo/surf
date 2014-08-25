# surf
--
    import "github.com/headzoo/surf"


## Usage

```go
const (
	// Name is used as the browser name in the default user agent.
	Name = "Surf"
	// Version is used as the version in the default user agent.
	Version = "0.4.3"
)
```

```go
var (
	// DefaultUserAgent is the global user agent value.
	DefaultUserAgent string = fmt.Sprintf("%s/%s (%s; %s)", Name, Version, runtime.Version(), osRelease())
	// DefaultSendRefererAttribute is the global value for the AttributeSendReferer attribute.
	DefaultSendRefererAttribute bool = true
	// DefaultMetaRefreshHandlingAttribute is the global value for the AttributeHandleRefresh attribute.
	DefaultMetaRefreshHandlingAttribute bool = true
	// DefaultFollowRedirectsAttribute is the global value for the AttributeFollowRedirects attribute.
	DefaultFollowRedirectsAttribute bool = true
)
```

#### type Attribute

```go
type Attribute int
```

Attribute represents a Browser capability.

```go
const (
	// SendRefererAttribute instructs a Browser to send the Referer header.
	SendRefererAttribute Attribute = iota
	// MetaRefreshHandlingAttribute instructs a Browser to handle the refresh meta tag.
	MetaRefreshHandlingAttribute
	// FollowRedirectsAttribute instructs a Browser to follow Location headers.
	FollowRedirectsAttribute
)
```

#### type AttributeMap

```go
type AttributeMap map[Attribute]bool
```

AttributeMap represents a map of Attribute values.

#### type Browsable

```go
type Browsable interface {
	Document
	// Get requests the given URL using the GET method.
	Get(url string) error
	// GetForm appends the data values to the given URL and sends a GET request.
	GetForm(url string, data url.Values) error
	// GetBookmark calls Get() with the URL for the bookmark with the given name.
	GetBookmark(name string) error
	// Post requests the given URL using the POST method.
	Post(url string, bodyType string, body io.Reader) error
	// PostForm requests the given URL using the POST method with the given data.
	PostForm(url string, data url.Values) error
	// BookmarkPage saves the page URL in the bookmarks with the given name.
	BookmarkPage(name string) error
	// FollowLink finds an anchor tag within the current document matching the expr,
	// and calls Get() using the anchor href attribute value.
	FollowLink(expr string) error
	// Links returns an array of every link found in the page.
	Links() []*Link
	// Form returns the form in the current page that matches the given expr.
	Form(expr string) (FormElement, error)
	// Forms returns an array of every form in the page.
	Forms() []FormElement
	// Back loads the previously requested page.
	Back() bool
	// Reload duplicates the last successful request.
	Reload() error
	// SiteCookies returns the cookies for the current site.
	SiteCookies() []*http.Cookie
	// SetAttribute sets a browser instruction attribute.
	SetAttribute(a Attribute, v bool)
	// ResolveUrl returns an absolute URL for a possibly relative URL.
	ResolveUrl(u *url.URL) *url.URL
	// ResolveStringUrl works just like ResolveUrl, but the argument and return value are strings.
	ResolveStringUrl(u string) (string, error)
}
```

Browsable represents an HTTP web browser.

#### type Browser

```go
type Browser struct {
	*Page
	// UserAgent is the User-Agent header value sent with requests.
	UserAgent string
	// Cookies stores cookies for every site visited by the browser.
	Cookies http.CookieJar
	// Bookmarks stores the saved bookmarks.
	Bookmarks jars.BookmarksJar
	// History stores the visited pages.
	History *PageStack
}
```

Browser is the default Browser implementation.

#### func  NewBrowser

```go
func NewBrowser() (*Browser, error)
```
NewBrowser creates and returns a *Browser type.

#### func (*Browser) Back

```go
func (b *Browser) Back() bool
```
Back loads the previously requested page.

#### func (*Browser) BookmarkPage

```go
func (b *Browser) BookmarkPage(name string) error
```
BookmarkPage saves the page URL in the bookmarks with the given name.

#### func (*Browser) FollowLink

```go
func (b *Browser) FollowLink(expr string) error
```
FollowLink finds an anchor tag within the current document matching the expr,
and calls Get() using the anchor href attribute value.

The expr can be any valid goquery expression, and the "a" tag is implied. The
method can be called using only ":contains('foo')" and the expr is automatically
converted to "a:contains('foo')". A complete expression can still be used, for
instance "p.title a.foo".

#### func (*Browser) Form

```go
func (b *Browser) Form(expr string) (FormElement, error)
```
Form returns the form in the current page that matches the given expr.

The expr can be any valid goquery expression, and the "form" tag is implied. The
method can be called using only ".login-form" and the expr is automatically
converted to "form.login-form". A complete expression can still be used, for
instance "div.login form".

#### func (*Browser) Forms

```go
func (b *Browser) Forms() []FormElement
```
Forms returns an array of every form in the page.

#### func (*Browser) Get

```go
func (b *Browser) Get(u string) error
```
Get requests the given URL using the GET method.

#### func (*Browser) GetBookmark

```go
func (b *Browser) GetBookmark(name string) error
```
GetBookmark calls Get() with the URL for the bookmark with the given name.

#### func (*Browser) GetForm

```go
func (b *Browser) GetForm(u string, data url.Values) error
```
GetForm appends the data values to the given URL and sends a GET request.

#### func (*Browser) Links

```go
func (b *Browser) Links() []*Link
```
Links returns an array of every link found in the page.

#### func (*Browser) Post

```go
func (b *Browser) Post(u string, bodyType string, body io.Reader) error
```
Post requests the given URL using the POST method.

#### func (*Browser) PostForm

```go
func (b *Browser) PostForm(u string, data url.Values) error
```
PostForm requests the given URL using the POST method with the given data.

#### func (*Browser) Reload

```go
func (b *Browser) Reload() error
```
Reload duplicates the last successful request.

#### func (*Browser) ResolveStringUrl

```go
func (b *Browser) ResolveStringUrl(u string) (string, error)
```
ResolveStringUrl works just like ResolveUrl, but the argument and return value
are strings.

#### func (*Browser) ResolveUrl

```go
func (b *Browser) ResolveUrl(u *url.URL) *url.URL
```
ResolveUrl returns an absolute URL for a possibly relative URL.

#### func (*Browser) SetAttribute

```go
func (b *Browser) SetAttribute(a Attribute, v bool)
```
SetAttribute sets a browser instruction attribute.

#### func (*Browser) SiteCookies

```go
func (b *Browser) SiteCookies() []*http.Cookie
```
SiteCookies returns the cookies for the current site.

#### type Document

```go
type Document interface {
	Url() *url.URL
	StatusCode() int
	Title() string
	Headers() http.Header
	Body() string
	Query() *goquery.Document
}
```

Document represents a web document loaded in a browser.

#### type Element

```go
type Element struct {
	Value *Page
	Next  *Element
}
```

Element holds stack values and points to the next element.

#### type Form

```go
type Form struct {
}
```

Form is the default form element.

#### func  NewForm

```go
func NewForm(b Browsable, s *goquery.Selection) *Form
```
NewForm creates and returns a *Form type.

#### func (*Form) Action

```go
func (f *Form) Action() string
```
Action returns the form action URL. The URL will always be absolute.

#### func (*Form) Click

```go
func (f *Form) Click(button string) error
```
Click submits the form by clicking the button with the given name.

#### func (*Form) Input

```go
func (f *Form) Input(name, value string) error
```
Input sets the value of a form field.

#### func (*Form) Method

```go
func (f *Form) Method() string
```
Method returns the form method, eg "GET" or "POST".

#### func (*Form) Query

```go
func (f *Form) Query() *goquery.Selection
```
Query returns the inner *goquery.Selection.

#### func (*Form) Submit

```go
func (f *Form) Submit() error
```
Submit submits the form. Clicks the first button in the form, or submits the
form without using any button when the form does not contain any buttons.

#### type FormElement

```go
type FormElement interface {
	Method() string
	Action() string
	Input(name, value string) error
	Click(button string) error
	Submit() error
	Query() *goquery.Selection
}
```

FormElement represents a single form element from a page.

#### type Link

```go
type Link struct {
	// ID is the value of the id attribute or empty when there is no id.
	ID string
	// Href is the value of the href attribute.
	Href string
	// Text is the text appearing between the opening and closing anchor tag.
	Text string
}
```

Link stores the properties of a page link.

#### type Page

```go
type Page struct {
}
```

Page represents a web page document.

#### func  NewPage

```go
func NewPage(r *http.Response, d *goquery.Document) *Page
```
NewPage creates and returns a *Page type.

#### func (*Page) Body

```go
func (p *Page) Body() string
```
Body returns the page body as a string of html.

#### func (*Page) Headers

```go
func (p *Page) Headers() http.Header
```
Headers returns the page headers.

#### func (*Page) Query

```go
func (p *Page) Query() *goquery.Document
```
Query returns the inner *goquery.Document.

#### func (*Page) StatusCode

```go
func (p *Page) StatusCode() int
```
StatusCode returns the response status code.

#### func (*Page) Title

```go
func (p *Page) Title() string
```
Title returns the page title.

#### func (*Page) Url

```go
func (p *Page) Url() *url.URL
```
Url returns the page URL as a string.

#### type PageStack

```go
type PageStack struct {
}
```

PageStack stores Page types in a LIFO stack.

#### func  NewPageStack

```go
func NewPageStack() *PageStack
```
NewPageStack creates and returns a new PageHeap type.

#### func (*PageStack) Len

```go
func (stack *PageStack) Len() int
```
Len returns the number of pages in the stack.

#### func (*PageStack) Pop

```go
func (stack *PageStack) Pop() *Page
```
Pop removes and returns the Page at the front of the stack.

#### func (*PageStack) Push

```go
func (stack *PageStack) Push(p *Page) int
```
Push adds a new Page at the front of the stack.

#### func (*PageStack) Top

```go
func (stack *PageStack) Top() *Page
```
Top returns the Page at the front of the stack without removing it.
