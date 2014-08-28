package browser

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/errors"
	"github.com/headzoo/surf/jar"
	"io"
	"io/ioutil"
	"net/http"
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
	SendReferer Attribute = iota

	// MetaRefreshHandlingAttribute instructs a Browser to handle the refresh meta tag.
	MetaRefreshHandling

	// FollowRedirectsAttribute instructs a Browser to follow Location headers.
	FollowRedirects
)

// Browsable represents an HTTP web browser.
type Browsable interface {
	// SetUserAgent sets the user agent.
	SetUserAgent(ua string)

	// SetAttribute sets a browser instruction attribute.
	SetAttribute(a Attribute, v bool)

	// SetAttributes is used to set all the browser attributes.
	SetAttributes(a AttributeMap)

	// SetBookmarksJar sets the bookmarks jar the browser uses.
	SetBookmarksJar(bj jar.BookmarksJar)

	// SetCookieJar is used to set the cookie jar the browser uses.
	SetCookieJar(cj http.CookieJar)

	// SetHistoryJar is used to set the history jar the browser uses.
	SetHistoryJar(hj jar.History)

	// AddHeader adds a header the browser sends with each request.
	AddHeader(name, value string)

	// SetHeaders sets the headers the browser sends with each request.
	SetHeaders(h http.Header)

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

	// Images returns an array of every image found in the page.
	Images() []*Image

	// Stylesheets returns an array of every stylesheet linked to the document.
	Stylesheets() []*Stylesheet

	// Scripts returns an array of every script linked to the document.
	Scripts() []*Script

	// SiteCookies returns the cookies for the current site.
	SiteCookies() []*http.Cookie

	// ResolveUrl returns an absolute URL for a possibly relative URL.
	ResolveUrl(u *url.URL) *url.URL

	// ResolveStringUrl works just like ResolveUrl, but the argument and return value are strings.
	ResolveStringUrl(u string) (string, error)

	// Download writes the contents of the document to the given writer.
	Download(o io.Writer) (int64, error)

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

	// Find returns the dom selections matching the given expression.
	Find(expr string) *goquery.Selection
}

// Default is the default Browser implementation.
type Browser struct {
	// state is the current browser state.
	state *jar.State

	// userAgent is the User-Agent header value sent with requests.
	userAgent string

	// cookies stores cookies for every site visited by the browser.
	cookies http.CookieJar

	// bookmarks stores the saved bookmarks.
	bookmarks jar.BookmarksJar

	// history stores the visited pages.
	history jar.History

	// headers are additional headers to send with each request.
	headers http.Header

	// attributes is the set browser attributes.
	attributes AttributeMap

	// refresh is a timer used to meta refresh pages.
	refresh *time.Timer
}

// Open requests the given URL using the GET method.
func (bow *Browser) Open(u string) error {
	return bow.sendGet(u, "")
}

// OpenForm appends the data values to the given URL and sends a GET request.
func (bow *Browser) OpenForm(u string, data url.Values) error {
	ul, err := url.Parse(u)
	if err != nil {
		return err
	}
	ul.RawQuery = data.Encode()

	return bow.Open(ul.String())
}

// OpenBookmark calls Open() with the URL for the bookmark with the given name.
func (bow *Browser) OpenBookmark(name string) error {
	url, err := bow.bookmarks.Read(name)
	if err != nil {
		return err
	}
	return bow.Open(url)
}

// Post requests the given URL using the POST method.
func (bow *Browser) Post(u string, bodyType string, body io.Reader) error {
	return bow.sendPost(u, bodyType, body, "")
}

// PostForm requests the given URL using the POST method with the given data.
func (bow *Browser) PostForm(u string, data url.Values) error {
	return bow.Post(u, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// Back loads the previously requested page.
//
// Returns a boolean value indicating whether a previous page existed, and was
// successfully loaded.
func (bow *Browser) Back() bool {
	if bow.history.Len() > 0 {
		bow.state = bow.history.Pop()
		return true
	}
	return false
}

// Reload duplicates the last successful request.
func (bow *Browser) Reload() error {
	if bow.state.Request != nil {
		return bow.send(bow.state.Request)
	}
	return errors.NewPageNotLoaded("Cannot reload, the previous request failed.")
}

// BookmarkPage saves the page URL in the bookmarks with the given name.
func (bow *Browser) BookmarkPage(name string) error {
	return bow.bookmarks.Save(name, bow.ResolveUrl(bow.Url()).String())
}

// Click clicks on the page element matched by the given expression.
//
// Currently this is only useful for click on links, which will cause the browser
// to load the page pointed at by the link. Future versions of Surf may support
// JavaScript and clicking on elements will fire the click event.
func (bow *Browser) Click(expr string) error {
	sel := bow.Find(expr)
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
	href, err := bow.ResolveStringUrl(href)
	if err != nil {
		return err
	}

	return bow.sendGet(href, bow.Url().String())
}

// Form returns the form in the current page that matches the given expr.
func (bow *Browser) Form(expr string) (Submittable, error) {
	sel := bow.Find(expr)
	if sel.Length() == 0 {
		return nil, errors.NewElementNotFound(
			"Form not found matching expr '%s'.", expr)
	}
	if !sel.Is("form") {
		return nil, errors.NewElementNotFound(
			"Expr '%s' does not match a form tag.", expr)
	}

	return NewForm(bow, sel), nil
}

// Forms returns an array of every form in the page.
//
// Returns nil when the page does not contain any forms.
func (bow *Browser) Forms() []Submittable {
	sel := bow.Find("form")
	len := sel.Length()
	if len == 0 {
		return nil
	}

	forms := make([]Submittable, len)
	sel.Each(func(_ int, s *goquery.Selection) {
		forms = append(forms, NewForm(bow, s))
	})
	return forms
}

// Links returns an array of every link found in the page.
func (bow *Browser) Links() []*Link {
	links := make([]*Link, 0, InitialAssetsSliceSize)
	bow.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, err := bow.attributeToUrl("href", s)
		if err == nil {
			links = append(links, &Link{
				ID:   bow.attributeOrDefault("id", "", s),
				URL:  href,
				Text: s.Text(),
			})
		}
	})

	return links
}

