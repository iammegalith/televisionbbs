package chatroom

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"unicode/utf8"
)

const (
	TVCHAT_VERSION = "0.1"
)

var LoggedInUsers = make(map[net.Conn]string)
var ConnectionChannels = make(map[net.Conn]chan string)

type Chat struct {
	Connections []net.Conn
	Messages    chan string
}

func (c *Chat) Run() {
	for msg := range c.Messages {
		for _, conn := range c.Connections {
			ConnectionChannels[conn] <- msg
			fmt.Println("Message sent to", LoggedInUsers[conn])
		}
	}
	for _, conn := range c.Connections {
		delete(LoggedInUsers, conn)
		delete(ConnectionChannels, conn)
	}
}

func (c *Chat) AddConnection(conn net.Conn) {
	c.Connections = append(c.Connections, conn)
	ConnectionChannels[conn] = make(chan string)
}

func (c *Chat) RemoveConnection(conn net.Conn) {
	for i, cn := range c.Connections {
		if cn == conn {
			fmt.Println("\r\nRemoving connection")
			c.Connections = append(c.Connections[:i], c.Connections[i+1:]...)
			delete(LoggedInUsers, conn)
			delete(ConnectionChannels, conn)
			return
		}
	}
}

func TrimFirstChar(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func promptUser(conn net.Conn) {
	fmt.Fprint(conn, "> ")
}

func MultiUserChat(conn net.Conn, conUser string) (ExitStatus string) {
	// Create a new chat instance
	chat := &Chat{
		Connections: []net.Conn{},
		Messages:    make(chan string),
	}
	// Add the connection to the chat
	chat.AddConnection(conn)
	go func() {
		for {
			msg := <-ConnectionChannels[conn]
			fmt.Fprintln(conn, msg)
		}
	}()

	// Start the chat
	go chat.Run()
	fmt.Fprint(conn, "[ TVChat v"+TVCHAT_VERSION+" ]\r\n")
	// Read messages from the connection and send them to the chat
	if _, ok := LoggedInUsers[conn]; !ok {
		LoggedInUsers[conn] = conUser
	}
	fmt.Println("User:", LoggedInUsers[conn])
	username := LoggedInUsers[conn]
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println(username, "disconnected")
				delete(LoggedInUsers, conn)
			} else {
				fmt.Println("Error reading message:", err)
			}
			// Remove the connection from the chat
			chat.RemoveConnection(conn)
			break
		}
		if len(message) > 0 && message[0] != '/' {
			for _, c := range chat.Connections {
				message = strings.TrimSuffix(message, "\n")
				ConnectionChannels[c] <- fmt.Sprintf("[%s] : %s", username, message)
			}
		}
		// Check if the message is a command
		if len(message) > 0 && message[0] == '/' {

			parts := strings.SplitN(message, " ", 3)
			var commandString = strings.TrimSpace(TrimFirstChar(parts[0]))
			fmt.Println("Command:", commandString)
			switch commandString {
			case "who":
				// Send a list of logged-in users to the connection
				fmt.Fprintln(conn, "\r\nLogged-in Users:")
				for _, user := range LoggedInUsers {
					fmt.Fprintln(conn, "\r\n - [ "+user+" ]\r\n")
				}
			case "q":
				fmt.Fprintln(conn, "Goodbye!")
				// Remove the connection from the chat
				chat.RemoveConnection(conn)
				return "EXIT"
			}
		}
	}
	return "EXIT"
}
