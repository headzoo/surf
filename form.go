package surf

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/errors"
	"net/url"
	"strings"
)

// FormElement represents a single form element from a page.
type FormElement interface {
	Method() string
	Action() string
	Input(name, value string) error
	Click(button string) error
	Submit() error
	Query() *goquery.Selection
}

// Form is the default form element.
type Form struct {
	browser   Browsable
	selection *goquery.Selection
	method    string
	action    string
	fields    url.Values
	buttons   url.Values
}

// NewForm creates and returns a *Form type.
func NewForm(b Browsable, s *goquery.Selection) *Form {
	fields, buttons := serializeForm(s)
	method, action := formAttributes(b, s)

	return &Form{
		browser:   b,
		selection: s,
		method:    method,
		action:    action,
		fields:    fields,
		buttons:   buttons,
	}
}

// Method returns the form method, eg "GET" or "POST".
func (f *Form) Method() string {
	return f.method
}

// Action returns the form action URL.
// The URL will always be absolute.
func (f *Form) Action() string {
	return f.action
}

// Input sets the value of a form field.
func (f *Form) Input(name, value string) error {
	if _, ok := f.fields[name]; ok {
		f.fields.Set(name, value)
		return nil
	}
	return errors.NewElementNotFound(
		"No input found with name '%s'.", name)
}

// Submit submits the form.
// Clicks the first button in the form, or submits the form without using
// any button when the form does not contain any buttons.
func (f *Form) Submit() error {
	if len(f.buttons) > 0 {
		for name := range f.buttons {
			return f.Click(name)
		}
	}
	return f.send("", "")
}

// Click submits the form by clicking the button with the given name.
func (f *Form) Click(button string) error {
	if _, ok := f.buttons[button]; !ok {
		return errors.NewInvalidFormValue(
			"Form does not contain a button with the name '%s'.", button)
	}
	return f.send(button, f.buttons[button][0])
}

// Query returns the inner *goquery.Selection.
func (f *Form) Query() *goquery.Selection {
	return f.selection
}

// send submits the form.
func (f *Form) send(buttonName, buttonValue string) error {
	method, ok := f.selection.Attr("method")
	if !ok {
		method = "GET"
	}
	action, ok := f.selection.Attr("action")
	if !ok {
		action = f.browser.Url().String()
	}
	aurl, err := url.Parse(action)
	if err != nil {
		return err
	}
	aurl = f.browser.ResolveUrl(aurl)

	values := make(url.Values, len(f.fields)+1)
	for name, vals := range f.fields {
		values[name] = vals
	}
	if buttonName != "" {
		values.Set(buttonName, buttonValue)
	}

	if strings.ToUpper(method) == "GET" {
		return f.browser.GetForm(aurl.String(), values)
	} else {
		return f.browser.PostForm(aurl.String(), values)
	}
}

// Serialize converts the form fields into a url.Values type.
// Returns two url.Value types. The first is the form field values, and the
// second is the form button values.
func serializeForm(sel *goquery.Selection) (url.Values, url.Values) {
	input := sel.Find("input,button")
	if input.Length() == 0 {
		return url.Values{}, url.Values{}
	}

	fields := make(url.Values)
	buttons := make(url.Values)
	input.Each(func(_ int, s *goquery.Selection) {
		name, ok := s.Attr("name")
		if ok {
			typ, ok := s.Attr("type")
			if ok {
				if typ == "submit" {
					val, ok := s.Attr("value")
					if ok {
						buttons.Add(name, val)
					} else {
						buttons.Add(name, "")
					}
				} else {
					val, ok := s.Attr("value")
					if ok {
						fields.Add(name, val)
					}
				}
			}
		}
	})

	return fields, buttons
}

func formAttributes(b Browsable, s *goquery.Selection) (string, string) {
	method, ok := s.Attr("method")
	if !ok {
		method = "GET"
	}
	action, ok := s.Attr("action")
	if !ok {
		action = b.Url().String()
	}
	aurl, err := url.Parse(action)
	if err != nil {
		return "", ""
	}
	aurl = b.ResolveUrl(aurl)

	return strings.ToUpper(method), aurl.String()
}
