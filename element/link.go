package element

// Link stores the properties of a page link.
type Link struct {
	// ID is the value of the id attribute or empty when there is no id.
	ID string

	// Href is the value of the href attribute.
	Href string

	// Text is the text appearing between the opening and closing anchor tag.
	Text string
}
