package browser

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/base64"

	"net/url"

	"strings"

	"io/ioutil"

	surferrors "github.com/headzoo/surf/errors"
	"github.com/headzoo/surf/jar"
	"github.com/headzoo/ut"
)

func newBrowser() *Browser {
	return &Browser{
		headers: make(http.Header, 10),
		history: jar.NewMemoryHistory(),
	}
}

func TestBrowserFormDefaultNotSelected(t *testing.T) {
	ts := setupTestServer(`
<!doctype html>
<html>

<head>
	<title>Echo Form</title>
</head>

<body>
	<form method="post" name="default">
		<input type="text" name="age" value="" />
		<input type="radio" name="gender" value="male" />
		<input type="radio" name="gender" value="female" />
		<input type="checkbox" name="option1" value="on" />
		<input type="checkbox" name="option2" value="on" />
		<select name="count">
            <option></option>
            <option value="1">One</option>
            <option value="2">Two</option>
            <option value="3">Three</option>
            <option value="4">Four</option>
            <option value="5">Five</option>
        </select>
		<input type="submit" name="submit1" value="submitted1" />
		<input type="submit" name="submit2" value="submitted2" />
	</form>
</body>

</html>`, t)
	defer ts.Close()

	bow := newBrowser()
	err := bow.Open(ts.URL)
	ut.AssertNil(err)

	f, err := bow.Form("[name='default']")
	ut.AssertNil(err)

	ut.AssertEquals(false, f.(*Form).selects["count"].multiple)

	// Initial state should not have any radio or checkbox inputs selected
	// submit with second button
	err = f.Click("submit2")
	ut.AssertEquals("age=&submit2=submitted2", string(bow.body))

	// Change text intput for age
	// submit with first button
	err = f.Input("age", "55")
	ut.AssertNil(err)
	err = f.Click("submit1")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&submit1=submitted1`, string(bow.body))

	// gender does not exist in the form, so Set() is required to add it to the form
	err = f.Set("gender", "male")
	ut.AssertNil(err)
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&gender=male&submit2=submitted2`, string(bow.body))

	// Change gender
	err = f.Input("gender", "female")
	ut.AssertNil(err)
	err = f.Click("submit1")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&gender=female&submit1=submitted1`, string(bow.body))

	err = f.Set("option1", "on")
	ut.AssertNil(err)
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&gender=female&option1=on&submit2=submitted2`, string(bow.body))

	err = f.Set("option2", "on")
	ut.AssertNil(err)
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&gender=female&option1=on&option2=on&submit2=submitted2`, string(bow.body))

	// uncheck option1
	f.Remove("option1")
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&gender=female&option2=on&submit2=submitted2`, string(bow.body))

	// uncheck option2
	f.Remove("option2")
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&gender=female&submit2=submitted2`, string(bow.body))

	err = f.Check("option1")
	ut.AssertNil(err)
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&gender=female&option1=on&submit2=submitted2`, string(bow.body))
	b, err := f.IsChecked("option1")
	ut.AssertNil(err)
	ut.AssertEquals(true, b)

	// option2 does not exist in the form, so Set() is required to add it to the form
	err = f.Check("option2")
	ut.AssertNil(err)
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&gender=female&option1=on&option2=on&submit2=submitted2`, string(bow.body))

	// uncheck option1
	err = f.UnCheck("option1")
	ut.AssertNil(err)
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&gender=female&option2=on&submit2=submitted2`, string(bow.body))
	b, err = f.IsChecked("option1")
	ut.AssertNil(err)
	ut.AssertEquals(false, b)

	// uncheck option2
	err = f.UnCheck("option2")
	ut.AssertNil(err)
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&gender=female&submit2=submitted2`, string(bow.body))
	_, err = f.IsChecked("option3")
	ut.AssertEquals(surferrors.NewElementNotFound("No checkbox found with name 'option3'."), err)

	// select count by label
	err = f.SelectByOptionLabel("count", "Two")
	ut.AssertNil(err)
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&count=2&gender=female&submit2=submitted2`, string(bow.body))

	// select multi count by label
	err = f.SelectByOptionLabel("count", "Two", "Three")
	ut.AssertEquals(surferrors.NewElementNotFound("The select element with name 'count' is not a select miltiple."), err)

	// select count by value
	err = f.SelectByOptionValue("count", "5")
	ut.AssertNil(err)
	err = f.Click("submit2")
	ut.AssertNil(err)
	ut.AssertEquals(`age=55&count=5&gender=female&submit2=submitted2`, string(bow.body))

	// select multi count by value
	err = f.SelectByOptionValue("count", "5", "3")
	ut.AssertEquals(surferrors.NewElementNotFound("The select element with name 'count' is not a select miltiple."), err)
}

