package channel

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"thetelevision/system"
	"time"
)

type Client struct {
	conn net.Conn
	send chan string
}

type Message struct {
	sender  net.Conn
	content string
}

type ChannelServer struct {
	clients    map[net.Conn]*Client
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
	mutex      sync.Mutex
}

func NewChannelServer() *ChannelServer {
	return &ChannelServer{
		clients:    make(map[net.Conn]*Client),
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (cs *ChannelServer) AddClient(conn net.Conn) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	cs.clients[conn] = &Client{conn: conn, send: make(chan string)}
}

func (cs *ChannelServer) RemoveClient(conn net.Conn) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	delete(cs.clients, conn)
	conn.Close()
}

func (cs *ChannelServer) Broadcast(sender net.Conn, message string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for conn, client := range cs.clients {
		if conn != sender {
			client.send <- fmt.Sprintf("[%s] %s", time.Now().Format("15:04:05"), message)
		}
	}
}

func (cs *ChannelServer) run() {
	for {
		select {
		case client := <-cs.register:
			cs.clients[client.conn] = client
		case client := <-cs.unregister:
			if _, ok := cs.clients[client.conn]; ok {
				delete(cs.clients, client.conn)
				close(client.send)
			}
		case msg := <-cs.broadcast:
			for _, client := range cs.clients {
				if client.conn != msg.sender {
					client.send <- msg.content
				}
			}
		}
	}
}

func (cs *ChannelServer) ListenAndServe() error {

	go cs.run() // Start the goroutine to handle channels
    var addr string = net.JoinHostPort(system.BBSConfig.ChannelsListenAddr, strconv.Itoa(system.BBSConfig.ChannelsPort))

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer listener.Close()

	log.Printf("TeleVision Channels tuned in on %s\n", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		client := &Client{conn: conn, send: make(chan string)}
		cs.register <- client

		go func(client *Client) {
			defer func() {
				cs.unregister <- client
				client.conn.Close()
			}()

			buf := make([]byte, 1024)

			for {
				n, err := client.conn.Read(buf)
				if err != nil {
					fmt.Printf("Error reading from connection: %v\n", err)
					break
				}

				message := string(buf[:n])
				cs.broadcast <- Message{sender: client.conn, content: message}
			}
		}(client)
	}
}

func (cs *ChannelServer) SendPrivateMessage(sender net.Conn, recipientName string, message string) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	for conn, client := range cs.clients {
		if conn != sender {
			continue
		}
		// Find the client with the specified recipient name
		for rConn, rClient := range cs.clients {
			if rConn == conn {
				continue
			}
			if rClient.conn.RemoteAddr().String() == client.conn.RemoteAddr().String() {
				if rClient.conn == client.conn {
					continue
				}
				if rClient.conn != nil {
					rClient.send <- fmt.Sprintf("[%s][Private Message from %s]: %s", time.Now().Format("15:04:05"), client.conn.RemoteAddr().String(), message)
				}
				return
			}
			if strings.EqualFold(strings.ToLower(rClient.conn.RemoteAddr().String()), strings.ToLower(recipientName)) {
				rClient.send <- fmt.Sprintf("[%s][Private Message from %s]: %s", time.Now().Format("15:04:05"), client.conn.RemoteAddr().String(), message)
				return
			}
		}
		// If recipient is not found
		client.send <- fmt.Sprintf("[%s][Server]: Private message to %s not delivered as recipient is not connected.", time.Now().Format("15:04:05"), recipientName)
	}
}

func (cs *ChannelServer) WhoIsOnline(sender net.Conn) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	onlineUsers := "Online users:\n"

	for conn, client := range cs.clients {
		if conn != sender {
			onlineUsers += fmt.Sprintf("- %s\n", client.conn.RemoteAddr().String())
		}
	}

	// If there is only one user online (the sender), adjust message accordingly
	if onlineUsers == "Online users:\n" {
		onlineUsers += "- You are the only user online.\n"
	}

	sender.Write([]byte(fmt.Sprintf("\r\n%s", onlineUsers)))
}
