package browser

import (
	"github.com/headzoo/surf/util"
	"io"
)

// Link stores the properties of a page link.
type Link struct {
	// ID is the value of the id attribute if available.
	ID string

	// Href is the value of the href attribute.
	Href string

	// Text is the text appearing between the opening and closing anchor tag.
	Text string
}

// Image stores the properties of an image.
type Image struct {
	// ID is the value of the id attribute if available.
	ID string

	// Src is the value of the image src attribute.
	Src string

	// Alt is the value of the image alt attribute if available.
	Alt string

	// Title is the value of the image title attribute if available.
	Title string
}

// Download writes the image to the given io.Writer type.
func (i *Image) Download(out io.Writer) (int64, error) {
	return util.Download(i.Src, out)
}
