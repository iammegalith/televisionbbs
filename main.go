package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gcfg.v1"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id       int64
	Username string
	Password string
	Level    int
}

type Config struct {
	Mainconfig struct {
		Port            int
		Bbsname         string
		Sysopname       string
		Prelogin        bool
		Bulletins       bool
		Newregistration bool
		Defaultlevel    int
		Configpath      string
		Ansipath        string
		Asciipath       string
		Modulepath      string
		Datapath        string
		Filespath       string
	}
}

const (
	BBS_VERSION = "1Q2023.1"
)

var (
	username        string = ""
	port            string = "8080"
	bbsname         string = "TelevisionBBS"
	sysopname       string = "Sysop"
	prelogin        bool   = true
	bulletins       bool   = true
	newregistration bool   = true
	defaultlevel    int    = 0
	configpath      string = "config/"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}


// called by: handleDoor(conn, "modules/door.exe")
func handleDoor(conn net.Conn, doorPath string) {
	cmd := exec.Command(doorPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		return
	}

	// Communicate with the "door" game
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				break
			}
			conn.Write(buf[:n])
		}
	}()

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		stdin.Write(buf[:n])
	}

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
}

func newUser(conn net.Conn, db *sql.DB) {
	var attempts int = 0
	for {
		fmt.Fprint(conn, "\r\nPlease choose a username: ")
		reader := bufio.NewReader(conn)
		for {
			char, _, err := reader.ReadRune()
			if err != nil {
				break
			}
			if char == '\r' || char == '\n' {
				break
			}
			fmt.Fprint(conn, string(char))
			username += string(char)
		}

		var user User
		fmt.Println("SELECT id FROM users WHERE username = ?", username)
		err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&user.Id)
		fmt.Println(err)
		if err != sql.ErrNoRows {
			attempts++
			if attempts >= 3 {
				fmt.Fprint(conn, "\r\nToo many attempts. Please try again later.\r\n")
				logout(conn)
				return
			}

			fmt.Fprint(conn, "\r\nUsername already exists, please choose another: ")
			username = ""
			continue
		}
		break
	}
	fmt.Println("Username: " + username)
	fmt.Fprint(conn, "\r\nPassword: ")
	passwordScanner := bufio.NewScanner(conn)
	passwordScanner.Scan()
	password := passwordScanner.Text()
	fmt.Fprint(conn, "\r\nRe-enter Password: ")
	rePasswordScanner := bufio.NewScanner(conn)
	rePasswordScanner.Scan()
	rePassword := rePasswordScanner.Text()

	if password != rePassword {
		fmt.Fprint(conn, "Passwords do not match. Please try again.\r\n")
		return
	}

	// You should hash the password and store the hashed password
	hashedPassword, _ := hashPassword(password)
	fmt.Println("INSERT INTO users(username,password,level) values('"+username+"','"+hashedPassword, 0)
	stmt, err := db.Prepare("INSERT INTO users(username,password,level) values(?,?,?)")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec(username, hashedPassword, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(conn, "Thank you for registering %s.\r\n", username)
}

func login(conn net.Conn, db *sql.DB) bool {
	fmt.Fprint(conn, "\033[2J\033[1;1H") // Clear screen
	fmt.Fprint(conn, "Welcome to the BBS. Are you a new user? (y/n): ")
	reader := bufio.NewReader(conn)
	isNewUser, _ := reader.ReadByte()
	if isNewUser == 'y' {
		newUser(conn, db)
		return true
	}
	fmt.Fprint(conn, "Username: ")
	var username string
	for {
		char, err := reader.ReadByte()
		if err != nil {
			break
		}
		if char == '\r' || char == '\n' {

			break
		}
		fmt.Fprint(conn, string(char))
		username += string(char)
	}

	fmt.Fprint(conn, "\r\nPassword: ")
	password, _, _ := bufio.NewReader(conn).ReadLine()

	var user User
	err := db.QueryRow("SELECT id, username, password, level FROM users WHERE username = ?", username).Scan(&user.Id, &user.Username, &user.Password, &user.Level)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprint(conn, "Incorrect username or password.\r\n")
			return false
		}
		fmt.Println(err)
		return false
	}
	// Compare the plain text password with the hashed password stored in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Fprint(conn, "Incorrect username or password.\r\n")
		return false
	}
	fmt.Fprintf(conn, "Welcome, %s.\r\n", user.Username)
	return true
}

//FIX THIS SHIT. IT'S AWFUL AND IT DOES NOT WORK.
func showAnsiFile(conn net.Conn, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(conn, "Error opening file: %s", err)
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		fmt.Fprintf(conn, "Error reading file: %s", err)
		return
	}

	// remove BOM from file
	content := bytes.TrimPrefix(buf.Bytes(), []byte("\xef\xbb\xbf"))

	// copy file contents to conn
	_, err = conn.Write(content)
	if err != nil {
		fmt.Fprintf(conn, "Error sending file: %s", err)
	}
}

func showTextFile(conn net.Conn, filePath string) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(conn, "Error: %v\n", err)
		return
	}
	defer file.Close()

	// Send the file contents to the user
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && (line[len(line)-1] == '\r' || line[len(line)-1] == '\n') {
			line = line[:len(line)-1]
		}
		fmt.Fprintln(conn, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(conn, "Error: %v\n", err)
	}
}

func logout(conn net.Conn) {
	fmt.Fprint(conn, "Goodbye!\r\n")
	username = ""
	conn.Close()
}

func handleConnection(conn net.Conn, db *sql.DB) {
	username = ""
	defer conn.Close()
	if login(conn, db) {
		// Handle authenticated user
		// Example: show main menu, handle commands, etc.
		// This is just for testing the login
		showTextFile(conn, "example.txt")
		showAnsiFile(conn, "example.ans")
		logout(conn)
		// end of testing
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./bbs.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	// Read Config File
	var cfg Config
	err = gcfg.ReadFileInto(&cfg, "bbs.config")
	if err != nil {
		fmt.Println("Failed to parse config file:", err)
		return
	}

	// set config values
	port = strconv.Itoa(cfg.Mainconfig.Port)
	bbsname = cfg.Mainconfig.Bbsname
	sysopname = cfg.Mainconfig.Sysopname
	prelogin = cfg.Mainconfig.Prelogin
	bulletins = cfg.Mainconfig.Bulletins
	newregistration = cfg.Mainconfig.Newregistration
	defaultlevel = cfg.Mainconfig.Defaultlevel
	configpath = cfg.Mainconfig.Configpath

	// start the listener
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	fmt.Println("Television BBS v" + BBS_VERSION + "\r\n" + bbsname + " Listening on TCP port " + port + "\r\nSysOp: " + sysopname + "\r\nConfig File Path: " + configpath)
	if newregistration {
		fmt.Println("New User Registration is enabled.")
	} else {
		fmt.Println("New User Registration is disabled.")
	}
	if bulletins {
		fmt.Println("Bulletins are enabled.")
	} else {
		fmt.Println("Bulletins are disabled.")
	}
	if prelogin {
		fmt.Println("Prelogin is enabled.")
	} else {
		fmt.Println("Prelogin is disabled.")
	}
	fmt.Println("Default userlevel: ", defaultlevel)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn, db)
	}
}
