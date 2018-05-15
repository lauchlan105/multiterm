//Adding basic go package to allow for 'go get'
package main

import (
	"time"

	"github.com/nsf/termbox-go"

	"github.com/lauchlan105/multiterm/multiterm"
)

type t struct {
	tabs []tab
}

type tab struct {
	id int
}

func main() {

	window, _ := multiterm.Init(termbox.ColorDefault, termbox.ColorDefault)
	defer window.Stop()
	time.Sleep(10 * time.Second)

}
