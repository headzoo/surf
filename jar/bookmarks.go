package jar

import "github.com/headzoo/surf/errors"

// BookmarksJar is a container for storage and retrieval of bookmarks.
type BookmarksJar interface {
	Save(name, url string) error
	Read(name string) (string, error)
	Remove(name string) bool
	Has(name string) bool
}

// MemoryBookmarks is an in-memory implementation of BookmarksJar.
type MemoryBookmarks struct {
	bookmarks map[string]string
}

// NewMemoryBookmarks creates and returns a new *BookmarkMemoryJar type.
func NewMemoryBookmarks() *MemoryBookmarks {
	return &MemoryBookmarks{
		bookmarks: make(map[string]string, 20),
	}
}

// Save saves a bookmark with the given name.
//
// Returns an error when a bookmark with the given name already exists. Use the
// Has() or Remove() methods first to avoid errors.
func (b *MemoryBookmarks) Save(name, url string) error {
	if b.Has(name) {
		return errors.New(
			"Bookmark with the name '%s' already exists.", name)
	}
	b.bookmarks[name] = url
	return nil
}

// Read returns the URL for the bookmark with the given name.
//
// Returns an error when a bookmark does not exist with the given name. Use the
// Has() method first to avoid errors.
func (b *MemoryBookmarks) Read(name string) (string, error) {
	if !b.Has(name) {
		return "", errors.New(
			"A bookmark does not exist with the name '%s'.", name)
	}
	return b.bookmarks[name], nil
}

// Remove deletes the bookmark with the given name.
//
// Returns a boolean value indicating whether a bookmark existed with the given
// name and was removed. This method may be safely called even when a bookmark
// with the given name does not exist.
func (b *MemoryBookmarks) Remove(name string) bool {
	if b.Has(name) {
		delete(b.bookmarks, name)
		return true
	}
	return false
}

// Has returns a boolean value indicating whether a bookmark exists with the given name.
func (b *MemoryBookmarks) Has(name string) bool {
	_, ok := b.bookmarks[name]
	return ok
}
