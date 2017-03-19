package surf

import (
	"fmt"
	"github.com/headzoo/surf/jar"
	"github.com/jeffail/tunny"
	"io"
	"io/ioutil"
	"net/http"
	gourl "net/url"
	"os"
	"path"
	"sync"
)

// Browser...
type Browser struct {
	*EventTarget

	Document  *Document
	Navigator *Navigator
	Location  *gourl.URL
	Headers   http.Header
	Response  *http.Response
	bookmarks jar.BookmarksJar
	history   jar.History
	state     *jar.State
}

// SendHEAD requests the given URL using the HEAD method.
func (b *Browser) SendHEAD(url string) error {
	return b.sendHEAD(url, "")
}

// SendGET requests the given URL using the GET method.
func (b *Browser) SendGET(url string) error {
	return b.sendGET(url, "")
}

// SendPOST requests the given URL using the POST method.
func (b *Browser) SendPOST(url string, contentType string, body io.Reader) error {
	// @todo
	return nil
}

// SendFormGET appends the data values to the given URL and sends a GET request.
func (b *Browser) SendFormGET(url string, data gourl.Values) error {
	// @todo
	return nil
}

// SendFormPOST requests the given URL using the POST method with the given data.
func (b *Browser) SendFormPOST(url string, data gourl.Values) error {
	// @todo
	return nil
}

// SendMultipartPOST requests the given URL using the POST method with the given data using multipart/form-data format.
func (b *Browser) SendMultipartPOST(u string, fields gourl.Values, files FileSet) error {
	// @todo
	return nil
}

// BookmarkOpen calls SendGET() with the URL for the bookmark with the given name.
func (b *Browser) BookmarkOpen(name string) error {
	// @todo
	return nil
}

// BookmarkSave saves the page URL in the bookmarks with the given name.
func (b *Browser) BookmarkSave(name string) error {
	// @todo
	return nil
}

// Back loads the previously requested page.
func (b *Browser) Back() bool {
	// @todo
	return true
}

// Reload duplicates the last successful request.
func (b *Browser) Reload() error {
	// @todo
	return nil
}

// SavePage the current page and all assets to the given directory.
func (b *Browser) SavePage(dir string, perm os.FileMode) (saveFile string, errs []error) {
	var err error
	errs = []error{}
	if err = os.MkdirAll(dir, perm); err != nil {
		errs = append(errs, err)
		return
	}

	filename := path.Base(b.Location.Path)
	if filename == "." {
		filename = "index.html"
	}
	saveFile = path.Join(dir, filename)
	assetsSubDir := fmt.Sprintf("%s_assets", filename)
	assetsFullDir := path.Join(dir, assetsSubDir)
	if err = os.MkdirAll(assetsFullDir, perm); err != nil {
		errs = append(errs, err)
		return
	}

	assets := b.Document.findDownloadableAssets()
	if len(assets) > 0 {
		wg := sync.WaitGroup{}
		pool, _ := tunny.CreatePool(NumDownloadWorkers, func(object interface{}) interface{} {
			defer wg.Done()
			err := b.downloadAsset(object.(downloadable), assetsFullDir, assetsSubDir, perm)
			if err != nil {
				errs = append(errs, err)
			}
			return nil
		}).Open()
		defer pool.Close()
		for _, a := range assets {
			wg.Add(1)
			go pool.SendWork(a)
		}
		wg.Wait()
	}

	html := b.Document.InnerHTML()
	if err = ioutil.WriteFile(saveFile, []byte(html), perm); err != nil {
		errs = append(errs, err)
	}
	return
}

// History returns the browser history.
// See https://developer.mozilla.org/en-US/docs/Web/API/Window/history
func (b *Browser) History() jar.History {
	return b.history
}

// sendHEAD makes an HTTP HEAD request for the given URL.
// When via is not nil, and AttributeSendReferer is true, the Referer header will
// be set to ref.
func (b *Browser) sendHEAD(url, referer string) error {
	req, err := b.buildRequest(MethodHEAD, url, referer, nil)
	if err != nil {
		return err
	}
	return b.sendRequest(req)
}

// sendGET makes an HTTP GET request for the given URL.
// When via is not nil, and AttributeSendReferer is true, the Referer header will
// be set to ref.
func (b *Browser) sendGET(url, referer string) error {
	req, err := b.buildRequest(MethodGET, url, referer, nil)
	if err != nil {
		return err
	}
	return b.sendRequest(req)
}

