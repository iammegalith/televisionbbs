package messagebases

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/ssh/terminal"
)

func messageEditor(message *Message) bool {
	state, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(int(os.Stdin.Fd()), state)

	reader := bufio.NewReader(os.Stdin)
	var newMessage strings.Builder
	var saveMessage = true
	var newLine = true
	var row, col, _ = terminal.GetSize(int(os.Stdin.Fd()))
	row -= 1

	// cursor positions
	var cursorX, cursorY int
	cursorX, cursorY = 0, 6

	// slice of lines in message
	lines := strings.Split(message.Body, "\n")
	topLine := 0

	for {
		// clear the screen and reprint the message
		fmt.Printf("\033c")
		fmt.Printf("From: %s\n", message.From)
		fmt.Printf("To: %s\n", message.To)
		fmt.Printf("Subject: %s\n", message.Subject)
		fmt.Printf("Message Base Name: %s\n", message.MessageBaseName)
		for i := topLine; i < len(lines) && i < topLine+row-6; i++ {
			line := lines[i]
			for len(line) > col {
				fmt.Println(line[:col])
				line = line[col:]
			}
			fmt.Println(line)
		}
		fmt.Printf("Commands: /s save and exit | /x exit without saving\n")
		fmt.Printf("\033[%d;%dH", cursorY, cursorX)

		char, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		if char == 27 { // Esc key, starts an arrow key escape sequence
			char, _, err = reader.ReadRune()
			if err != nil {
				if err == io.EOF {
					break
				}
				panic(err)
			}
			if char == '[' {
				char, _, err = reader.ReadRune()
				if err != nil {
					if err == io.EOF {
						break
					}
					panic(err)
				}

				if char == 'A' { // Up arrow key
					cursorY--
					if cursorY < 6 {
						cursorY = 6
						topLine = math.max(0, topLine-1)

					} else if char == 'B' { // Down arrow key
						cursorY++
						if cursorY >= row {
							cursorY = row - 1
							topLine = math.min(len(lines)-(row-6), topLine+1)
						}
					} else if char == 'C' { // Right arrow key
						cursorX++
					} else if char == 'D' { // Left arrow key
						cursorX--
					}
					if cursorY > topLine+row-6 {
						cursorY = topLine + row - 6
					}
					if cursorY < 6 {
						cursorY = 6
					}
					fmt.Printf("\033[%d;%dH", cursorY, cursorX)
				}
			} else if char == '\n' {
				newLine = true
			} else if newLine && char == '/' {
				char, _, err = reader.ReadRune()
				if err != nil {
					if err == io.EOF {
						break
					}
					panic(err)
				}
				if char == 's' {
					break
				} else if char == 'x' {
					saveMessage = false
					break
				} else {
					newLine = false
					newMessage.WriteRune('/')
					newMessage.WriteRune(char)
				}
			} else if char == 127 || char == '\b' {
				newMessage.Truncate(newMessage.Len() - utf8.RuneLen(char))
				if cursorX > 0 {
					cursorX--
				} else if cursorY > 6 {
					cursorX = col - 1
					cursorY--
				} else {
					cursorX = 0
				}
				fmt.Printf("\033[%d;%dH", cursorY, cursorX)
				fmt.Printf(" \033[%d;%dH", cursorY, cursorX)
			} else {
				newLine = false
				newMessage.WriteRune(char)
			}
		}
	}

	message.Body = newMessage.String()
	return saveMessage
}
