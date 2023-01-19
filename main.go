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
	"strings"

	"github.com/go-ini/ini"
	"github.com/k0kubun/go-ansi"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gcfg.v1"
)

type User struct {
	Id          int64
	Username    string
	Password    string
	Level       int
	Linefeeds   int
	Translation int
	Active      bool
	Clearscreen bool
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
	prevmenu           string = ""
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

// user data
var (
	user          User
	cuid          int64
	cuname        string
	culevel       int
	culinefeeds   int
	cutranslation int
	cuactive      bool
	cuclearscreen bool
)

// General Command and Functions

func alreadyLoggedIn(conn net.Conn, username string, db *sql.DB) bool {
	var sessionExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM sessions WHERE username = ? and active = 1)", username).Scan(&sessionExists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if sessionExists {
		reader := bufio.NewReader(conn)
		fmt.Fprint(conn, "\r\nYou are already logged in. Would you like to end your previous session? (y/n): ")
		endSession, _ := reader.ReadByte()
		if endSession == 'y' {
			// end previous session
			_, err := db.Exec("UPDATE sessions SET active = 0 WHERE username = ?", username)
			if err != nil {
				fmt.Println(err)
				return false
			}
			return false
		} else {
			return true
		}
	}
	return false
}

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
	stmt, err := db.Prepare("INSERT INTO users(username,password,level,linefeeds,translation, active, clearscreen) values(?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec(username, hashedPassword, 0, 0, 0, 1, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Fprintf(conn, "Thank you for registering %s.\r\n", username)
}

func login(conn net.Conn, db *sql.DB) bool {
	fmt.Fprint(conn, "\r\nWelcome to the BBS. Are you a new user? (y/n): ")
	reader := bufio.NewReader(conn)
	isNewUser, _ := reader.ReadByte()
	if isNewUser == 'y' {
		fmt.Fprintf(conn, "\r\n")
		newUser(conn, db)
		return true
	}
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprint(conn, "Username: ")
	var username string
	for {
		char, err := reader.ReadByte()
		fmt.Fprint(conn, string(char))
		if err != nil {
			break
		}
		if char == '\r' || char == '\n' {
			break
		}
		username += string(char)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		fmt.Fprintln(conn, "\r\nUsername cannot be blank.")
		return login(conn, db)
	}
	fmt.Fprint(conn, "\r\nPassword: ")
	password, _, _ := bufio.NewReader(conn).ReadLine()
	err := db.QueryRow("SELECT id, username, password, level, linefeeds, translation, active, clearscreen FROM users WHERE username = ?", username).Scan(&user.Id, &user.Username, &user.Password, &user.Level, &user.Linefeeds, &user.Translation, &user.Active, &user.Clearscreen)

	cuid = user.Id
	cuname = user.Username
	culevel = user.Level
	culinefeeds = user.Linefeeds
	cutranslation = user.Translation
	cuactive = user.Active
	cuclearscreen = user.Clearscreen

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(conn, "\r\n")
			fmt.Fprint(conn, "Incorrect username or password.\r\n")
			return false
		}
		fmt.Println(err)
		return false
	}
	// Compare the plain text password with the hashed password stored in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Fprintf(conn, "\r\n")
		fmt.Fprint(conn, "Incorrect username or password.\r\n")
		return false
	}
	if alreadyLoggedIn(conn, username, db) {
		fmt.Fprint(conn, "\r\nThis account is already logged in. Please try again later.\r\n")
		return false
	}
	fmt.Fprintf(conn, "\r\n")
	fmt.Fprintf(conn, "Welcome, %s.\r\n", cuname)
	return true
}

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

	// Interpret ANSI codes
	ansiContent, err := ansi.Printf(string(content))
	if err != nil {
		fmt.Fprintf(conn, "Error interpreting ANSI codes: %s", err)
		return
	}
	fmt.Fprint(conn, ansiContent)
}

