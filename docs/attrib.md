# attrib
--
    import "github.com/headzoo/surf/attrib"


## Usage

```go
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
```

#### type Attribute

```go
type Attribute int
```

Attribute represents a Browser capability.

```go
const (
	// SendRefererAttribute instructs a Browser to send the Referer header.
	SendReferer Attribute = iota

	// MetaRefreshHandlingAttribute instructs a Browser to handle the refresh meta tag.
	MetaRefreshHandling

	// FollowRedirectsAttribute instructs a Browser to follow Location headers.
	FollowRedirects
)
```

#### type AttributeMap

```go
type AttributeMap map[Attribute]bool
```

AttributeMap represents a map of Attribute values.