// sendPOST makes an HTTP POST request for the given URL.
// When via is not nil, and AttributeSendReferer is true, the Referer header will
// be set to ref.
func (b *Browser) sendPOST(url, referer, contentType string, body io.Reader) error {
	req, err := b.buildRequest(MethodPOST, url, referer, body)
	if err != nil {
		return err
	}
	req.Header.Set(HeaderContentType, contentType)
	return b.sendRequest(req)
}

// sendRequest uses the given *http.Request to make an HTTP request.
func (b *Browser) sendRequest(req *http.Request) error {
	var err error
	if b.Document != nil {
		b.Document.unload()
	}

	b.DispatchEvent(OnRequest, b, NewEventArgs(EventArgValues{
		"request": req,
	}))
	debugRequest(req)
	if b.Response, err = b.buildClient().Do(req); err != nil {
		b.DispatchEvent(OnError, b, &EventArgs{Error: err})
		return err
	}
	debugResponse(b.Response)
	b.DispatchEvent(OnResponse, b, NewEventArgs(EventArgValues{
		"response": b.Response,
	}))

	b.Location = req.URL
	b.history.Push(b.state)
	b.state = jar.NewHistoryState(req, b.Response)
	if req.Method != MethodHEAD {
		if err = b.Document.load(req, b.Response); err != nil {
			return err
		}
	}

	return nil
}

// buildRequest creates and returns a *http.Request type.
// Sets any headers that need to be sent with the request.
func (b *Browser) buildRequest(method, url, referer string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		b.DispatchEvent(OnError, b, &EventArgs{Error: err})
		return nil, err
	}
	req.Header = b.makeRequestHeaders()
	if host := req.Header.Get(HeaderHost); host != "" {
		req.Host = host
	}
	req.Header.Set(HeaderUserAgent, UserAgent)
	if SendReferer && referer != "" {
		req.Header.Set(HeaderReferer, referer)
	}
	return req, nil
}

// buildClient creates, configures, and returns a *http.Client type.
func (b *Browser) buildClient() *http.Client {
	client := &http.Client{}
	client.Jar = JarCookies
	client.CheckRedirect = b.clientCheckRedirect
	if Transport != nil {
		client.Transport = Transport
	}
	return client
}

// clientCheckRedirect is used as the value to http.Client.CheckRedirect.
func (b *Browser) clientCheckRedirect(req *http.Request, _ []*http.Request) error {
	if FollowRedirects {
		return nil
	}
	return fmt.Errorf("Redirects are disabled. Cannot follow '%s'.", req.URL.String())
}

// makeRequestHeaders creates and returns a copy of the default headers and browser headers.
func (b *Browser) makeRequestHeaders() http.Header {
	mk := make(http.Header, len(RequestHeaders)+len(b.Headers))
	for k, v := range RequestHeaders {
		mk[k] = v
	}
	for k, v := range b.Headers {
		mk[k] = v
	}
	return mk
}

// resolveUrl returns an absolute URL for a possibly relative URL.
func (b *Browser) resolveUrl(u *gourl.URL) *gourl.URL {
	// @todo
	return &gourl.URL{}
}

// downloadAsset downloads the given asset and saves it to saveDir.
func (b *Browser) downloadAsset(a downloadable, saveDir, relativePath string, perm os.FileMode) error {
	if src, ok := a.node.Attr(a.urlAttr); ok {
		srcUrl, err := gourl.Parse(src)
		if err != nil {
			return err
		}
		srcUrl = b.Location.ResolveReference(srcUrl)
		imgFile := path.Join(saveDir, srcUrl.Path)
		if imgFile != "" {
			debugMessage(`Downloading "%s" to "%s"`, srcUrl.String(), imgFile)
			if err = os.MkdirAll(path.Dir(imgFile), perm); err != nil {
				return err
			}
			var file *os.File
			if file, err = os.OpenFile(imgFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm); err != nil {
				return err
			}
			defer file.Close()
			resp, err := http.Get(srcUrl.String())
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if _, err = io.Copy(file, resp.Body); err != nil {
				return err
			}
			a.node.SetAttr(a.urlAttr, path.Join(relativePath, srcUrl.Path))
		}
	}
	return nil
}
