// Package surf ensembles other packages into a usable browser.
package surf

import (
	"github.com/headzoo/surf/browser"
	"github.com/headzoo/surf/jar"
	"net/http"
	"net/http/cookiejar"
)

// NewBrowser creates and returns a *browser.Browser type.
func NewBrowser() (*browser.Browser, error) {
	cookies, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	bow := &browser.Browser{}
	bow.SetUserAgent(browser.DefaultUserAgent)
	bow.SetCookieJar(cookies)
	bow.SetBookmarksJar(jar.NewMemoryBookmarks())
	bow.SetHistoryJar(jar.NewMemoryHistory())
	bow.SetHeaders(make(http.Header, 10))
	bow.SetAttributes(browser.AttributeMap{
		browser.SendReferer:         browser.DefaultSendReferer,
		browser.MetaRefreshHandling: browser.DefaultMetaRefreshHandling,
		browser.FollowRedirects:     browser.DefaultFollowRedirects,
	})

	return bow, nil
}
