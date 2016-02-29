// Package surf ensembles other packages into a usable browser.
package surf

import (
	"github.com/emgfc/surf/browser"
)

// NewBrowser creates and returns a *browser.Browser type.
func NewBrowser() *browser.Browser {
	return browser.NewBrowser()
}
