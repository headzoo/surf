package browser

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

// InitialDownloadBufferSize is the initial buffer size for downloaded assets.
var InitialDownloadBufferSize = 100000

// InitialAssetsArraySize is the initial size when allocating a slice.
var InitialAssetsSliceSize = 20

// AssetType describes a type of page asset, such as an image or stylesheet.
type AssetType uint16

const (
	// ImageAsset describes an *Image asset.
	ImageAsset AssetType = iota

	// StylesheetAsset describes a *Stylesheet asset.
	StylesheetAsset

	// ScriptAsset describes a *Script asset.
	ScriptAsset
)

// AsyncDownloadResult has the results of an asynchronous download.
type AsyncDownloadResult struct {
	// Asset is a pointer to the downloaded asset, such as an *Image or
	// *Stylesheet.
	Asset interface{}

	// Type is the type of asset that was downloaded, such as ImageAsset
	// or StylesheetAsset.
	Type AssetType

	// Data is the downloaded data such as the image data.
	Data []byte

	// Error contains any error that occurred during the download or nil.
	Error error
}

// AsyncDownloadChannel is a channel upon which the results of an async download
// are passed.
type AsyncDownloadChannel chan *AsyncDownloadResult

// Downloadable represents an asset that may be downloaded.
type Downloadable interface {
	// Download writes the contents of the element to the given writer.
	//
	// Returns the number of bytes written.
	Download(out io.Writer) (int64, error)

	// DownloadAsync downloads the contents of the element asynchronously.
	//
	// An instance of AsyncDownloadResult will be sent down the given channel
	// when the download is complete.
	DownloadAsync(ch AsyncDownloadChannel)
}

// Link stores the properties of a page link.
type Link struct {
	// ID is the value of the id attribute if available.
	ID string

	// URL is the asset URL.
	URL *url.URL

	// Text is the text appearing between the opening and closing anchor tag.
	Text string
}

// Image stores the properties of an image.
type Image struct {
	// ID is the value of the id attribute if available.
	ID string

	// URL is the asset URL.
	URL *url.URL

	// Alt is the value of the image alt attribute if available.
	Alt string

	// Title is the value of the image title attribute if available.
	Title string
}

// Download writes the image to the given io.Writer type.
func (i *Image) Download(out io.Writer) (int64, error) {
	return downloadAsset(i.URL, out)
}

// DownloadAsync downloads the image asynchronously.
func (i *Image) DownloadAsync(ch AsyncDownloadChannel) {
	downloadAssetAsync(i.URL, i, ImageAsset, ch)
}

// Stylesheet stores the properties of a linked stylesheet.
type Stylesheet struct {
	// ID is the value of the id attribute if available.
	ID string

	// URL is the asset URL.
	URL *url.URL

	// Media is the value of the media attribute. Defaults to "all" when not specified.
	Media string

	// Type is the value of the type attribute. Defaults to "text/css" when not specified.
	Type string
}

// Download writes the stylesheet to the given io.Writer type.
func (s *Stylesheet) Download(out io.Writer) (int64, error) {
	return downloadAsset(s.URL, out)
}

// DownloadAsync downloads the stylesheet asynchronously.
func (s *Stylesheet) DownloadAsync(ch AsyncDownloadChannel) {
	downloadAssetAsync(s.URL, s, StylesheetAsset, ch)
}

// Script stores the properties of a linked script.
type Script struct {
	// ID is the value of the id attribute if available.
	ID string

	// URL is the asset URL.
	URL *url.URL

	// Type is the value of the type attribute. Defaults to "text/javascript" when not specified.
	Type string
}

// Download writes the script to the given io.Writer type.
func (s *Script) Download(out io.Writer) (int64, error) {
	return downloadAsset(s.URL, out)
}

// DownloadAsync downloads the stylesheet asynchronously.
func (s *Script) DownloadAsync(ch AsyncDownloadChannel) {
	downloadAssetAsync(s.URL, s, ScriptAsset, ch)
}

// downloadAsset copies a remote file to the given writer.
func downloadAsset(u *url.URL, out io.Writer) (int64, error) {
	resp, err := http.Get(u.String())
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return io.Copy(out, resp.Body)
}

// downloadAssetAsync downloads an asset asynchronously and notifies the given channel
// when the download is complete.
func downloadAssetAsync(u *url.URL, asset interface{}, typ AssetType, c AsyncDownloadChannel) {
	go func() {
		init := make([]byte, 0, InitialDownloadBufferSize)
		buff := bytes.NewBuffer(init)
		results := &AsyncDownloadResult{Asset: asset, Type: typ}

		_, err := downloadAsset(u, buff)
		if err != nil {
			results.Error = err
		} else {
			results.Data = buff.Bytes()
		}
		c <- results
	}()
}
