package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Generate a random number from 1 to 10
	rand.Seed(42)
	number := rand.Intn(10) + 1

	// Play the guessing game
	turns := 0
	for {
		turns++

		// Prompt the user to guess a number
		fmt.Print("Pick a number from 1 to 10: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return
		}

		// Clean up the input and convert it to an integer
		input = strings.TrimSpace(input)
		guess, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input")
			continue
		}

		// Check if the guess is correct, too high, or too low
		if guess == number {
			// Win message
			fmt.Printf("You won! You guessed %d, which is correct. It took you %d turns to guess correctly!\r\n\r\nPress [ R ] to play again or [ Q ] to quit.\r\n", guess, turns)
			break
		} else if guess < number {
			fmt.Println("Try higher")
		} else {
			fmt.Println("Try lower")
		}
	}
}
