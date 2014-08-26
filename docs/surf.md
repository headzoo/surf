# surf
--
    import "github.com/headzoo/surf"


## Usage

#### type Browser

```go
type Browser struct {
	*element.Page

	// UserAgent is the User-Agent header value sent with requests.
	UserAgent string

	// Cookies stores cookies for every site visited by the browser.
	Cookies http.CookieJar

	// Bookmarks stores the saved bookmarks.
	Bookmarks jar.BookmarksJar

	// History stores the visited pages.
	History *element.PageStack

	// RequestHeaders are additional headers to send with each request.
	RequestHeaders http.Header
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

Returns a boolean value indicating whether a previous page existed, and was
successfully loaded.

#### func (*Browser) BookmarkPage

```go
func (b *Browser) BookmarkPage(name string) error
```
BookmarkPage saves the page URL in the bookmarks with the given name.

#### func (*Browser) Click

```go
func (b *Browser) Click(expr string) error
```
Click clicks on the page element matched by the given expression.

Currently this is only useful for click on links, which will cause the browser
to load the page pointed at by the link. Future versions of Surf may support
JavaScript and clicking on elements will fire the click event.

#### func (*Browser) Form

```go
func (b *Browser) Form(expr string) (element.Submittable, error)
```
Form returns the form in the current page that matches the given expr.

The expr can be any valid goquery expression, and the "form" tag is implied. The
method can be called using only ".login-form" and the expr is automatically
converted to "form.login-form". A complete expression can still be used, for
instance "div.login form".

#### func (*Browser) Forms

```go
func (b *Browser) Forms() []element.Submittable
```
Forms returns an array of every form in the page.

Returns nil when the page does not contain any forms.

#### func (*Browser) Links

```go
func (b *Browser) Links() []*element.Link
```
Links returns an array of every link found in the page.

#### func (*Browser) Open

```go
func (b *Browser) Open(u string) error
```
Open requests the given URL using the GET method.

#### func (*Browser) OpenBookmark

```go
func (b *Browser) OpenBookmark(name string) error
```
OpenBookmark calls Open() with the URL for the bookmark with the given name.

#### func (*Browser) OpenForm

```go
func (b *Browser) OpenForm(u string, data url.Values) error
```
OpenForm appends the data values to the given URL and sends a GET request.

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
func (b *Browser) SetAttribute(a attrib.Attribute, v bool)
```
SetAttribute sets a browser instruction attribute.

#### func (*Browser) SiteCookies

```go
func (b *Browser) SiteCookies() []*http.Cookie
```
SiteCookies returns the cookies for the current site.
