package surf

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/errors"
	"github.com/headzoo/surf/jars"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	// Name is used as the browser name in the default user agent.
	Name = "Surf"
	// Version is used as the version in the default user agent.
	Version = "0.4.2"
)

// Attribute represents a Browser capability.
type Attribute int

// AttributeMap represents a map of Attribute values.
type AttributeMap map[Attribute]bool

const (
	// SendRefererAttribute instructs a Browser to send the Referer header.
	SendRefererAttribute Attribute = iota
	// MetaRefreshHandlingAttribute instructs a Browser to handle the refresh meta tag.
	MetaRefreshHandlingAttribute
	// FollowRedirectsAttribute instructs a Browser to follow Location headers.
	FollowRedirectsAttribute
)

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

// exprPrefixesImplied are strings a selection expr may start with, and the tag is implied.
var exprPrefixesImplied = []string{":", ".", "["}

// Browsable represents an HTTP web browser.
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
	SiteCookies() []*http.Cookie
	SetAttribute(a Attribute, v bool)
	ResolveUrl(u *url.URL) *url.URL
}

// Browser is the default Browser implementation.
type Browser struct {
	*Page
	UserAgent   string
	Cookies   http.CookieJar
	Bookmarks   jars.BookmarksJar
	History     *PageStack
	lastRequest *http.Request
	attributes  AttributeMap
	refresh     *time.Timer
}

// NewBrowser creates and returns a *Browser type.
func NewBrowser() (*Browser, error) {
	cookies, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Browser{
		UserAgent: DefaultUserAgent,
		Cookies: cookies,
		Bookmarks: jars.NewMemoryBookmarks(),
		History:   NewPageStack(),
		attributes: AttributeMap{
			SendRefererAttribute:         DefaultSendRefererAttribute,
			MetaRefreshHandlingAttribute: DefaultMetaRefreshHandlingAttribute,
			FollowRedirectsAttribute:     DefaultFollowRedirectsAttribute,
		},
	}, nil
}

// Get requests the given URL using the GET method.
func (b *Browser) Get(u string) error {
	return b.sendGet(u, nil)
}

// GetForm appends the data values to the given URL and sends a GET request.
func (b *Browser) GetForm(u string, data url.Values) error {
	ul, err := url.Parse(u)
	if err != nil {
		return err
	}
	ul.RawQuery = data.Encode()

	return b.Get(ul.String())
}

// GetBookmark calls Get() with the URL for the bookmark with the given name.
func (b *Browser) GetBookmark(name string) error {
	url, err := b.Bookmarks.Read(name)
	if err != nil {
		return err
	}
	return b.Get(url)
}

// Post requests the given URL using the POST method.
func (b *Browser) Post(u string, bodyType string, body io.Reader) error {
	return b.sendPost(u, bodyType, body, nil)
}

// PostForm requests the given URL using the POST method with the given data.
func (b *Browser) PostForm(u string, data url.Values) error {
	return b.Post(u, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// Bookmark saves the page URL in the bookmarks with the given name.
func (b *Browser) Bookmark(name string) error {
	return b.Bookmarks.Save(name, b.ResolveUrl(b.Page.Url()).String())
}

// FollowLink finds an anchor tag within the current document matching the expr,
// and calls Get() using the anchor href attribute value.
//
// The expr can be any valid goquery expression, and the "a" tag is implied. The
// method can be called using only ":contains('foo')" and the expr is automatically
// converted to "a:contains('foo')". A complete expression can still be used, for
// instance "p.title a.foo".
func (b *Browser) FollowLink(expr string) error {
	sel := b.Page.doc.Find(prefixSelection(expr, "a"))
	if sel.Length() == 0 {
		return errors.NewElementNotFound(
			"Anchor not found matching expr '%s'.", expr)
	}
	if !sel.Is("a") {
		return errors.NewElementNotFound(
			"Expr '%s' does not match an anchor tag.", expr)
	}

	href, ok := sel.Attr("href")
	if !ok {
		return errors.NewLinkNotFound("No link found matching expr %s.", expr)
	}
	hurl, err := url.Parse(href)
	if err != nil {
		return err
	}
	hurl = b.ResolveUrl(hurl)

	return b.sendGet(hurl.String(), b.Page)
}

// Links returns an array of every anchor tag href value found in the current page.
func (b *Browser) Links() []string {
	sel := b.Page.doc.Find("a")
	links := make([]string, 0, sel.Length())
	sel.Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			links = append(links, href)
		}
	})

	return links
}

// SiteCookies returns the cookies for the current site.
func (b *Browser) SiteCookies() []*http.Cookie {
	return b.Cookies.Cookies(b.Page.Url())
}

// Form returns the form in the current page that matches the given expr.
func (b *Browser) Form(expr string) (FormElement, error) {
	sel := b.Page.doc.Find(prefixSelection(expr, "form"))
	if sel.Length() == 0 {
		return nil, errors.NewElementNotFound(
			"Form not found matching expr '%s'.", expr)
	}
	if !sel.Is("form") {
		return nil, errors.NewElementNotFound(
			"Expr '%s' does not match a form tag.", expr)
	}

	return NewForm(b, sel), nil
}

