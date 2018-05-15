//Adding basic go package to allow for 'go get'
package main

import (
	"github.com/lauchlan105/multiterm/multiterm"
)

type t struct {
	tabs []tab
}

type tab struct {
	id int
}

func main() {

	window, _ := multiterm.Init()
	defer window.Stop()
	window.Start()

}
