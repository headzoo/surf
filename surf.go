package surf

import (
	"fmt"
	"github.com/headzoo/surf/agent"
	"github.com/headzoo/surf/jar"
	"net/http"
	"net/http/httputil"
	"os"
)

const (
	Name    = "Surf"
	Version = "2.0"
)

const (
	DefaultSendReferer         = true
	DefaultMetaRefreshHandling = true
	DefaultFollowRedirects     = true
	DefaultNumDownloadWorkers  = 4
)

const (
	HeaderUserAgent   = "User-Agent"
	HeaderHost        = "Host"
	HeaderReferer     = "Referer"
	HeaderContentType = "Content-Type"
)

const (
	MethodGET  = "GET"
	MethodPOST = "POST"
	MethodHEAD = "HEAD"
)

const (
	ContentTypeTextPlain = "text/plain"
	ContentTypeTextHtml  = "text/html"
)

var (
	// Debugging turns debugging messages on and off.
	Debugging bool = false

	// UserAgent is the User-Agent header value sent with requests.
	UserAgent string = agent.Create(Name, Version)

	// JarState is the current browser state.
	JarState *jar.State = &jar.State{}

	// JarCookies stores cookies for every site visited by the browser.
	JarCookies http.CookieJar = jar.NewMemoryCookies()

	// JarBookmarks stores the saved bookmarks.
	JarBookmarks jar.BookmarksJar = jar.NewMemoryBookmarks()

	// JarHistory stores the visited pages.
	JarHistory jar.History = jar.NewMemoryHistory()

	// transport specifies the mechanism by which individual HTTP
	// requests are made.
	Transport http.RoundTripper

	// RequestHeaders are additional headers to send with each request.
	RequestHeaders http.Header = jar.NewMemoryHeaders()

	// NumDownloadWorkers is the number of workers to download page assets.
	NumDownloadWorkers int = DefaultNumDownloadWorkers

	// SendReferer instructs a Browser to send the Referer header.
	SendReferer bool = DefaultSendReferer

	// MetaRefreshHandling instructs a Browser to handle the refresh meta tag.
	MetaRefreshHandling bool = DefaultMetaRefreshHandling

	// FollowRedirects instructs a Browser to follow Location headers.
	FollowRedirects bool = DefaultFollowRedirects
)

// init the package.
func init() {
	Debugging = os.Getenv("SURF_DEBUG") != ""
}

// NewBrowser returns a new *Browser instance.
func NewBrowser() *Browser {
	b := &Browser{
		EventTarget: NewEventTarget(),
		Headers:     http.Header{},
		bookmarks:   JarBookmarks,
		history:     JarHistory,
		state:       JarState,
		Response:    &http.Response{},
	}
	b.Document = NewDocument(b)
	return b
}

// debugMessage prints a message when debugging is turned on.
func debugMessage(f string, v ...interface{}) {
	if Debugging {
		fmt.Fprintln(os.Stderr, "===== [SURF] =====")
		fmt.Fprintf(os.Stderr, f+"\n", v...)
		fmt.Fprintln(os.Stderr, "==================")
	}
}

// debugRequest prints the req data when debugging is turned on.
func debugRequest(req *http.Request) {
	if Debugging {
		d, _ := httputil.DumpRequest(req, false)
		debugMessage("REQUEST\n%s", d)
	}
}

// debugResponse prints the res data when debugging is turned on.
func debugResponse(res *http.Response) {
	if Debugging {
		d, _ := httputil.DumpResponse(res, false)
		debugMessage("RESPONSE\n%s", d)
	}
}
