// chat.go
package main

import (
	"bufio"
	"net"
	"strings"
	"sync"
	"thetelevision/channel"
	"thetelevision/system"
)

const VARI_VERSION = "v1.0.25042023"

type ChatDoor struct{}

func (d *ChatDoor) Name() string {
	return "Chat"
}

func (d *ChatDoor) Description() string {
	return "TeleVision VariChat"
}

func (d *ChatDoor) Play(conn net.Conn, cs *channel.ChannelServer, user system.UserInfo) {
	// Add the user to the chat room
	cs.AddClient(conn)

	// Send the welcome message to the user
	conn.Write([]byte("\r\nWelcome to the VariChat " + VARI_VERSION + "\r\nType '/q' to exit.\r\n"))

	// Start listening for messages
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		reader := bufio.NewReader(conn)
		for {
			// Read input from the user
			userInput, err := reader.ReadString('\n')
			if err != nil {
				break
			}
	
			// Remove the newline character from the input
			userInput = strings.TrimSpace(userInput)

			// If the user wants to quit, break out of the loop
			if userInput == "/q" {
				break
			}

			// Handle commands
			if strings.HasPrefix(userInput, "/") {
				command := strings.Split(userInput, " ")[0]
				switch command {
				case "/w":
					cs.WhoIsOnline(conn)
					continue
				case "/q":
					// Exit the chat room
					goto exitLoop
				case "/p":
					// Send a private message to a user
					parts := strings.SplitN(userInput, " ", 2)
					if len(parts) != 2 {
						conn.Write([]byte("\r\nInvalid command. Please try again.\r\n"))
						continue
					}
					toUser := parts[0][2:]
					message := "[P][ " + user.Username + " ]: " + parts[1]
					cs.SendPrivateMessage(conn, toUser, message)
				default:
					conn.Write([]byte("\r\nInvalid command. Please try again.\r\n"))
					continue
				}
			} else {
				// Build message
				message := "[ " + user.Username + " ]: " + userInput
				// Broadcast the message to all clients
				cs.Broadcast(conn, message)
			}
		}
	exitLoop:
		wg.Done()
	}()

	// Wait for the user to exit the chat room
	wg.Wait()

	// Remove the user from the chat room
	cs.RemoveClient(conn)

	// Exit the module and return to the previous loop
	conn.Write([]byte("\r\nExiting chat room...\r\n"))
}

var Door ChatDoor

func main() {}