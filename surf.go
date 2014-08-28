// Package surf ensembles other packages into a usable browser.
package surf

import (
	"github.com/headzoo/surf/browser"
	"github.com/headzoo/surf/jar"
)

// NewBrowser creates and returns a *browser.Browser type.
func NewBrowser() *browser.Browser {
	bow := &browser.Browser{}
	bow.SetUserAgent(browser.DefaultUserAgent)
	bow.SetCookieJar(jar.NewMemoryCookies())
	bow.SetBookmarksJar(jar.NewMemoryBookmarks())
	bow.SetHistoryJar(jar.NewMemoryHistory())
	bow.SetHeaders(jar.NewMemoryHeaders())
	bow.SetAttributes(browser.AttributeMap{
		browser.SendReferer:         browser.DefaultSendReferer,
		browser.MetaRefreshHandling: browser.DefaultMetaRefreshHandling,
		browser.FollowRedirects:     browser.DefaultFollowRedirects,
	})

	return bow
}
