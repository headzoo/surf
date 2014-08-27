// Package element contains types related to web documents
// and document elements.
package element

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/attrib"
	"io"
	"net/http"
	"net/url"
)

// Document represents a web document loaded in a browser.
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

// Page represents a web page document.
type Page struct {
	resp *http.Response
	doc  *goquery.Document
}

// NewPage creates and returns a *Page type.
func NewPage(r *http.Response, d *goquery.Document) *Page {
	return &Page{
		resp: r,
		doc:  d,
	}
}

// Url returns the page URL as a string.
func (p *Page) Url() *url.URL {
	return p.resp.Request.URL
}

// StatusCode returns the response status code.
func (p *Page) StatusCode() int {
	return p.resp.StatusCode
}

// Title returns the page title.
func (p *Page) Title() string {
	return p.doc.Find("title").Text()
}

// Headers returns the page headers.
func (p *Page) Headers() http.Header {
	return p.resp.Header
}

// Body returns the page body as a string of html.
func (p *Page) Body() string {
	body, _ := p.doc.Find("body").Html()
	return body
}

// Dom returns the inner *goquery.Selection.
func (p *Page) Dom() *goquery.Selection {
	return p.doc.First()
}

// PageStack stores Page types in a LIFO stack.
type PageStack struct {
	top  *Element
	size int
}

// Element holds stack values and points to the next element.
type Element struct {
	Value *Page
	Next  *Element
}

// NewPageStack creates and returns a new PageHeap type.
func NewPageStack() *PageStack {
	return &PageStack{}
}

// Len returns the number of pages in the stack.
func (stack *PageStack) Len() int {
	return stack.size
}

// Push adds a new Page at the front of the stack.
func (stack *PageStack) Push(p *Page) int {
	stack.top = &Element{p, stack.top}
	stack.size++
	return stack.size
}

// Pop removes and returns the Page at the front of the stack.
func (stack *PageStack) Pop() *Page {
	if stack.size > 0 {
		value := stack.top.Value
		stack.top = stack.top.Next
		stack.size--
		return value
	}

	return nil
}

// Top returns the Page at the front of the stack without removing it.
func (stack *PageStack) Top() *Page {
	if stack.size == 0 {
		return nil
	}
	return stack.top.Value
}
