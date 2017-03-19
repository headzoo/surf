package surf

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	gourl "net/url"
	"regexp"
	"strings"
	"time"
)

// downloadable stores the details of a node which points to downloadable content.
type downloadable struct {
	node    *goquery.Selection
	urlAttr string
}

// contentTypeRegexp pulls values out of the content type header.
var contentTypeRegexp = regexp.MustCompile(`([\w\-]+/[\w\-]+)(\s*;\s*charset=([\w\-]+))?`)

// Document stores the details of the current browser document.
type Document struct {
	*EventTarget

	// Location stores the current url.
	Location *gourl.URL

	// dom is the actual content of the document.
	dom *goquery.Document

	// browser owns the document.
	browser *Browser

	// contentType comes from the content-type header.
	contentType string

	// charSet is the document character set.
	charSet string

	// refreshTimer is a timer used to meta refresh pages.
	refreshTimer *time.Timer
}

// NewDocument returns a *Document instance.
func NewDocument(b *Browser) *Document {
	return &Document{
		EventTarget: NewEventTarget(),
		browser:     b,
		contentType: ContentTypeTextPlain,
	}
}

// Title returns the page title.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/title
func (doc *Document) Title() string {
	return strings.TrimSpace(doc.dom.Find("title").Text())
}

// Content type returns the document content type, e.g. "text/html".
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/contentType
func (doc *Document) ContentType() string {
	return doc.contentType
}

// CharacterSet returns the document character set, e.g. "utf-8".
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/characterSet
func (doc *Document) CharacterSet() string {
	return doc.charSet
}

// InnerHTML returns the document html.
func (doc *Document) InnerHTML() string {
	html, _ := doc.dom.Html()
	return html
}

// Body returns the page body as a string of html.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/body
func (doc *Document) Body() *goquery.Selection {
	return doc.dom.Find("body").First()
}

// See https://developer.mozilla.org/en-US/docs/Web/API/Document/head
func (doc *Document) Head() *goquery.Selection {
	return doc.dom.Find("head").First()
}

// Cookie returns the cookies for the document.
func (doc *Document) Cookie() []*http.Cookie {
	// @todo
	return []*http.Cookie{}
}

// Anchors returns an array of every anchor tag found in the page.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/anchors
func (doc *Document) Anchors() []*AnchorElement {
	// @todo
	return []*AnchorElement{}
}

// Images returns an array of every image found in the page.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/images
func (doc *Document) Images() []*ImageElement {
	images := make([]*ImageElement, 0, 20)
	doc.dom.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, err := doc.attrToResolvedURL("src", s)
		if err == nil {
			images = append(images, NewImageElement(
				src,
				s.AttrOr("id", ""),
				s.AttrOr("alt", ""),
				s.AttrOr("title", ""),
			))
		}
	})

	return images
}

// Stylesheets returns an array of every stylesheet linked to the document.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/styleSheets
func (doc *Document) Stylesheets() []*StylesheetElement {
	// @todo
	return []*StylesheetElement{}
}

// Scripts returns an array of every script linked to the document.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/scripts
func (doc *Document) Scripts() []*ScriptElement {
	// @todo
	return []*ScriptElement{}
}

// Forms returns an array of every form in the page.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/forms
func (doc *Document) Forms() []Submittable {
	// @todo
	return []Submittable{}
}

// GetElementById returns a reference to the element by its ID.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/getElementById
func (doc *Document) GetElementByID(id string) *goquery.Selection {
	return doc.dom.Find(fmt.Sprintf(`#%s`, id)).First()
}

// GetElementsByClassName returns every element with the given class.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/getElementsByClassName
func (doc *Document) GetElementsByClassName(className string) *goquery.Selection {
	return doc.dom.Find(fmt.Sprintf(`.%s`, className))
}

// GetElementsByName returns the elements with the given name.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/getElementsByName
func (doc *Document) GetElementsByName(name string) *goquery.Selection {
	return doc.dom.Find(fmt.Sprintf(`[%s]`, name))
}

// GetElementsByTagName returns the elements with the given tag.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/getElementsByTagName
func (doc *Document) GetElementsByTagName(name string) *goquery.Selection {
	return doc.dom.Find(name)
}

