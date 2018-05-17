//Adding basic go package to allow for 'go get'
package main

import (
	"fmt"

	"github.com/lauchlan105/multiterm/multiterm"
)

type testFunc struct {
	f func()
}

func test() {
	fmt.Println("tester func")
}

func main() {

	window, _ := multiterm.Init()
	defer window.Stop()
	window.Start()

}
