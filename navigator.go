package surf

// Navigator represents the state and the identity of the user agent.
// See https://developer.mozilla.org/en-US/docs/Web/API/Navigator
type Navigator struct {
}

// NewNavigator returns a *Navigator instance.
func NewNavigator() *Navigator {
	return &Navigator{}
}

// AppName returns the name of the browser.
// See https://developer.mozilla.org/en-US/docs/Web/API/NavigatorID/appName
func (n *Navigator) AppName() string {
	return Name
}

// AppVersion returns the browser version.
// See https://developer.mozilla.org/en-US/docs/Web/API/NavigatorID/appVersion
func (n *Navigator) AppVersion() string {
	return Version
}

// UserAgent returns the browser user agent.
// See https://developer.mozilla.org/en-US/docs/Web/API/NavigatorID/userAgent
func (n *Navigator) UserAgent() string {
	return UserAgent
}
