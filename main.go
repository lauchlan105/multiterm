//Adding basic go package to allow for 'go get'
package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
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

	inR, inW := stdio.Pipe()
	outR, outW := stdio.Pipe()

	newPipe := io{
		in: pipe{
			inR,
			inW,
		},
		out: pipe{
			outR,
			outW,
		},
	}

	go runCommand([]string{"go", "run", "./maintwo/test.go"}, &newPipe)

	//Print anything printed to the output pipe
	go func() {
		scanner := bufio.NewScanner(newPipe.out.r)
		for scanner.Scan() {
			fmt.Println("Output: " + scanner.Text())
		}
	}()

	//Write any input to the input pipe
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			b, _ := reader.ReadBytes('\n')
			newPipe.in.w.Write(b)
		}
	}()

	/*
		window := multiterm.Init()
		defer window.Wait()
		window.Start()
	*/
}

func runCommand(s []string, io *io) {
	cmd := exec.Command(s[0], s[1:]...)

	cmd.Stdout = io.out.w
	cmd.Stdin = io.in.r

	cmd.Run()
}

func stayAlive() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
