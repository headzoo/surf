Surf
====

[![Build Status](https://img.shields.io/travis/headzoo/surf/master.svg)](https://travis-ci.org/headzoo/surf)
[![Documentation](https://img.shields.io/badge/documentation-latest-blue.svg)](http://surf.readthedocs.org/)
[![MIT license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/headzoo/surf/master/LICENSE.md)

Surf is a Go (golang) library that implements a virtual browser that you can control pragmatically. Just like a
real browser you can open pages, follow links, bookmark pages, submit forms, and many other things.

### Installation
Download the library using go.  
`go get github.com/headzoo/surf`

Import the library into your project.  
`import "github.com/headzoo/surf"`


### Quick Start
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

Complete documentation is available on [Read the Docs](http://surf.readthedocs.org/)

### Credits
Surf uses the awesome [goquery](https://github.com/PuerkitoBio/goquery) by Martin Angers, and
was written using [Intellij](http://www.jetbrains.com/idea/) and
the [golang plugin](http://plugins.jetbrains.com/plugin/5047).

Contributing authors:

* [Haruyama Seigo](https://github.com/haruyama)
* [Tatsushi Demachi](https://github.com/tatsushid)
* [Charl Matthee](https://github.com/charl)
* [Matt Holt](https://github.com/mholt)


### Use Cases
* Interacting with sites that do not have public APIs.
* Testing/Stressing your sites.
* Scraping sites.
* Creating a web crawler.


### License
Surf is released open source software released under The MIT License (MIT).
See [LICENSE.md](https://raw.githubusercontent.com/headzoo/surf/master/LICENSE.md) for more information.