func TestBrowserFormDefaultsSelected(t *testing.T) {
	ts := setupTestServer(`
<!doctype html>
<html>

<head>
	<title>Echo Form</title>
</head>

<body>
	<form method="post" name="default">
		<input type="radio" name="gender" value="male" />
		<input type="radio" name="gender" checked="checked" value="female" />
		<input type="checkbox" name="option1" value="on" />
		<input type="checkbox" name="option2" checked="checked" value="on" />
		<select name="count">
            <option></option>
            <option value="1">One</option>
            <option value="2">Two</option>
            <option value="3" Selected="SELECTED">Three</option>
            <option value="4">Four</option>
            <option value="5">Five</option>
            <option value="5,6"> Five &amp; Six </option>
        </select>
		<input type="submit" name="submit" value="submitted" />
	</form>
</body>

</html>`, t)
	defer ts.Close()

	bow := newBrowser()
	err := bow.Open(ts.URL)
	ut.AssertNil(err)

	f, err := bow.Form("[name='default']")
	ut.AssertNil(err)

	ut.AssertEquals(false, f.(*Form).selects["count"].multiple)

	// Initial state should have defaults selected
	err = f.Submit()
	ut.AssertNil(err)
	ut.AssertEquals("count=3&gender=female&option2=on&submit=submitted", string(bow.body))

	val, err := f.Value("option2")
	ut.AssertNil(err)
	ut.AssertEquals("on", val)
}

func TestBrowserFormSelectMultiple(t *testing.T) {
	ts := setupTestServer(`
<!doctype html>
<html>

<head>
	<title>Echo Form</title>
</head>

<body>
	<form method="post" name="default">
		<select name="count" multiple>
            <option></option>
            <option value="1" selected="selected">One</option>
            <option value="2">Two</option>
            <option value="3" selected="selected">Three</option>
            <option value="4">Four</option>
            <option value="5">Five</option>
        </select>
		<input type="submit" name="submit" value="submitted" />
	</form>
</body>

</html>`, t)
	defer ts.Close()

	bow := newBrowser()
	err := bow.Open(ts.URL)
	ut.AssertNil(err)

	f, err := bow.Form("[name='default']")
	ut.AssertNil(err)

	ut.AssertEquals(true, f.(*Form).selects["count"].multiple)

	// Check initial values
	vals, err := f.SelectValues("count")
	ut.AssertNil(err)
	ut.AssertEquals(2, len(vals))
	ut.AssertEquals("1", vals[0])
	ut.AssertEquals("3", vals[1])

	// Check initial values
	vals, err = f.SelectLabels("count")
	ut.AssertNil(err)
	ut.AssertEquals(2, len(vals))
	ut.AssertEquals("One", vals[0])
	ut.AssertEquals("Three", vals[1])

	// Initial state should have defaults selected
	err = f.Submit()
	ut.AssertNil(err)
	ut.AssertEquals("count=1&count=3&submit=submitted", string(bow.body))

	// select multi count by value
	err = f.SelectByOptionValue("count", "5", "1")
	ut.AssertNil(err)
	err = f.Submit()
	ut.AssertNil(err)
	ut.AssertEquals(`count=5&count=1&submit=submitted`, string(bow.body))

	// select multi count by label
	err = f.SelectByOptionLabel("count", "Two", "Three")
	ut.AssertNil(err)
	err = f.Submit()
	ut.AssertNil(err)
	ut.AssertEquals(`count=2&count=3&submit=submitted`, string(bow.body))

	// select multi count by label
	err = f.RemoveValue("count", "2")
	ut.AssertNil(err)
	err = f.Submit()
	ut.AssertNil(err)
	ut.AssertEquals(`count=3&submit=submitted`, string(bow.body))
}

func TestBrowserFormClickByValue(t *testing.T) {
	ts := setupTestServer(`
<!doctype html>
<html>
	<head>
		<title>Echo Form</title>
	</head>
	<body>
		<form method="post" action="/" name="default">
			<input type="text" name="age" value="" />
			<input type="submit" name="submit" value="submitted1" />
			<input type="SUBMIT" name="submit" value="submitted2" />
		</form>
	</body>
</html>`, t)
	defer ts.Close()

	bow := newBrowser()

	err := bow.Open(ts.URL)
	ut.AssertNil(err)

	f, err := bow.Form("[name='default']")
	ut.AssertNil(err)

	err = f.Input("age", "55")
	ut.AssertNil(err)
	err = f.ClickByValue("submit", "submitted2")
	ut.AssertNil(err)
	ut.AssertContains("age=55", bow.Body())
	ut.AssertContains("submit=submitted2", bow.Body())
}