// Images returns an array of every image found in the page.
func (bow *Browser) Images() []*Image {
	images := make([]*Image, 0, InitialAssetsSliceSize)
	bow.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, err := bow.attributeToUrl("src", s)
		if err == nil {
			images = append(images, &Image{
				ID:    bow.attributeOrDefault("id", "", s),
				URL:   src,
				Alt:   bow.attributeOrDefault("alt", "", s),
				Title: bow.attributeOrDefault("title", "", s),
			})
		}
	})

	return images
}

// Stylesheets returns an array of every stylesheet linked to the document.
func (bow *Browser) Stylesheets() []*Stylesheet {
	stylesheets := make([]*Stylesheet, 0, InitialAssetsSliceSize)
	bow.Find("link").Each(func(_ int, s *goquery.Selection) {
		rel, ok := s.Attr("rel")
		if ok && rel == "stylesheet" {
			href, err := bow.attributeToUrl("href", s)
			if err == nil {
				stylesheets = append(stylesheets, &Stylesheet{
					ID:    bow.attributeOrDefault("id", "", s),
					URL:   href,
					Media: bow.attributeOrDefault("media", "all", s),
					Type:  bow.attributeOrDefault("type", "text/css", s),
				})
			}
		}
	})

	return stylesheets
}

// Scripts returns an array of every script linked to the document.
func (bow *Browser) Scripts() []*Script {
	scripts := make([]*Script, 0, InitialAssetsSliceSize)
	bow.Find("script").Each(func(_ int, s *goquery.Selection) {
		src, err := bow.attributeToUrl("src", s)
		if err == nil {
			scripts = append(scripts, &Script{
				ID:   bow.attributeOrDefault("id", "", s),
				URL:  src,
				Type: bow.attributeOrDefault("type", "text/javascript", s),
			})
		}
	})

	return scripts
}

// SiteCookies returns the cookies for the current site.
func (bow *Browser) SiteCookies() []*http.Cookie {
	return bow.cookies.Cookies(bow.Url())
}

// SetCookieJar is used to set the cookie jar the browser uses.
func (bow *Browser) SetCookieJar(cj http.CookieJar) {
	bow.cookies = cj
}

// SetUserAgent sets the user agent.
func (bow *Browser) SetUserAgent(userAgent string) {
	bow.userAgent = userAgent
}

// SetAttribute sets a browser instruction attribute.
func (bow *Browser) SetAttribute(a Attribute, v bool) {
	bow.attributes[a] = v
}

// SetAttributes is used to set all the browser attributes.
func (bow *Browser) SetAttributes(a AttributeMap) {
	bow.attributes = a
}

// SetBookmarksJar sets the bookmarks jar the browser uses.
func (bow *Browser) SetBookmarksJar(bj jar.BookmarksJar) {
	bow.bookmarks = bj
}

// SetHistoryJar is used to set the history jar the browser uses.
func (bow *Browser) SetHistoryJar(hj jar.History) {
	bow.history = hj
}

// AddHeader sets a header the browser sends with each request.
func (bow *Browser) AddHeader(name, value string) {
	bow.headers.Add(name, value)
}

// SetHeaders sets the headers the browser sends with each request.
func (bow *Browser) SetHeaders(h http.Header) {
	bow.headers = h
}

// ResolveUrl returns an absolute URL for a possibly relative URL.
func (bow *Browser) ResolveUrl(u *url.URL) *url.URL {
	return bow.Url().ResolveReference(u)
}

// ResolveStringUrl works just like ResolveUrl, but the argument and return value are strings.
func (bow *Browser) ResolveStringUrl(u string) (string, error) {
	pu, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	pu = bow.Url().ResolveReference(pu)
	return pu.String(), nil
}