// Forms returns an array of every form in the page.
func (b *Browser) Forms() []FormElement {
	sel := b.Page.doc.Find("form")
	len := sel.Length()
	if len == 0 {
		return nil
	}

	forms := make([]FormElement, len)
	sel.Each(func(_ int, s *goquery.Selection) {
		forms = append(forms, NewForm(b, s))
	})
	return forms
}

// Back loads the previously requested page.
func (b *Browser) Back() bool {
	if b.History.Len() > 0 {
		b.Page = b.History.Pop()
		return true
	}
	return false
}

// Reload duplicates the last successful request.
func (b *Browser) Reload() error {
	if b.lastRequest != nil {
		return b.send(b.lastRequest)
	}
	return errors.NewPageNotLoaded("Cannot reload, the previous request failed.")
}

// SetAttribute sets a browser instruction attribute.
func (b *Browser) SetAttribute(a Attribute, v bool) {
	b.attributes[a] = v
}

// ResolveUrl returns an absolute URL for a possibly relative URL.
func (b *Browser) ResolveUrl(u *url.URL) *url.URL {
	return b.Url().ResolveReference(u)
}

// client creates, configures, and returns a *http.Client type.
func (b *Browser) client() *http.Client {
	client := &http.Client{}
	client.Jar = b.Cookies
	client.CheckRedirect = b.shouldRedirect
	return client
}

// request creates and returns a *http.Request type.
// Sets any headers that need to be sent with the request.
func (b *Browser) request(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header["User-Agent"] = []string{b.UserAgent}

	return req, nil
}

// sendGet makes an HTTP GET request for the given URL.
// When via is not nil, and AttributeSendReferer is true, the Referer header will
// be set to via's URL.
func (b *Browser) sendGet(url string, via *Page) error {
	req, err := b.request("GET", url)
	if err != nil {
		return err
	}
	if b.attributes[SendRefererAttribute] && via != nil {
		req.Header["Referer"] = []string{via.Url().String()}
	}

	return b.send(req)
}

// sendPost makes an HTTP POST request for the given URL.
// When via is not nil, and AttributeSendReferer is true, the Referer header will
// be set to via's URL.
func (b *Browser) sendPost(url string, bodyType string, body io.Reader, via *Page) error {
	req, err := b.request("POST", url)
	if err != nil {
		return err
	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	req.Body = rc
	req.Header["Content-Type"] = []string{bodyType}
	if b.attributes[SendRefererAttribute] && via != nil {
		req.Header["Referer"] = []string{via.Url().String()}
	}

	return b.send(req)
}

// send uses the given *http.Request to make an HTTP request.
func (b *Browser) send(req *http.Request) error {
	b.preSend()
	resp, err := b.client().Do(req)
	if err != nil {
		return err
	}
	body, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	b.lastRequest = req
	b.History.Push(b.Page)
	b.Page = NewPage(resp, body)
	b.postSend()

	return nil
}

// preSend sets browser state before sending a request.
func (b *Browser) preSend() {
	if b.refresh != nil {
		b.refresh.Stop()
	}
}

// postSend sets browser state after sending a request.
func (b *Browser) postSend() {
	if b.attributes[MetaRefreshHandlingAttribute] {
		sel := b.Page.doc.Find("meta[http-equiv='refresh']")
		if sel.Length() > 0 {
			attr, ok := sel.Attr("content")
			if ok {
				dur, err := time.ParseDuration(attr + "s")
				if err == nil {
					b.refresh = time.NewTimer(dur)
					go func() {
						<-b.refresh.C
						b.Reload()
					}()
				}
			}
		}
	}
}

// shouldRedirect is used as the value to http.Client.CheckRedirect.
func (b *Browser) shouldRedirect(req *http.Request, _ []*http.Request) error {
	if b.attributes[FollowRedirectsAttribute] {
		return nil
	}
	return errors.NewLocation(
		"Redirects are disabled. Cannot follow '%s'.", req.URL.String())
}

// prefixSelection prefixes a selection expr with elm when sel is prefixed with
// one of the values from exprPrefixesImplied.
func prefixSelection(sel, elm string) string {
	for _, prefix := range exprPrefixesImplied {
		if strings.HasPrefix(sel, prefix) {
			return elm + sel
		}
	}
	return sel
}

// osRelease returns the name of the OS and it's release version.
func osRelease() string {
	buf := &syscall.Utsname{}
	err := syscall.Uname(buf)
	if err != nil {
		return "0.0"
	}

	return charsToString(buf.Sysname) + "/" + charsToString(buf.Release)
}

// charsToString converts a [65]int8 byte array into a string.
func charsToString(ca [65]int8) string {
	s := make([]byte, len(ca))
	var lens int
	for ; lens < len(ca); lens++ {
		if ca[lens] == 0 {
			break
		}
		s[lens] = uint8(ca[lens])
	}
	return string(s[0:lens])
}
