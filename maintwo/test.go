package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	fmt.Println("Starting subroutine...")

	reader := bufio.NewReader(os.Stdin)
	for {
		str, _ := reader.ReadString('\n')
		fmt.Println("Test recieved: \"" + str[:len(str)-1] + "\"")
	}
}