// Download writes the contents of the document to the given writer.
func (bow *Browser) Download(o io.Writer) (int64, error) {
	h, err := bow.state.Dom.Html()
	if err != nil {
		return 0, err
	}
	l, err := io.WriteString(o, h)
	return int64(l), err
}

// Url returns the page URL as a string.
func (bow *Browser) Url() *url.URL {
	return bow.state.Request.URL
}

// StatusCode returns the response status code.
func (bow *Browser) StatusCode() int {
	return bow.state.Response.StatusCode
}

// Title returns the page title.
func (bow *Browser) Title() string {
	return bow.state.Dom.Find("title").Text()
}

// Headers returns the page headers.
func (bow *Browser) Headers() http.Header {
	return bow.state.Response.Header
}

// Body returns the page body as a string of html.
func (bow *Browser) Body() string {
	body, _ := bow.state.Dom.Find("body").Html()
	return body
}

// Dom returns the inner *goquery.Selection.
func (bow *Browser) Dom() *goquery.Selection {
	return bow.state.Dom.First()
}

// Find returns the dom selections matching the given expression.
func (bow *Browser) Find(expr string) *goquery.Selection {
	return bow.state.Dom.Find(expr)
}

// client creates, configures, and returns a *http.Client type.
func (bow *Browser) client() *http.Client {
	client := &http.Client{}
	client.Jar = bow.cookies
	client.CheckRedirect = bow.shouldRedirect
	return client
}

// request creates and returns a *http.Request type.
// Sets any headers that need to be sent with the request.
func (bow *Browser) request(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = bow.headers
	req.Header["User-Agent"] = []string{bow.userAgent}
	return req, nil
}

// sendGet makes an HTTP GET request for the given URL.
// When via is not nil, and AttributeSendReferer is true, the Referer header will
// be set to via's URL.
func (bow *Browser) sendGet(url string, via string) error {
	req, err := bow.request("GET", url)
	if err != nil {
		return err
	}
	if bow.attributes[SendReferer] && via != "" {
		req.Header["Referer"] = []string{via}
	}

	return bow.send(req)
}

// sendPost makes an HTTP POST request for the given URL.
// When via is not nil, and AttributeSendReferer is true, the Referer header will
// be set to via's URL.
func (bow *Browser) sendPost(url string, bodyType string, body io.Reader, via string) error {
	req, err := bow.request("POST", url)
	if err != nil {
		return err
	}
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	req.Body = rc
	req.Header["Content-Type"] = []string{bodyType}
	if bow.attributes[SendReferer] && via != "" {
		req.Header["Referer"] = []string{via}
	}

	return bow.send(req)
}

// send uses the given *http.Request to make an HTTP request.
func (bow *Browser) send(req *http.Request) error {
	bow.preSend()
	resp, err := bow.client().Do(req)
	if err != nil {
		return err
	}
	body, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}
	bow.history.Push(bow.state)
	bow.state = &jar.State{
		Request:  req,
		Response: resp,
		Dom:      body,
	}
	bow.postSend()

	return nil
}

// preSend sets browser state before sending a request.
func (bow *Browser) preSend() {
	if bow.refresh != nil {
		bow.refresh.Stop()
	}
}

// postSend sets browser state after sending a request.
func (bow *Browser) postSend() {
	if bow.attributes[MetaRefreshHandling] {
		sel := bow.Find("meta[http-equiv='refresh']")
		if sel.Length() > 0 {
			attr, ok := sel.Attr("content")
			if ok {
				dur, err := time.ParseDuration(attr + "s")
				if err == nil {
					bow.refresh = time.NewTimer(dur)
					go func() {
						<-bow.refresh.C
						bow.Reload()
					}()
				}
			}
		}
	}
}

// shouldRedirect is used as the value to http.Client.CheckRedirect.
func (bow *Browser) shouldRedirect(req *http.Request, _ []*http.Request) error {
	if bow.attributes[FollowRedirects] {
		return nil
	}
	return errors.NewLocation(
		"Redirects are disabled. Cannot follow '%s'.", req.URL.String())
}

// attributeToUrl reads an attribute from an element and returns a url.
func (bow *Browser) attributeToUrl(name string, sel *goquery.Selection) (*url.URL, error) {
	src, ok := sel.Attr(name)
	if !ok {
		return nil, errors.NewAttributeNotFound(
			"Attribute '%s' not found.", name)
	}
	ur, err := url.Parse(src)
	if err != nil {
		return nil, err
	}

	return bow.ResolveUrl(ur), nil
}

// attributeOrDefault reads an attribute and returns it or the default value when it's empty.
func (bow *Browser) attributeOrDefault(name, def string, sel *goquery.Selection) string {
	a, ok := sel.Attr(name)
	if ok {
		return a
	}
	return def
}
