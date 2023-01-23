package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func main() {
	var name = "The Crawl"
	var numRooms = 3
	var roomSize = 5
	var numMonsters = 5
	var playerPos = [2]int{0, 0}
	var monsters = make([][2]int, numMonsters)
	var dungeon = make([][]rune, roomSize*numRooms)
	var gameOver = false
	var victory = false

	fmt.Println("\r\nWelcome to", name, "!")
	// Initialize dungeon
	initializeDungeon(numRooms, roomSize, numMonsters, playerPos, monsters, dungeon)

	// Main game loop
	for !gameOver {
		// Print dungeon
		drawDungeon(name, dungeon, playerPos, monsters)

		// Read input from user
		fmt.Println("Enter your move (w/a/s/d): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Process input
		switch input {
		case "w":
			if playerPos[0] < 1 {
				playerPos[0] = 2
			}
			playerPos[0]--
		case "a":
			if playerPos[1] < 1 {
				playerPos[1] = 2
			}
			playerPos[1]--
		case "s":
			playerPos[0]++
		case "d":
			playerPos[1]++
		}

		// Check for win or loss conditions
		if dungeon[playerPos[0]][playerPos[1]] == 'E' {
			victory = true
			gameOver = true
		}
		for _, monster := range monsters {
			if monster[0] == playerPos[0] && monster[1] == playerPos[1] {
				gameOver = true
				break
			}
		}
	}

	// Print game over message
	drawDungeon(name, dungeon, playerPos, monsters)
	if victory {
		fmt.Println("Congratulations! You have defeated the dungeon and retrieved the treasure!")
	} else {
		fmt.Println("You have been defeated by the monsters in the dungeon. Better luck next time!")
	}
}

func initializeDungeon(numRooms int, roomSize int, numMonsters int, playerPos [2]int, monsters [][2]int, dungeon [][]rune) {
	room := rand.Intn(numRooms)
	for i := 0; i < numMonsters; i++ {
		monsters[i][0] = rand.Intn(roomSize) + (room * roomSize)
		monsters[i][1] = rand.Intn(roomSize) + (room * roomSize)
	}
	for i := 0; i < len(dungeon); i++ {
		for j := 0; j < len(dungeon[i]); j++ {
			if i == playerPos[0] && j == playerPos[1] {
				dungeon[i][j] = 'P'
			} else if i == room*roomSize && j == room*roomSize {
				dungeon[i][j] = 'E'
			} else {
				dungeon[i][j] = '.'
			}
		}
	}
}

func drawDungeon(name string, dungeon [][]rune, playerPos [2]int, monsters [][2]int) {
	// Clear screen
	fmt.Print("\x1b[2J")
	fmt.Println(name)

	// Draw dungeon
	for i := 0; i < len(dungeon); i++ {
		for j := 0; j < len(dungeon[i]); j++ {
			if i == playerPos[0] && j == playerPos[1] {
				fmt.Print("P")
			} else {
				fmt.Printf("%c", dungeon[i][j])
			}
		}
		fmt.Println()
	}
	for _, monster := range monsters {
		if dungeon[monster[0]][monster[1]] == '.' {
			fmt.Println("M")
		}
	}
}
