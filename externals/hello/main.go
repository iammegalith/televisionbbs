package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1) // Add 1 to the WaitGroup
	done := make(chan bool)

	// Print "Hello, world!" to the screen
	fmt.Println("Hello, world!")

	// Signal the handleDoor function that the program has finished
	done <- true

	// Wait for the handleDoor function to receive the signal
	<-done
	wg.Done() // Signal that the program has finished
	wg.Wait() // Wait for all goroutines to finish executing
	os.Exit(0)
}
