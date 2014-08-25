# surf
--
    import "github.com/headzoo/surf"


## Usage

```go
const (
	// Name is used as the browser name in the default user agent.
	Name = "Surf"
	// Version is used as the version in the default user agent.
	Version = "0.4.2"
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
	Get(url string) error
	GetForm(url string, data url.Values) error
	GetBookmark(name string) error
	Post(url string, bodyType string, body io.Reader) error
	PostForm(url string, data url.Values) error
	Bookmark(name string) error
	FollowLink(expr string) error
	Links() []string
	Form(expr string) (FormElement, error)
	Forms() []FormElement
	Back() bool
	Reload() error
	Cookies() []*http.Cookie
	SetAttribute(a Attribute, v bool)
	ResolveUrl(u *url.URL) *url.URL
	Stop() error
}
```

Browsable represents an HTTP web browser.

#### type Browser

```go
type Browser struct {
	*Page
	UserAgent string
	CookieJar http.CookieJar
	Bookmarks jars.BookmarksJar
	History   *PageStack
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

#### func (*Browser) Bookmark

```go
func (b *Browser) Bookmark(name string) error
```
Bookmark saves the page URL in the bookmarks with the given name.

#### func (*Browser) Cookies

```go
func (b *Browser) Cookies() []*http.Cookie
```
Cookies returns the cookies for the current page.

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
func (b *Browser) Links() []string
```
Links returns an array of every anchor tag href value found in the current page.

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

#### func (*Browser) Stop

```go
func (b *Browser) Stop() error
```
Stop releases resources held by the browser.

This method is called automatically by the runtime, but is safe to call
repeatedly without any errors.

The browser should not be used after Stop is called. Doing so will cause
unexpected behavior.

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
