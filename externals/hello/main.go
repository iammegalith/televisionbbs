package main

import (
	"fmt"
)

func main() {
	// Create a channel to communicate with the handleDoor function
	done := make(chan bool)

	// Print "Hello, world!" to the screen
	fmt.Println("Hello, world!")

	// Signal the handleDoor function that the program has finished
	done <- true

	// Wait for the handleDoor function to receive the signal
	<-done
}
