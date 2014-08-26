package surf

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/agents"
	"github.com/headzoo/surf/errors"
	"github.com/headzoo/surf/jars"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
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
	DefaultUserAgent = agents.Create()

	// DefaultSendRefererAttribute is the global value for the AttributeSendReferer attribute.
	DefaultSendRefererAttribute = true

	// DefaultMetaRefreshHandlingAttribute is the global value for the AttributeHandleRefresh attribute.
	DefaultMetaRefreshHandlingAttribute = true

	// DefaultFollowRedirectsAttribute is the global value for the AttributeFollowRedirects attribute.
	DefaultFollowRedirectsAttribute = true
)

// exprPrefixesImplied are strings a selection expr may start with, and the tag is implied.
var exprPrefixesImplied = []string{":", ".", "["}

// Link stores the properties of a page link.
type Link struct {
	// ID is the value of the id attribute or empty when there is no id.
	ID string

	// Href is the value of the href attribute.
	Href string

	// Text is the text appearing between the opening and closing anchor tag.
	Text string
}

// Browsable represents an HTTP web browser.
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
	Form(expr string) (FormElement, error)

	// Forms returns an array of every form in the page.
	Forms() []FormElement

	// Links returns an array of every link found in the page.
	Links() []*Link

	// SiteCookies returns the cookies for the current site.
	SiteCookies() []*http.Cookie

	// SetAttribute sets a browser instruction attribute.
	SetAttribute(a Attribute, v bool)

	// ResolveUrl returns an absolute URL for a possibly relative URL.
	ResolveUrl(u *url.URL) *url.URL

	// ResolveStringUrl works just like ResolveUrl, but the argument and return value are strings.
	ResolveStringUrl(u string) (string, error)
}

// Browser is the default Browser implementation.
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

	// RequestHeaders are additional headers to send with each request.
	RequestHeaders http.Header

	// lastRequest is the *http.Request for the last successful request.
	lastRequest *http.Request

	// attributes is the set browser attributes.
	attributes AttributeMap

	// refresh is a timer used to meta refresh pages.
	refresh *time.Timer
}

// NewBrowser creates and returns a *Browser type.
func NewBrowser() (*Browser, error) {
	cookies, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Browser{
		UserAgent:      DefaultUserAgent,
		Cookies:        cookies,
		Bookmarks:      jars.NewMemoryBookmarks(),
		History:        NewPageStack(),
		RequestHeaders: make(http.Header, 10),
		attributes: AttributeMap{
			SendRefererAttribute:         DefaultSendRefererAttribute,
			MetaRefreshHandlingAttribute: DefaultMetaRefreshHandlingAttribute,
			FollowRedirectsAttribute:     DefaultFollowRedirectsAttribute,
		},
	}, nil
}

// Open requests the given URL using the GET method.
func (b *Browser) Open(u string) error {
	return b.sendGet(u, nil)
}

// OpenForm appends the data values to the given URL and sends a GET request.
func (b *Browser) OpenForm(u string, data url.Values) error {
	ul, err := url.Parse(u)
	if err != nil {
		return err
	}
	ul.RawQuery = data.Encode()

	return b.Open(ul.String())
}

// OpenBookmark calls Open() with the URL for the bookmark with the given name.
func (b *Browser) OpenBookmark(name string) error {
	url, err := b.Bookmarks.Read(name)
	if err != nil {
		return err
	}
	return b.Open(url)
}

// Post requests the given URL using the POST method.
func (b *Browser) Post(u string, bodyType string, body io.Reader) error {
	return b.sendPost(u, bodyType, body, nil)
}

// PostForm requests the given URL using the POST method with the given data.
func (b *Browser) PostForm(u string, data url.Values) error {
	return b.Post(u, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
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

// BookmarkPage saves the page URL in the bookmarks with the given name.
func (b *Browser) BookmarkPage(name string) error {
	return b.Bookmarks.Save(name, b.ResolveUrl(b.Page.Url()).String())
}

// Click clicks on the page element matched by the given expression.
//
// Currently this is only useful for click on links, which will cause the browser
// to load the page pointed at by the link. Future versions of Surf may support
// JavaScript and clicking on elements will fire the click event.
func (b *Browser) Click(expr string) error {
	sel := b.Page.doc.Find(prefixSelection(expr, "a"))
	if sel.Length() == 0 {
		return errors.NewElementNotFound(
			"Element not found matching expr '%s'.", expr)
	}
	if !sel.Is("a") {
		return errors.NewElementNotFound(
			"Expr '%s' must match an anchor tag.", expr)
	}

	href, ok := sel.Attr("href")
	if !ok {
		return errors.NewLinkNotFound(
			"No link found matching expr '%s'.", expr)
	}
	href, err := b.ResolveStringUrl(href)
	if err != nil {
		return err
	}

	return b.sendGet(href, b.Page)
}

// Form returns the form in the current page that matches the given expr.
//
// The expr can be any valid goquery expression, and the "form" tag is implied. The
// method can be called using only ".login-form" and the expr is automatically
// converted to "form.login-form". A complete expression can still be used, for
// instance "div.login form".
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

// Links returns an array of every link found in the page.
func (b *Browser) Links() []*Link {
	sel := b.Page.doc.Find("a")
	links := make([]*Link, 0, sel.Length())

	sel.Each(func(_ int, s *goquery.Selection) {
		id, _ := s.Attr("id")
		href, ok := s.Attr("href")
		if ok {
			href, err := b.ResolveStringUrl(href)
			if err == nil {
				links = append(links, &Link{
					ID:   id,
					Href: href,
					Text: s.Text(),
				})
			}
		}
	})

	return links
}

// SiteCookies returns the cookies for the current site.
func (b *Browser) SiteCookies() []*http.Cookie {
	return b.Cookies.Cookies(b.Page.Url())
}

// SetAttribute sets a browser instruction attribute.
func (b *Browser) SetAttribute(a Attribute, v bool) {
	b.attributes[a] = v
}

// ResolveUrl returns an absolute URL for a possibly relative URL.
func (b *Browser) ResolveUrl(u *url.URL) *url.URL {
	return b.Url().ResolveReference(u)
}

// ResolveStringUrl works just like ResolveUrl, but the argument and return value are strings.
func (b *Browser) ResolveStringUrl(u string) (string, error) {
	pu, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	pu = b.Url().ResolveReference(pu)
	return pu.String(), nil
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
	req.Header = b.RequestHeaders
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
