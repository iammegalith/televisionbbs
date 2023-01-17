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

	"github.com/gdamore/tcell"
	"github.com/go-ini/ini"
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
		StringsFile     string
	}
}

type BbsStrings struct {
	General struct {
		Menuprompt         string
		Welcomestring      string
		Pressanykey        string
		Pressreturn        string
		Entername          string
		Enterpassword      string
		Enternewpassword   string
		Enterpasswordagain string
		Areyounew          string
		Areyousure         string
		Ansimode           string
		Asciimode          string
		Invalidoption      string
		Invalidname        string
		Invalidpassword    string
		Passwordmismatch   string
		Userexists         string
		Usercreated        string
		Enterchat          string
		Leavechat          string
		Pagesysop          string
		Ispagingyou        string
		Sysopishere        string
		Sysopisaway        string
		Sysopisbusy        string
	}
}

const (
	BBS_VERSION = "1Q2023.1"
)

var (
	username           string = ""
	port               string = "8080"
	bbsname            string = "TelevisionBBS"
	sysopname          string = "Sysop"
	prelogin           bool   = true
	bulletins          bool   = true
	newregistration    bool   = true
	defaultlevel       int    = 0
	configpath         string = "config/"
	stringsfile        string = "strings.config"
	ansipath           string = "ansi/"
	asciipath          string = "ascii/"
	modulepath         string = "modules/"
	datapath           string = "data/"
	filespath          string = "files/"
	menuprompt         string = "Please select an option:"
	welcomestring      string = "Welcome to TelevisionBBS!"
	pressanykey        string = "Press any key to continue..."
	pressreturn        string = "Press return to continue..."
	entername          string = "Please enter your name:"
	enterpassword      string = "Please enter your password:"
	enternewpassword   string = "Please enter a new password:"
	enterpasswordagain string = "Please enter your password again:"
	areyounew          string = "Are you new to TelevisionBBS? (Y/N)"
	areyousure         string = "Are you sure? (Y/N)"
	ansimode           string = "ANSI mode"
	asciimode          string = "ASCII mode"
	invalidoption      string = "Invalid option."
	invalidname        string = "Invalid name."
	invalidpassword    string = "Invalid password."
	passwordmismatch   string = "Passwords do not match."
	userexists         string = "User already exists."
	usercreated        string = "User created."
	enterchat          string = "Enter chat"
	leavechat          string = "Leave chat"
	pagesysop          string = "Page Sysop"
	ispagingyou        string = " is paging you."
	sysopishere        string = "Sysop is here."
	sysopisaway        string = "Sysop is away."
	sysopisbusy        string = "Sysop is busy."
)

// General Command and Functions

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
				logout(conn, user)
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

// FIX THIS SHIT. IT'S AWFUL AND IT DOES NOT WORK.
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

func logout(conn net.Conn, user User) {
	fmt.Fprint(conn, "Goodbye!\r\n")
	username = ""
	conn.Close()
}

// Menu System
func getMenu(conn net.Conn, user User, menuName string) {
	var err error
	var mcfg *ini.File
	mcfg, err = ini.Load(fmt.Sprintf("%s.ini", menuName))
	if err != nil {
		fmt.Printf("Error loading config file: %v", err)
		return
	}

	// Get the menu options from the INI file
	options := mcfg.Sections()

	// Display the menu
	fmt.Println("Please select an option:")
	for i, option := range options {
		cmd := option.Key("command").String()
		desc := option.Key("description").String()
		fmt.Printf("%d) %s - %s\n", i+1, cmd, desc)
	}

	// Get user input for menu selection
	var selection int
	fmt.Scan(&selection)

	// Get the selected menu option
	selectedOption := options[selection-1]
	typ := selectedOption.Key("type").String()
	args := selectedOption.Key("arguments").String()
	ilvl := selectedOption.Key("level").String()
	lvl, _ := strconv.Atoi(ilvl)

	if typ == "submenu" {
		getMenu(conn, user, menuName)
	} else {
		handleSelection(conn, user, typ, args, lvl)
	}
}

