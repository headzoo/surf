package gosurf

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
)

// WebPage represents a single web page.
type WebPage interface {
	Url() *url.URL
	StatusCode() int
	Title() string
	Headers() http.Header
	Body() string
	Query() *goquery.Document
}

// Page represents the attributes of a single web page.
type Page struct {
	resp *http.Response
	doc *goquery.Document
}

// NewPage creates and returns a *Page type.
func NewPage(r *http.Response, d *goquery.Document) *Page {
	return &Page{
		resp: r,
		doc: d,
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

// Query returns the inner *goquery.Document.
func (p *Page) Query() *goquery.Document {
	return p.doc
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
