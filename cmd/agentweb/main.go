package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("agent-project-manager: console started")
	fmt.Println("Type something and press Enter (type 'exit' to quit).")

	in := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !in.Scan() {
			break // EOF / error
		}

		text := in.Text()
		if text == "exit" || text == "quit" {
			fmt.Println("bye")
			return
		}

		fmt.Printf("you typed: %s\n", text)
	}

	if err := in.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "input error: %v\n", err)
	}
}
