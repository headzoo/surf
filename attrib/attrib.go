package attrib

import "github.com/headzoo/surf/agent"

// Attribute represents a Browser capability.
type Attribute int

// AttributeMap represents a map of Attribute values.
type AttributeMap map[Attribute]bool

const (
	// SendRefererAttribute instructs a Browser to send the Referer header.
	SendReferer Attribute = iota

	// MetaRefreshHandlingAttribute instructs a Browser to handle the refresh meta tag.
	MetaRefreshHandling

	// FollowRedirectsAttribute instructs a Browser to follow Location headers.
	FollowRedirects
)

var (
	// DefaultUserAgent is the global user agent value.
	DefaultUserAgent = agent.Create()

	// DefaultSendRefererAttribute is the global value for the AttributeSendReferer attribute.
	DefaultSendReferer = true

	// DefaultMetaRefreshHandlingAttribute is the global value for the AttributeHandleRefresh attribute.
	DefaultMetaRefreshHandling = true

	// DefaultFollowRedirectsAttribute is the global value for the AttributeFollowRedirects attribute.
	DefaultFollowRedirects = true
)
