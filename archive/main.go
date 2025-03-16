package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func main() {
	// Save the original terminal state
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	// Restore the terminal state when the program exits
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	// Read single bytes directly
	buffer := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(buffer)
		if err != nil {
			break
		}
		// Print the actual byte value received
		fmt.Printf("Received byte: %d\n", buffer[0])

		// Exit on 'q'
		if buffer[0] == 'q' {
			break
		}
	}
}
