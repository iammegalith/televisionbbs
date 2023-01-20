package chatroom

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"televisionbbs/util"
)

type Chat struct {
	Connections []net.Conn
	Messages    chan string
}

func (c *Chat) Run() {
	for {
		select {
		case msg := <-c.Messages:
			for _, conn := range c.Connections {
				fmt.Fprintln(conn, msg)
			}
		}
	}
}

func (c *Chat) AddConnection(conn net.Conn) {
	c.Connections = append(c.Connections, conn)
}

func (c *Chat) RemoveConnection(conn net.Conn) {
	for i, cn := range c.Connections {
		if cn == conn {
			c.Connections = append(c.Connections[:i], c.Connections[i+1:]...)
			break
		}
	}
}

func MultiUserChat(conn net.Conn) {
	// Create a new chat instance
	chat := &Chat{
		Connections: []net.Conn{conn},
		Messages:    make(chan string),
	}
	// Add the connection to the chat
	chat.AddConnection(conn)
	// Start the chat
	go chat.Run()
	fmt.Fprint(conn, "[ Entering TVChat ]"+util.CR_LF)
	// Read messages from the connection and send them to the chat
	username := util.LoggedInUsers[conn]
	chat.Messages <- fmt.Sprintf("%s: %s", username, " has entered TVChat"+util.CR_LF)
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println(username, "disconnected")
			} else {
				fmt.Println("Error reading message:", err)
			}
			// Remove the connection from the chat
			chat.RemoveConnection(conn)
			break
		}
		// Check if the message is a command
		if len(message) > 0 && message[0] == '/' {
			parts := strings.SplitN(message, " ", 2)
			fmt.Println("Command:", parts[0])
			switch parts[0] {
			case "/pm":
				// Send the private message to the specified user
				recipient := parts[1]
				msg := fmt.Sprintf("%s: %s", username, parts[2])
				for _, c := range chat.Connections {
					if util.LoggedInUsers[c] == recipient {
						fmt.Fprintln(c, msg)
						break
					}
				}
			case "/who":
				// Send a list of logged-in users to the connection
				fmt.Fprintln(conn, "Logged-in users:")
				for _, c := range chat.Connections {
					fmt.Fprintln(conn, util.LoggedInUsers[c])
				}
			case "/q":
				// Remove the connection from the chat
				chat.RemoveConnection(conn)
			}
		} else {
			// Send the message to all connections in the chat
			chat.Messages <- fmt.Sprintf("%s: %s", username, message)
		}
	}
}