func showTextFile(conn net.Conn, filePath string, linefeeds int) {
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
		if linefeeds == 1 {
			fmt.Println("Linefeeds: 1")
			fmt.Fprint(conn, line+"\r\n")
		} else {
			fmt.Fprint(conn, line)
		}
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
	var currentMenu string = menuName
	var err error
	var mcfg *ini.File
	mcfg, err = ini.Load(fmt.Sprintf(configpath+"%s.config", menuName))
	if err != nil {
		fmt.Printf("Error loading config file: %v", err)
		return
	}

	// Get the menu options from the INI file
	options := mcfg.Sections()

	// Display the menu
	showTextFile(conn, fmt.Sprintf(asciipath+"%s.asc", menuName), culinefeeds)
	fmt.Fprint(conn, menuprompt)

	// Get user input for menu selection
	user_input, _ := bufio.NewReader(conn).ReadString('\n')
	user_input = strings.TrimSpace(user_input)

	// Get the selected menu option
	var selection *ini.Section
	for _, option := range options {
		if option.Key("fast").String() == user_input {
			selection = option
			break
		}
	}

	if selection == nil {
		fmt.Fprint(conn, "\r\nInvalid selection")
		getMenu(conn, user, currentMenu)
		return
	}

	typ := selection.Key("type").String()
	args := selection.Key("arguments").String()
	lvl := selection.Key("level").MustInt(0)

	if typ == "menu" {
		prevmenu = currentMenu
		getMenu(conn, user, args)
	} else {
		handleSelection(conn, user, args, lvl, currentMenu)
	}
}

func pressReturn(conn net.Conn) error {
	fmt.Fprint(conn, "\r\n--- Press [ RETURN ]] ---")
	// Wait for the user to press RETURN
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	return nil
}

func pressKey(conn net.Conn) error {
	fmt.Fprint(conn, "\r\n--- Press any key ---")
	_, err := conn.Read(make([]byte, 1))
	if err != nil {
		return err
	}
	return nil
}

func askYesNo(conn net.Conn, question string) (bool, error) {
	var response string
	fmt.Fprint(conn, question+" (y/n): ")
	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		response = scanner.Text()
	} else {
		return false, scanner.Err()
	}
	response = strings.ToLower(response)
	if response == "y" || response == "yes" {
		return true, nil
	} else if response == "n" || response == "no" {
		return false, nil
	} else {
		return false, fmt.Errorf("\r\nInvalid response. Please enter 'y' or 'n'.")
	}
}

func handleSelection(conn net.Conn, user User, args string, lvl int, currentMenu string) {
	if lvl > user.Level {
		fmt.Fprintf(conn, "Sorry, you don't have permission to access this feature.")
	} else {
		switch args {
		case "info":
			showTextFile(conn, asciipath+args+".asc", culinefeeds)
			pressKey(conn)
			getMenu(conn, user, currentMenu)
		case "goodbye":
			result, err := askYesNo(conn, "Are you sure you want to log out?")
			if err != nil {
				fmt.Fprintf(conn, "Error: %v", err)
				return
			}
			if result {
				logout(conn, user)
			} else {
				getMenu(conn, user, currentMenu)
			}
		case "bye":
			logout(conn, user)
		case "teleconference":
			// code to handle teleconference feature
		case "userlist":
			// code to handle userlist feature
		case "obbs":
			// code to handle obbs feature
		case "whosonline":
			// code to handle whosonline feature
		case "pagesysop":
			// code to handle pagesysop feature
		case "userconfig":
			// code to handle userconfig feature
		case "sysop":
			// code to handle sysop feature
		default:
			if args == "" {
				getMenu(conn, user, currentMenu)
				return
			} else {
				fmt.Fprintf(conn, "\r\nInvalid option selected")
				pressKey(conn)
				getMenu(conn, user, currentMenu)
				return
			}
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
		go handleConnection(conn, db)
	}
}
