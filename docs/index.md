# Surf

Surf is a Go (golang) library that implements a virtual browser that you control pragmatically. Just like a real
browser you can open pages, follow links, bookmark pages, submit forms, and many other things. 

[![Build Status](https://img.shields.io/travis/headzoo/surf/master.svg?style=flat-square)](https://travis-ci.org/headzoo/surf)
[![Github](https://img.shields.io/badge/source-github-blue.svg?style=flat-square)](https://github.com/headzoo/surf/)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://raw.githubusercontent.com/headzoo/surf/master/LICENSE.md)
[![GitHub Stars](https://img.shields.io/github/stars/headzoo/surf.svg?style=flat-square)](https://github.com/headzoo/surf/stargazers)
[![GitHub Forks](https://img.shields.io/github/forks/headzoo/surf.svg?style=flat-square)](https://github.com/headzoo/surf/network)

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