// QuerySelector returns the first element matching the given selector.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/querySelector
func (doc *Document) QuerySelector(selector string) *goquery.Selection {
	return doc.dom.Find(selector).First()
}

// QuerySelectorAll returns all of the elements matching the given selector.
// See https://developer.mozilla.org/en-US/docs/Web/API/Document/querySelectorAll
func (doc *Document) QuerySelectorAll(selector string) *goquery.Selection {
	return doc.dom.Find(selector)
}

// Click clicks on the page element matched by the given expression.
func (doc *Document) Click(expr string) error {
	// @todo
	return nil
}

// Form returns the form in the current page that matches the given expr.
func (doc *Document) Form(expr string) (Submittable, error) {
	// @todo
	return nil, nil
}

// Write writes the contents of the document to the given writer.
func (doc *Document) Write(w io.Writer) (int64, error) {
	// @todo
	return 0, nil
}

// load is called to parse the document.
func (doc *Document) load(req *http.Request, res *http.Response) (err error) {
	var body []byte
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		doc.DispatchEvent(OnError, doc, &EventArgs{Error: err})
		return
	}
	buff := bytes.NewBuffer(body)
	if doc.dom, err = goquery.NewDocumentFromReader(buff); err != nil {
		doc.DispatchEvent(OnError, doc, &EventArgs{Error: err})
		return
	}

	doc.Location = req.URL
	ct := res.Header.Get(HeaderContentType)
	matches := contentTypeRegexp.FindStringSubmatch(ct)
	if len(matches) > 0 {
		doc.contentType = matches[1]
		if len(matches) > 2 {
			doc.charSet = matches[3]
		}
	}
	if doc.charSet == "" {
		sel := doc.dom.Find(`meta[charset]`)
		if sel.Length() > 0 {
			if attr, ok := sel.Attr("charset"); ok {
				doc.charSet = attr
			}
		}
	}

	if doc.contentType == ContentTypeTextHtml && MetaRefreshHandling {
		sel := doc.dom.Find(`meta[http-equiv="refresh"]`)
		if sel.Length() > 0 {
			if attr, ok := sel.Attr("content"); ok {
				dur, err := time.ParseDuration(attr + "s")
				if err == nil {
					doc.refreshTimer = time.NewTimer(dur)
					go func() {
						<-doc.refreshTimer.C
						doc.browser.Reload()
					}()
				}
			}
		}
	}

	doc.DispatchEvent(OnLoad, doc, nil)
	return
}

// unload is called before the document is destroyed.
func (doc *Document) unload() {
	doc.DispatchEvent(OnUnload, doc, nil)
	doc.dom = nil
	if doc.refreshTimer != nil {
		doc.refreshTimer.Stop()
	}
}

// findDownloadableAssets returns document assets which may be downloaded.
func (doc *Document) findDownloadableAssets() []downloadable {
	assets := []downloadable{}
	host := doc.Location.Host
	find := map[string]string{
		`img`:                    "src",
		`script`:                 "src",
		`link[rel="stylesheet"]`: "href",
		`link[rel="icon"]`:       "href",
	}
	for selector, urlAttr := range find {
		doc.dom.Find(selector).Each(func(_ int, node *goquery.Selection) {
			if attr := node.AttrOr(urlAttr, ""); attr != "" {
				if url, err := gourl.Parse(attr); err == nil {
					url = doc.Location.ResolveReference(url)
					if url.Host == host {
						assets = append(assets, downloadable{
							node:    node,
							urlAttr: urlAttr,
						})
					}
				}
			}
		})
	}

	return assets
}

// attributeToUrl reads an attribute from an element and returns a url.
func (doc *Document) attrToResolvedURL(name string, sel *goquery.Selection) (*gourl.URL, error) {
	src, ok := sel.Attr(name)
	if !ok {
		return nil, fmt.Errorf("Attribute '%s' not found.", name)
	}
	ur, err := gourl.Parse(src)
	if err != nil {
		return nil, err
	}
	return doc.Location.ResolveReference(ur), nil
}
