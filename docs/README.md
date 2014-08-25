# gosurf
--
    import "github.com/headzoo/gosurf"


## Usage

```go
const (
	Name    = "GoSurf"
	Version = "0.1"
)
```
The browser name and version.

```go
var (
	DefaultUserAgent                string = fmt.Sprintf("%s/%s (%s; %s)", Name, Version, runtime.Version(), osRelease())
	DefaultAttributeSendReferer     bool   = true
	DefaultAttributeHandleRefresh   bool   = true
	DefaultAttributeFollowRedirects bool   = true
)
```
Default values.

#### type Attribute

```go
type Attribute int
```

Attribute represents a Browser capability.

```go
const (
	// AttributeSendReferer instructs a Browser to send the Referer header.
	AttributeSendReferer Attribute = iota
	// AttributeHandleRefresh instructs a Browser to handle the refresh meta tag.
	AttributeHandleRefresh
	// AttributeFollowRedirects instructs a Browser to follow Location headers.
	AttributeFollowRedirects
)
```

#### type AttributeMap

```go
type AttributeMap map[Attribute]bool
```

AttributeMap represents a map of Attribute values.

#### type Browser

```go
type Browser struct {
	*Page
	UserAgent string
	Jar       http.CookieJar
}
```

Browser is the default Browser implementation.

#### func  NewBrowser

```go
func NewBrowser() *Browser
```
NewBrowser creates and returns a *DefaultBrowser type.

#### func (*Browser) Back

```go
func (b *Browser) Back() bool
```
Back loads the previously requested page.

#### func (*Browser) Cookies

```go
func (b *Browser) Cookies() []*http.Cookie
```
Cookies returns the cookies for the current page.

#### func (*Browser) FollowLink

```go
func (b *Browser) FollowLink(expr string) error
```
FollowLink finds an anchor tag within the current page matching the expr, and
calls Open() on the anchor href attribute value.

#### func (*Browser) Form

```go
func (b *Browser) Form(expr string) (FormElement, error)
```
Form returns the form in the current page that matches the given expr.

#### func (*Browser) Get

```go
func (b *Browser) Get(u string) error
```
Get requests the given URL using the GET method.

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
func NewForm(b WebBrowser, s *goquery.Selection) *Form
```
NewForm creates and returns a *Form type.

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

#### func (*Form) Submit

```go
func (f *Form) Submit() error
```
Submit submits the form. Clicks the first button in the form, or submits the
form without using any button when the form does not contain any buttons.

#### type FormElement

```go
type FormElement interface {
	Input(name, value string) error
	Click(button string) error
	Submit() error
}
```

FormElement represents a single form element from a page.

#### type Page

```go
type Page struct {
}
```

Page represents the attributes of a single web page.

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

#### type WebBrowser

```go
type WebBrowser interface {
	WebPage
	Get(url string) error
	GetForm(url string, data url.Values) error
	Post(url string, bodyType string, body io.Reader) error
	PostForm(url string, data url.Values) error
	FollowLink(expr string) error
	Links() []string
	Form(expr string) (FormElement, error)
	Back() bool
	Reload() error
	Cookies() []*http.Cookie
	SetAttribute(a Attribute, v bool)
	ResolveUrl(u *url.URL) *url.URL
}
```

WebBrowser represents an HTTP browser.

#### type WebPage

```go
type WebPage interface {
	Url() *url.URL
	StatusCode() int
	Title() string
	Headers() http.Header
	Body() string
	Query() *goquery.Document
}
```

WebPage represents a single web page.
