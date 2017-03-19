package surf

import (
	"io"
	"net/http"
	gourl "net/url"
)

// ElementType describes a type of page element, such as an image or stylesheet.
type ElementType uint16

const (
	// ElementTypeLink describes a *Link element.
	ElementTypeLink ElementType = iota

	// ElementTypeImage describes an *Image element.
	ElementTypeImage

	// ElementTypeStylesheet describes a *Stylesheet element.
	ElementTypeStylesheet

	// ElementTypeScript describes a *Script element.
	ElementTypeScript
)

// Element represents a page element, such as an image or stylesheet.
type Element interface {
	// URL returns the asset URL.
	URL() *gourl.URL

	// ID returns the asset ID or an empty string when not available.
	ID() string

	// Type returns the type of element.
	Type() ElementType
}

// BaseElement implements Element.
type BaseElement struct {
	// ID is the value of the id attribute if available.
	id string

	// url of the element.
	url *gourl.URL

	// typ describes the type of element.
	typ ElementType
}

// URL returns the asset URL.
func (at *BaseElement) URL() *gourl.URL {
	return at.url
}

// ID returns the asset ID or an empty string when not available.
func (at *BaseElement) ID() string {
	return at.id
}

// Type returns the asset type.
func (at *BaseElement) Type() ElementType {
	return at.typ
}

// AnchorElement stores the properties of a page link.
type AnchorElement struct {
	BaseElement

	// Text is the text appearing between the opening and closing anchor tag.
	Text string
}

// NewAnchorElement creates and returns a new *AnchorElement instance.
func NewAnchorElement(u *gourl.URL, id, text string) *AnchorElement {
	return &AnchorElement{
		BaseElement: BaseElement{
			url: u,
			id:  id,
			typ: ElementTypeLink,
		},
		Text: text,
	}
}

// ImageElement stores the properties of an image.
type ImageElement struct {
	BaseDownloadableElement

	// Alt is the value of the image alt attribute if available.
	Alt string

	// Title is the value of the image title attribute if available.
	Title string
}

// NewImageElement creates and returns a new *ImageElement instance.
func NewImageElement(url *gourl.URL, id, alt, title string) *ImageElement {
	return &ImageElement{
		BaseDownloadableElement: BaseDownloadableElement{
			BaseElement: BaseElement{
				url: url,
				id:  id,
				typ: ElementTypeImage,
			},
		},
		Alt:   alt,
		Title: title,
	}
}

// StylesheetElement stores the properties of a linked stylesheet.
type StylesheetElement struct {
	BaseDownloadableElement

	// Media is the value of the media attribute. Defaults to "all" when not specified.
	Media string

	// TypeAttr is the value of the type attribute. Defaults to "text/css" when not specified.
	TypeAttr string
}

// NewStylesheetElement creates and returns a new *StylesheetElement instance.
func NewStylesheetElement(url *gourl.URL, id, media, typ string) *StylesheetElement {
	return &StylesheetElement{
		BaseDownloadableElement: BaseDownloadableElement{
			BaseElement: BaseElement{
				url: url,
				typ: ElementTypeStylesheet,
				id:  id,
			},
		},
		Media:    media,
		TypeAttr: typ,
	}
}

// ScriptElement stores the properties of a linked script.
type ScriptElement struct {
	BaseDownloadableElement

	// Type is the value of the type attribute. Defaults to "text/javascript" when not specified.
	TypeAttr string
}

// NewScriptElement creates and returns a new *ScriptElement instance.
func NewScriptElement(url *gourl.URL, id, typ string) *ScriptElement {
	return &ScriptElement{
		BaseDownloadableElement: BaseDownloadableElement{
			BaseElement: BaseElement{
				url: url,
				typ: ElementTypeScript,
				id:  id,
			},
		},
		TypeAttr: typ,
	}
}

// DownloadableElement represents an element that may be downloaded.
type DownloadableElement interface {
	Element

	// Download writes the contents of the element to the given writer.
	//
	// Returns the number of bytes written.
	Download(out io.Writer) (int64, error)

	// DownloadAsync downloads the contents of the element asynchronously.
	//
	// An instance of AsyncDownloadResult will be sent down the given channel
	// when the download is complete.
	DownloadAsync(out io.Writer, ch AsyncDownloadChannel)
}

// BaseDownloadableElement is an element that may be downloaded.
type BaseDownloadableElement struct {
	BaseElement
}

// Download writes the element to the given io.Writer type.
func (at *BaseDownloadableElement) Download(out io.Writer) (int64, error) {
	return DownloadElement(at, out)
}

// DownloadAsync downloads the element asynchronously.
func (at *BaseDownloadableElement) DownloadAsync(out io.Writer, ch AsyncDownloadChannel) {
	DownloadElementAsync(at, out, ch)
}

// AsyncDownloadResult has the results of an asynchronous download.
type AsyncDownloadResult struct {
	// Element is a pointer to the Downloadable asset that was downloaded.
	Element DownloadableElement

	// Writer where the asset data was written.
	Writer io.Writer

	// Size is the number of bytes written to the io.Writer.
	Size int64

	// Error contains any error that occurred during the download or nil.
	Error error
}

// AsyncDownloadChannel is a channel upon which the results of an async download
// are passed.
type AsyncDownloadChannel chan *AsyncDownloadResult

// DownloadElement copies a remote file to the given writer.
func DownloadElement(asset DownloadableElement, out io.Writer) (int64, error) {
	resp, err := http.Get(asset.URL().String())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return io.Copy(out, resp.Body)
}

// DownloadElementAsync downloads an element asynchronously and notifies the given channel
// when the download is complete.
func DownloadElementAsync(asset DownloadableElement, out io.Writer, c AsyncDownloadChannel) {
	go func() {
		results := &AsyncDownloadResult{Element: asset, Writer: out}
		size, err := DownloadElement(asset, out)
		if err != nil {
			results.Error = err
		} else {
			results.Size = size
		}
		c <- results
	}()
}
