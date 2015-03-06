# Surf

Surf is a Go (golang) library that implements a virtual browser that you control pragmatically. Just like a real
browser you can open pages, follow links, bookmark pages, submit forms, and many other things. 

```go
package main

import (
	"github.com/headzoo/surf"
	"fmt"
)

func main() {
	bow := surf.NewBrowser()
	err := bow.Open("http://golang.org")
	if err != nil {
		panic(err)
	}
	
	// Outputs: "The Go Programming Language"
	fmt.Println(bow.Title())
}
```


### Installation
Download the library using go.

```sh
$ go get github.com/headzoo/surf
```

Import the library into your project.

```go
import "github.com/headzoo/surf"
```