func TestSubmitMultipart(t *testing.T) {
	ts := setupTestServer(`<!doctype html>
<html>
	<head>
		<title>multipart form</title>
	</head>
	<body>
		<form method="post" action="/" name="default" enctype="multipart/form-data">
			<input type="text" name="comment" value="" />
			<INPUT type="FILE" name="image" />
			<input type="submit" name="submit" value="submitted1" />
		</form>
	</body>
</html>
`, t)
	defer ts.Close()

	bow := newBrowser()

	err := bow.Open(ts.URL)
	ut.AssertNil(err)

	f, err := bow.Form("[name='default']")
	ut.AssertNil(err)

	f.Input("comment", "my profile picture")
	imgData, err := base64.StdEncoding.DecodeString(image)
	err = f.File("image", "profile.png", bytes.NewBuffer(imgData))
	ut.AssertNil(err)
	err = f.Submit()
	ut.AssertNil(err)
	ut.AssertContains("comment=my+profile+picture", bow.Body())
	ut.AssertContains("image=profile.png", bow.Body())
	ut.AssertContains(fmt.Sprintf("profile.png=%s", url.QueryEscape(image)), bow.Body())
}

func setupTestServer(html string, t *testing.T) *httptest.Server {
	ut.Run(t)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			fmt.Fprint(w, html)
		} else {
			ct := r.Header.Get("Content-Type")
			if strings.LastIndex(ct, "multipart/form-data") != -1 {

				r.ParseMultipartForm(1024 * 1024) // max 1MB ram
				values := url.Values{}
				for k, av := range r.MultipartForm.Value {
					for _, v := range av {
						values.Add(k, v)
					}
				}
				for k, af := range r.MultipartForm.File {
					for _, fh := range af {
						values.Add(k, fmt.Sprintf("%s", fh.Filename))
						f, _ := fh.Open()
						data, _ := ioutil.ReadAll(f)
						val := base64.StdEncoding.EncodeToString(data)
						values.Add(fh.Filename, val)
					}
				}
				fmt.Fprint(w, values.Encode())
			} else {
				r.ParseForm()
				fmt.Fprint(w, r.Form.Encode())
			}

		}
	}))

	return ts
}

var image = `iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAYAAABzenr0AAACjUlEQVRYR+2Wy6oiMRCG4x0V76CIim504fu/gk/hQlS8LhQVERR15stQPZl0a3IOwtlMgXSnU/nrq0p12thoNHqqH7T4D8bWof8DeFXg+fzTJnKldK57c/7dNjsBYrGY4hcFkUwmVSqV0vMmFL7y7F1w5pIuh8fjoeLxv5y5XE61Wi1VKpUUANj9fleHw0Etl0t1Op0CSR8QJ4BkD0i321XtdltXw8wQkGq1qmq1moaYTqeuvIJ5J4B4djodHVy2xI4gQM1mUwPOZjNdOVcvOHuAQJQdAF9ji4rFojM4el4ACPo2lWwZa3zMCUAJy+Wyj5b2wZ/SFwqFf5r3lYATIJ1OB93+SsR8LhUAIpPJOJc4AUTB1Ux2pFfNavs5AW63m+IVlMxsAXssryjXy+ViT4fGTgAC73Y7vdCnEfEBmIPJp2pOAMTW67UW456rCJtj7mUM7Gq1CmUb9cAJQEYcr4vFQldAgnA1zazOdrtV+/3eq2JOAIKQ8Xw+11lJL3A1g8rebzYbNZlMAtiorM1nTgBzH8nseDxGvt9SKbZLIF3BmY/Z/wklEzMwZwFHcb1ejwxuBmJrAKVi1+s1mDIrZ/o7P0b5fF71+339PfAxAjUaDX0SjsdjdT6f3y4LbQEC0mAEHw6HOrhZkbeKvyfxZQ1r0cDQNHtGNEIAOCYSCX38DgYD/Y8Hi1osIvZVfFmLBlpo2m8O60IAsrjX66lsNqsDmz87mD22/dFAC4tKIhKAstFwnzK00PQCYP/kb9enAN5phirAJ5TvfxTtd4HQQjPq8xwCqFQqQfCvdP4rONEAAm3bQgC8v5/MXgKiibZtIQCaxaS2F3x1LMmgKWeCqREC4Nj9ROltUDTRtu0X2hs2IkarWoAAAAAASUVORK5CYII=`
