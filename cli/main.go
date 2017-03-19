package main

import (
	"github.com/headzoo/surf"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		return
	}

	//surf.Debugging = true
	browser := surf.NewBrowser()
	browser.Document.AddEventListener(surf.OnLoad, func(e *surf.Event) {
		println("Title:", browser.Document.Title())
		println("Content Type:", browser.Document.ContentType())
		println("Character Set:", browser.Document.CharacterSet())
	})
	if err := browser.SendGET(os.Args[1]); err != nil {
		panic(err)
	}
}
