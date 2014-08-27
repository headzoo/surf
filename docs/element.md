# element
--
    import "github.com/headzoo/surf/element"


## Usage

#### type Browsable

```go
type Browsable interface {
	Document

	// Open requests the given URL using the GET method.
	Open(url string) error

	// OpenForm appends the data values to the given URL and sends a GET request.
	OpenForm(url string, data url.Values) error

	// OpenBookmark calls Get() with the URL for the bookmark with the given name.
	OpenBookmark(name string) error

	// Post requests the given URL using the POST method.
	Post(url string, bodyType string, body io.Reader) error

	// PostForm requests the given URL using the POST method with the given data.
	PostForm(url string, data url.Values) error

	// Back loads the previously requested page.
	Back() bool

	// Reload duplicates the last successful request.
	Reload() error

	// BookmarkPage saves the page URL in the bookmarks with the given name.
	BookmarkPage(name string) error

	// Click clicks on the page element matched by the given expression.
	Click(expr string) error

	// Form returns the form in the current page that matches the given expr.
	Form(expr string) (Submittable, error)

	// Forms returns an array of every form in the page.
	Forms() []Submittable

	// Links returns an array of every link found in the page.
	Links() []*Link

	// SiteCookies returns the cookies for the current site.
	SiteCookies() []*http.Cookie

	// SetAttribute sets a browser instruction attribute.
	SetAttribute(a attrib.Attribute, v bool)

	// ResolveUrl returns an absolute URL for a possibly relative URL.
	ResolveUrl(u *url.URL) *url.URL

	// ResolveStringUrl works just like ResolveUrl, but the argument and return value are strings.
	ResolveStringUrl(u string) (string, error)

	// Write writes the document to the given writer.
	Write(o io.Writer) (int, error)
}
```

Browsable represents an HTTP web browser.

#### type Document

```go
type Document interface {
	// Url returns the page URL as a string.
	Url() *url.URL

	// StatusCode returns the response status code.
	StatusCode() int

	// Title returns the page title.
	Title() string

	// Headers returns the page headers.
	Headers() http.Header

	// Body returns the page body as a string of html.
	Body() string

	// Dom returns the inner *goquery.Selection.
	Dom() *goquery.Selection
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

#### func (*Form) Dom

```go
func (f *Form) Dom() *goquery.Selection
```
Dom returns the inner *goquery.Selection.

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

#### func (*Form) Submit

```go
func (f *Form) Submit() error
```
Submit submits the form. Clicks the first button in the form, or submits the
form without using any button when the form does not contain any buttons.

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

#### func (*Page) Dom

```go
func (p *Page) Dom() *goquery.Selection
```
Dom returns the inner *goquery.Selection.

#### func (*Page) Headers

```go
func (p *Page) Headers() http.Header
```
Headers returns the page headers.

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

#### type Submittable

```go
type Submittable interface {
	Method() string
	Action() string
	Input(name, value string) error
	Click(button string) error
	Submit() error
	Dom() *goquery.Selection
}
```

Submittable represents an element that may be submitted, such as a form.
