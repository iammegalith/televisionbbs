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
	// Print the welcome message
	fmt.Println("This is an example of how to write a super simple door game for TeleVision BBS.")
	fmt.Println("Welcome to the guessing game!")
	fmt.Println("I'm thinking of a number between 1 and 100. Can you guess it?")
	fmt.Println("Enter 'quit' to exit.")

	// Generate a random number
	randNum := 1 + int(99*rand.Float64())

	// Use a loop to keep the game going until the user quits
	for {
		// Read the user's guess
		fmt.Print("Your guess: ")
		reader := bufio.NewReader(os.Stdin)
		guess, _ := reader.ReadString('\n')
		// Remove leading and trailing whitespace
		guess = strings.TrimSpace(guess)
		// Check if the user wants to quit
		if guess == "quit" {
			break
		}
		// Parse the user's guess as an integer
		guessInt, err := strconv.Atoi(guess)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number or 'quit'.")
			continue
		}
		// Check if the user's guess is correct
		if guessInt == randNum {
			fmt.Println("Correct! You win!")
			continue
		}
		// Tell the user if their guess was too high or too low
		if guessInt < randNum {
			fmt.Println("Your guess is too low.")
		} else {
			fmt.Println("Your guess is too high.")
		}
	}
	fmt.Println("Thanks for playing!")
	return
}
