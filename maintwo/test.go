package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {

	fmt.Println("Starting subroutine...")

	go func() {
		for {
			fmt.Println("POWER NAP!")
			time.Sleep(1 * time.Second)
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	go func() {
		for {
			str, _ := reader.ReadString('\n')
			fmt.Println("Test recieved: \"" + str[:len(str)-1] + "\"")
		}
	}()

	for {

	}
}
