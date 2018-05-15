//Adding basic go package to allow for 'go get'
package main

import (
	"time"

	"github.com/lauchlan105/multiterm/multiterm"
)

func main() {

	window := multiterm.Start()
	defer window.Stop()

	time.Sleep(10 * time.Second)

}
