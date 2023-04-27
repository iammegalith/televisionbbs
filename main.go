package main

import (
	"database/sql"
	"flag"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"thetelevision/bbsmenu"
	"thetelevision/channel"
	"thetelevision/onoff"
	"thetelevision/system"
	"time"

	"github.com/ebarkie/telnet"
	_ "github.com/mattn/go-sqlite3" // Import sqlite3 driver
)

var (
	channelServer *channel.ChannelServer
)

func HandleConnection(conn net.Conn, db *sql.DB, BBSConfig system.ConfigStruct) {
	defer conn.Close()
	defer log.Printf("Connection from %s closed", conn.RemoteAddr())

	tn := telnet.NewReadWriter(conn)

	tn.Write([]byte(".:: TeleVision BBS " + system.BBS_VERSION + " ::.\r\n\r\n"))
	if BBSConfig.PreLogin {
		system.ShowTextFile(tn, "ascii/prelogin.asc")
	}

	u, err := onoff.Login(db, tn, conn)
	if err != nil {
		log.Printf("Error logging in: %v", err)
	}
	// Create an example User struct. You should replace this with the actual user data after authentication.
	user := u

    bbsmenu.GetMenu(tn, conn, channelServer, user, "main")
}

func main() {
	daemonize := flag.Bool("d", false, "run as daemon")
	flag.Parse()

	if *daemonize {
		cmd := exec.Command(os.Args[0])
		err := cmd.Start()
		if err != nil {
			log.Fatalf("Error starting daemon process: %v", err)
		}
		log.Println("TeleVision BBS started as a daemon")
		os.Exit(0)
	}
	
	db, err := sql.Open("sqlite3", "data/userdata.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}
	system.ReadBBSConfig()

	// Set up the connection pool
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	// Set up the channel server
	channelServer = channel.NewChannelServer()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Error starting Channel Server:", err)
			}
		}()
		err := channelServer.ListenAndServe()
		if err != nil {
			log.Println("Error starting Channel Server:", err)
		}
	}()

	addr := net.JoinHostPort(system.BBSConfig.ListenAddr, strconv.Itoa(system.BBSConfig.Port))
	log.Printf("TeleVision BBS listening on %s", addr)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// Ping the database periodically to check if it's still available
	go func() {
		for {
			if err := db.Ping(); err != nil {
				log.Printf("Error pinging database: %v", err)
			}
			time.Sleep(time.Minute)
		}
	}()

	// Set up a waitgroup to keep track of all active connections
	wg := &sync.WaitGroup{}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
			}

		log.Printf("Accepted connection from %s", conn.RemoteAddr())

		// Add 1 to the waitgroup for the new connection
		wg.Add(1)

		go func(conn net.Conn) {
			defer wg.Done()

			// Add the new connection to the channel server
			channelServer.AddClient(conn)

			HandleConnection(conn, db, system.BBSConfig)

			// Remove the connection from the channel server after the user disconnects
			channelServer.RemoveClient(conn)
		}(conn)
	}
}
