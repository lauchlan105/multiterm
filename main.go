//Adding basic go package to allow for 'go get'
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/lauchlan105/multiterm/multiterm"
)

func main() {
	window := multiterm.Init()
	defer window.Wait()
	window.Start()
}

func stayAlive() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
