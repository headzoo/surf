# jars
--
    import "github.com/headzoo/surf/jars"


## Usage

#### type BookmarksJar

```go
type BookmarksJar interface {
	Save(name, url string) error
	Read(name string) (string, error)
	Remove(name string) bool
	Has(name string) bool
}
```

BookmarksJar is a container for storage and retrieval of bookmarks.

#### type MemoryBookmarks

```go
type MemoryBookmarks struct {
}
```

MemoryBookmarks is an in-memory implementation of BookmarksJar.

#### func  NewMemoryBookmarks

```go
func NewMemoryBookmarks() *MemoryBookmarks
```
NewMemoryBookmarks creates and returns a new *BookmarkMemoryJar type.

#### func (*MemoryBookmarks) Has

```go
func (b *MemoryBookmarks) Has(name string) bool
```
Has returns a boolean value indicating whether a bookmark exists with the given
name.

#### func (*MemoryBookmarks) Read

```go
func (b *MemoryBookmarks) Read(name string) (string, error)
```
Read returns the URL for the bookmark with the given name.

Returns an error when a bookmark does not exist with the given name. Use the
Has() method first to avoid errors.

#### func (*MemoryBookmarks) Remove

```go
func (b *MemoryBookmarks) Remove(name string) bool
```
Remove deletes the bookmark with the given name.

Returns a boolean value indicating whether a bookmark existed with the given
name and was removed. This method may be safely called even when a bookmark with
the given name does not exist.

#### func (*MemoryBookmarks) Save

```go
func (b *MemoryBookmarks) Save(name, url string) error
```
Save saves a bookmark with the given name.

Returns an error when a bookmark with the given name already exists. Use the
Has() or Remove() methods first to avoid errors.