func handleSelection(conn net.Conn, user User, typ string, args string, lvl int) {
	if lvl > user.Level {
		fmt.Println("Sorry, you don't have permission to access this feature.")
	} else {
		switch typ {
		case "function":
			switch args {
			case "bulletins":
				getMenu(conn, user, "bulletins")
			case "goodbye":
				logout(conn, user)
				break
			default:
				fmt.Println("Invalid option selected")
			}
			break
		case "display":
			switch args {
			case "info":
				showTextFile(conn, args)
			default:
				fmt.Println("Invalid option selected")
			}
			break
		default:
			fmt.Println("Invalid option selected")
		}
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func statusScreen(screen tcell.Screen, err error) {
	// Clear the screen
	screen.Clear()

	// Draw the status screen
	x, y := screen.Size()
	header := "TeleVision BBS Version " + BBS_VERSION
	headerX := (x / 2) - (len(header) / 2)
	drawText(screen, headerX, 0, x, y, tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack), header)

	// Draw the status information
	// Example: number of users online, number of messages, etc.
	usersOnline := "Users online: 12"
	numMessages := "Messages: 42"
	drawText(screen, headerX, 1, x, y, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack), usersOnline)
	drawText(screen, headerX, 2, x, y, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack), numMessages)
	// Draw the footer
	footer := "Press 'q' to quit"
	footerX := (x / 2) - (len(footer) / 2)
	drawText(screen, footerX, y-1, x, y, tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack), footer)

	// Show the screen
	screen.Show()

	// Wait for user input
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyRune:
				if ev.Rune() == 'q' {
					screen.Fini()
					return
				}
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}

// Handle Connection and Main functions
func handleConnection(conn net.Conn, db *sql.DB) (user User) {
	username = ""
	defer conn.Close()
	if login(conn, db) {
		getMenu(conn, user, "main")
		logout(conn, user)
		// end of testing
	}
	return
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
	stringsfile = cfg.Mainconfig.StringsFile

	// Get strings values
	var bbsStrings BbsStrings
	err = gcfg.ReadFileInto(&bbsStrings, configpath+stringsfile)
	if err != nil {
		fmt.Println("Failed to parse strings file:", err)
		return
	}

	// set strings values
	menuprompt = bbsStrings.General.Menuprompt
	welcomestring = bbsStrings.General.Welcomestring
	pressanykey = bbsStrings.General.Pressanykey
	pressreturn = bbsStrings.General.Pressreturn
	entername = bbsStrings.General.Entername
	enterpassword = bbsStrings.General.Enterpassword
	enternewpassword = bbsStrings.General.Enternewpassword
	enterpasswordagain = bbsStrings.General.Enterpasswordagain
	areyounew = bbsStrings.General.Areyounew
	areyousure = bbsStrings.General.Areyousure
	ansimode = bbsStrings.General.Ansimode
	asciimode = bbsStrings.General.Asciimode
	invalidoption = bbsStrings.General.Invalidoption
	invalidname = bbsStrings.General.Invalidname
	invalidpassword = bbsStrings.General.Invalidpassword
	passwordmismatch = bbsStrings.General.Passwordmismatch
	userexists = bbsStrings.General.Userexists
	usercreated = bbsStrings.General.Usercreated
	enterchat = bbsStrings.General.Enterchat
	leavechat = bbsStrings.General.Leavechat
	pagesysop = bbsStrings.General.Pagesysop
	ispagingyou = bbsStrings.General.Ispagingyou
	sysopishere = bbsStrings.General.Sysopishere
	sysopisaway = bbsStrings.General.Sysopisaway
	sysopisbusy = bbsStrings.General.Sysopisbusy

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
		screen, err := tcell.NewScreen()
		go statusScreen(screen, err)
		go handleConnection(conn, db)
	}
}
