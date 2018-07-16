//Adding basic go package to allow for 'go get'
package main

import (
	stdio "io"
	"os"
	"os/signal"
	"syscall"

	"github.com/lauchlan105/multiterm/multiterm"
)

type io struct {
	in  pipe
	out pipe
}

type pipe struct {
	r *stdio.PipeReader
	w *stdio.PipeWriter
}

func main() {

	// inR, inW := stdio.Pipe()
	// outR, outW := stdio.Pipe()

	// newPipe := io{
	// 	in: pipe{
	// 		inR,
	// 		inW,
	// 	},
	// 	out: pipe{
	// 		outR,
	// 		outW,
	// 	},
	// }

	// go runCommand([]string{"go", "run", "./maintwo/test.go"}, &newPipe)

	// //Print anything printed to the output pipe
	// go func() {
	// 	scanner := bufio.NewScanner(newPipe.out.r)
	// 	for scanner.Scan() {
	// 		fmt.Println("Output: " + scanner.Text())
	// 	}
	// }()

	// //Write any input to the input pipe
	// go func() {
	// 	reader := bufio.NewReader(os.Stdin)
	// 	for {
	// 		b, _ := reader.ReadBytes('\n')
	// 		newPipe.in.w.Write(b)
	// 	}
	// }()

	window := multiterm.Init()
	defer window.Wait()
	window.Start()

	tab := window.NewTab()
	tab.Open()

	// tab.RunCommand([]string{"go", "run", "maintwo/test.go"})

	// oneTab := window.NewTab()
	// oneTab.Open()
	// oneTab.RunCommand([]string{"go", "run", "maintwo/test.go"})
}

func stayAlive() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
