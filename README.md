Surf
====

[![Build Status](https://img.shields.io/travis/headzoo/surf/master.svg)](https://travis-ci.org/headzoo/surf)
[![Documentation](https://img.shields.io/badge/documentation-latest-blue.svg)](http://surf.readthedocs.org/)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/headzoo/surf/master/LICENSE.md)
[![GitHub Stars](https://img.shields.io/github/stars/headzoo/surf.svg)](https://github.com/headzoo/surf/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/headzoo/surf.svg)](https://github.com/headzoo/surf/network)

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
Go was started by [Sean Hickey](https://github.com/headzoo) (headzoo) to learn more about the Go programming language.
The idea to create Surf was born in [this Reddit thread](http://www.reddit.com/r/golang/comments/2efw1q/mechanize_in_go/cjz4lze).

Surf uses the awesome [goquery](https://github.com/PuerkitoBio/goquery) by Martin Angers, and
was written using [Intellij](http://www.jetbrains.com/idea/) and
the [golang plugin](http://plugins.jetbrains.com/plugin/5047).

Contributing authors:

* [Haruyama Seigo](https://github.com/haruyama)
* [Tatsushi Demachi](https://github.com/tatsushid)
* [Charl Matthee](https://github.com/charl)
* [Matt Holt](https://github.com/mholt)


### Contributing
See [CONTRIBUTING.md] for more information on contributing to the project.


### License
Surf is released open source software released under The MIT License (MIT).
See [LICENSE.md](https://raw.githubusercontent.com/headzoo/surf/master/LICENSE.md) for more information.
