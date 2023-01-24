package chatroom

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"televisionbbs/util"

	"gopkg.in/gcfg.v1"
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

type ActionConfig struct {
	Actions map[string]string
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

func ReadActions(configpath string, actionFile string) (map[string]string, error) {
	var cfg ActionConfig
	fmt.Println("Reading actions from", configpath+actionFile)
	if err := gcfg.ReadFileInto(&cfg, configpath+actionFile); err != nil {
		return nil, err
	}
	return cfg.Actions, nil
}

func handleActions(conn net.Conn, argString string, actions map[string]string, chat *Chat) (actionMessage string) {
	action, ok := actions[argString]
	if !ok {
		fmt.Fprintln(conn, "Invalid action.")
		return
	}

	// Get the action to perform from the actions.config file
	msg := action
	// Check if the specified action is valid
	if msg == "" {
		fmt.Fprintln(conn, "Invalid action.")
		return
	}
	return msg
}

func MultiUserChat(conn net.Conn, conUser string, cutranslation int, configpath string) (ExitStatus string) {
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
	if cutranslation == 1 {
		fmt.Fprint(conn, "\r\n"+util.ANSI_BLUE_BG+util.ANSI_WHITE+"[ TVChat v"+TVCHAT_VERSION+" ]"+util.ANSI_RESET+"\r\n/? For help\r\n")
	} else {
		fmt.Fprint(conn, "\r\n[ TVChat v"+TVCHAT_VERSION+" ]\r\n/? For help\r\n")
	}
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
			parts := strings.SplitN(message, " ", 4)
			var commandString = ""
			if len(parts) > 0 {
				commandString = strings.TrimPrefix(parts[0], "/")
				commandString = strings.TrimSpace(commandString)
			}
			fmt.Println("Command:", commandString)
			switch commandString {
			case "p":
				// Send the private message to the specified user
				recipient := strings.TrimSpace(parts[1])
				msgString := strings.TrimSpace(strings.Join(parts[2:], " "))
				msg := fmt.Sprintf("P[%s] : %s", username, msgString+"\r\n")
				for _, c := range chat.Connections {
					if LoggedInUsers[c] == recipient {
						fmt.Fprintln(c, msg)
						break
					}
				}
			case "w":
				// Handle the who command
				var userList string
				for _, c := range chat.Connections {
					userList += fmt.Sprintf("%s, ", LoggedInUsers[c])
				}
				fmt.Fprintln(conn, "Connected Users: "+userList+"\r\n")
			case "q":
				// Handle the quit command
				fmt.Println(username, "disconnected")
				delete(LoggedInUsers, conn)
				chat.RemoveConnection(conn)
				fmt.Fprintln(conn, "Goodbye!")
				return "EXIT"
			case "?":
				// Handle the help command
				fmt.Fprint(conn, "\r\nCommands:")
				fmt.Fprint(conn, "\r\n    /p [user] [message] - Send a private message to the specified user")
				fmt.Fprint(conn, "\r\n    /a [action] [target user] - Perform an action (e.g. /a wave John)")
				fmt.Fprint(conn, "\r\n    /w - Show connected users")
				fmt.Fprint(conn, "\r\n    /q - Quit the chat")
				fmt.Fprint(conn, "\r\n    /? - Show this help message\r\n")
			default:
				fmt.Fprint(conn, "\r\n>>> Unknown command. Type /? for help.\r\n")
			}
		}
	}
	return "success"
}
