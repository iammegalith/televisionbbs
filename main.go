package main

import (
	// Core packages
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// Specific to TeleVision BBS
	"televisionbbs/util"

	// External Packages
	"github.com/PatrickRudolph/telnet"
	"github.com/PatrickRudolph/telnet/options"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gcfg.v1"
)

// Constants
const (
	TELEVISION_VERSION = "1Q2023.1.31"
)

// Structs
// Television Config
type Config struct {
	Mainconfig struct {
		Port            string
		Listenaddr      string
		Bbsname         string
		Sysopname       string
		Prelogin        bool
		Showbulls       bool
		Newregistration bool
		Defaultlevel    int
		Configpath      string
		Ansipath        string
		Asciipath       string
		Modulepath      string
		Datapath        string
		Filespath       string
		StringsFile     string
		Textfiles       string
	}
}

// Television String
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
		Userlisthead       string
		Userlistfoot       string
	}
}

// Current User Struct
type CurrentUser struct {
	username      string
	conn          *telnet.Connection
	userStatus    string
	ansiSupported bool
	hasAnsi       bool
	level         int
	linefeeds     int
	hotkeys       int
	active        int
	clearscreen   int
}

// Global Variables
var (
	// Core Functionality

	TheUser  = make(map[string]CurrentUser)
	username string
	currmenu string
	prevmenu string
	hasAnsi  bool

	// Television Config
	port         string
	listenaddr   string
	bbsname      string
	sysopname    string
	prelogin     bool
	showbulls    bool
	newreg       bool
	defaultlevel int
	configpath   string
	ansipath     string
	asciipath    string
	modulepath   string
	datapath     string
	filespath    string
	stringsfile  string
	textfiles    string

	// Strings
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
	userlisthead       string = "User List"
	userlistfoot       string = "========="
)

// Initialization sequences
func init() {
	util.Connections = make(map[string]*telnet.Connection)
	util.UserStatus = make(map[string]string)
}

// Read BBS Config File
//
// Call with:
//
// getConfig()
func getConfig() {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "config/television.conf")
	if err != nil {
		fmt.Println("Failed to parse config file:", err)
		return
	} else {
		port = cfg.Mainconfig.Port
		listenaddr = cfg.Mainconfig.Listenaddr
		bbsname = cfg.Mainconfig.Bbsname
		sysopname = cfg.Mainconfig.Sysopname
		prelogin = cfg.Mainconfig.Prelogin
		showbulls = cfg.Mainconfig.Showbulls
		newreg = cfg.Mainconfig.Newregistration
		defaultlevel = cfg.Mainconfig.Defaultlevel
		configpath = cfg.Mainconfig.Configpath
		ansipath = cfg.Mainconfig.Ansipath
		asciipath = cfg.Mainconfig.Asciipath
		modulepath = cfg.Mainconfig.Modulepath
		datapath = cfg.Mainconfig.Datapath
		filespath = cfg.Mainconfig.Filespath
		stringsfile = cfg.Mainconfig.StringsFile
		textfiles = cfg.Mainconfig.Textfiles
	}
}

func getStrings() {
	var bbsStrings BbsStrings
	err := gcfg.ReadFileInto(&bbsStrings, configpath+stringsfile)
	if err != nil {
		fmt.Println("Failed to parse strings file:", err)
		return
	} else {
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
		userlisthead = bbsStrings.General.Userlisthead
		userlistfoot = bbsStrings.General.Userlistfoot
	}
}

// Converts error to a byte string
//
// Call with:
//
// writeLine(w, []byte(convertError(err)))
func convertError(err error) (errorString string) {
	errMessage := fmt.Sprintf("%s", err)
	return (errMessage)
}

// Converts error to a byte string
// Call with: writeLine(w, []byte(convertByteToString(strValue)))
func convertByteToString(byteThing byte) (strValue string) {
	strValue = string([]byte{byteThing})
	return (strValue)
}

// This is used to read a single keypress from the user
// called with:
//
//	reader := telnet.Reader(conn)
//	hotkey, err := readKey(reader)
func readKey(conn *telnet.Connection) (string, error) {
	char := make([]byte, 1)
	_, err := conn.Read(char)
	if err != nil {
		return "error", err
	}
	returnKey := convertByteToString(char[0])
	return returnKey, nil
}

// This is used to read the user's input
// called with:
func readLine(conn *telnet.Connection) (string, error) {
	reader := bufio.NewReader(conn)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(line)), nil
}

func writeLine(conn *telnet.Connection, buffer []byte) error {
	_, err := conn.Write(buffer)
	if err != nil {
		return err
	}
	return nil
}

func toggleAnsi(conn *telnet.Connection) bool {
	if hasAnsi {
		hasAnsi = false
		writeLine(conn, []byte("ANSI mode is OFF"))
		return hasAnsi
	} else {
		hasAnsi = true
		writeLine(conn, []byte("ANSI mode is "+util.ANSI_BOLD+util.ANSI_BLUE_BG+util.ANSI_WHITE+"ON"+util.ANSI_RESET))
		return hasAnsi
	}
}

func listUsers(conn *telnet.Connection) {
	util.LockMaps()
	defer util.UnlockMaps()
	fmt.Println("Connections:")
	for name, conn := range util.Connections {
		fmt.Printf("%s : %v\n", name, conn)
	}
	fmt.Println("UserStatus:")
	for name, status := range util.UserStatus {
		fmt.Printf("%s : %v\n", name, status)
	}
}

// call with:
//
//	answer, err := askYesNo(conn, "Are your sure?")
//
// if answer { this is code for yes }
//
// or
//
// if !answer { this is code for no }
func askYesNo(conn *telnet.Connection, question string) (bool, error) {
	var response string
	writeLine(conn, []byte(question+" (y/n): "))
	hotkey, err := readKey(conn)
	if err != nil {
		writeLine(conn, []byte(convertError(err)))
	}
	response = strings.ToLower(hotkey)
	if response == "y" || response == "yes" {
		return true, nil
	} else if response == "n" || response == "no" {
		return false, nil
	} else {
		return false, fmt.Errorf("\r\ninvalid response. Please enter 'y' or 'n'")
	}
}

func logout(conn *telnet.Connection) {
	util.UserStatus[util.Cuname] = "logging out"
	writeLine(conn, []byte("Goodbye!\r\n"))
	conn.Close()
}

// pressKey function reads a single keypress and returns an err or nothing
//
// Call with:
//
// pressKey(*conn, "Press a key")
func pressKey(conn *telnet.Connection, message string) error {
	writeLine(conn, []byte("\r\n"+message))
	_, err := conn.Read(make([]byte, 1))
	if err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func newUser(conn *telnet.Connection, db *sql.DB) {
	var attempts int = 0
	var username string

	for {
		writeLine(conn, []byte("\r\nPlease choose a username: "))
		username, _ = readLine(conn)
		var id int
		err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&id)
		fmt.Println("error: ", err)
		if err != sql.ErrNoRows {
			attempts++
			if attempts >= 3 {
				writeLine(conn, []byte("\r\nToo many attempts. Please try again later.\r\n"))
				logout(conn)
				return
			}

			writeLine(conn, []byte("\r\nUsername already exists, please choose another: "))
			continue
		}
		break
	}
	writeLine(conn, []byte("\r\nPassword: "))
	passOne, _ := readLine(conn)
	writeLine(conn, []byte("\r\nRe-enter Password: "))
	passTwo, _ := readLine(conn)
	if passOne != passTwo {
		writeLine(conn, []byte("Passwords do not match. Please try again.\r\n"))
		return
	}

	// You should hash the password and store the hashed password
	hashedPassword, _ := hashPassword(passOne)
	stmt, err := db.Prepare("INSERT INTO users(username,password,level,linefeeds,hotkeys,ansi, active, clearscreen) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = stmt.Exec(username, hashedPassword, 0, 0, 0, 0, 1, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	writeLine(conn, []byte("Thank you for registering "+username+".\r\n"))
}

func checkAttempts(conn *telnet.Connection, attempts int) {
	if attempts >= 3 {
		writeLine(conn, []byte("\r\nToo many attempts. Please try again later.\r\n"))
		util.UserStatus[util.Cuname] = "too many attempts"
		logout(conn)
	}
}

func login(conn *telnet.Connection, db *sql.DB, TheUser map[string]CurrentUser) bool {
	// check is prelogin is true, if so, print prelogin ascii message.
	if prelogin {
		showTextFile(conn, textfiles+"prelogin.txt")
	}
	// Check if ANSI is supported - set variable ansiSupported
	ansiSupported := checkForANSI(conn)
	if ansiSupported {
		writeLine(conn, []byte("\r\n"+util.ANSI_BOLD+util.ANSI_BLUE_BG+util.ANSI_WHITE+"ANSI Supported"+util.ANSI_RESET+"\r\n"))
	} else {
		writeLine(conn, []byte("\r\nASCII Mode\r\n"))
	}

	// reset login attempt counter to zero
	var attempts int = 0
	// this needs to be part of the username prompt - if new, type new, otherwise, use your username
	answer, _ := askYesNo(conn, "Are you a new user? ")
	if answer {
		writeLine(conn, []byte("\r\n"))
		newUser(conn, db)
		return true
	}

	// get username, set string var username to the username. one more "username" :)
	writeLine(conn, []byte("\r\n"))
	writeLine(conn, []byte("Username: "))
	var username string
	getUsername, _ := readLine(conn)
	username = strings.TrimSpace(getUsername)
	if username == "" {
		writeLine(conn, []byte("\r\nUsername cannot be blank."))
		attempts++
		checkAttempts(conn, attempts)
		return login(conn, db, TheUser)
	}
	// get the users password - set var password to the password. one more "password" :)
	writeLine(conn, []byte("\r\nPassword: "))
	password, _ := readLine(conn)

	// Check if the username exists in the database - populate the user struct with the users info
	//set temp vars for the query
	tUsername := ""
	tId := 0
	tPassword := ""

	err := db.QueryRow("SELECT id, username, password from users WHERE username = ?", username).Scan(tId, tUsername, tPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			attempts++
			checkAttempts(conn, attempts)
			writeLine(conn, []byte("\r\n"))
			writeLine(conn, []byte("Incorrect username or password.\r\n"))
			return false
		}
		fmt.Println(err)
		return false
	}
	// Compare the plain text password with the hashed password stored in the database
	err = bcrypt.CompareHashAndPassword([]byte(tPassword), []byte(strings.TrimSpace(password)))
	if err != nil {
		attempts++
		checkAttempts(conn, attempts)
		writeLine(conn, []byte("\r\n"))
		writeLine(conn, []byte("Incorrect username or password.\r\n"))
		return false
	}
	if alreadyLoggedIn(conn, username, db) {
		writeLine(conn, []byte("\r\nThis account is already logged in. Please try again later.\r\n"))
		return false
	}

	// go ahead and the the rest of the data, set the map.
	tUsername = ""
	tId = 0
	tPassword = ""
	tLevel := 0
	tLinefeeds := 0
	tHotkeys := 0
	tActive := 0
	tClearscreen := 0

	err = db.QueryRow("SELECT id, username, password from users WHERE username = ?", username).Scan(tId, tUsername, tPassword, tLevel, tLinefeeds, tHotkeys, tActive, tClearscreen)
	if err != nil {
		fmt.Println("error getting user data", err)
		return false
	}
	util.LockMaps()
	currentuser := TheUser[username]
	currentuser.username = username
	currentuser.level = tLevel
	currentuser.linefeeds = tLinefeeds
	currentuser.hotkeys = tHotkeys
	currentuser.active = tActive
	currentuser.clearscreen = tClearscreen
	currentuser.ansiSupported = ansiSupported
	currentuser.hasAnsi = ansiSupported
	currentuser.conn = conn
	currentuser.userStatus = "logged in"
	util.UnlockMaps()
	writeLine(conn, []byte("\r\n"))
	writeLine(conn, []byte("Welcome, \r\n"+util.Cuname))
	return true
}

func alreadyLoggedIn(conn *telnet.Connection, username string, db *sql.DB) bool {

	var sessionExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM sessions WHERE username = ? and active = 1)", username).Scan(&sessionExists)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if sessionExists {
		writeLine(conn, []byte("You are already logged in. Would you like to end your previous session? (y/n): "))
		buffer := make([]byte, 1)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return false
		}
		if n > 0 {
			endSession := buffer[0]
			if endSession == 'y' {
				// end previous session
				_, err := db.Exec("UPDATE sessions SET active = 0 WHERE username = ?", username)
				if err != nil {
					fmt.Println(err)
					return false
				}
				return false
			}
		} else {
			return true
		}
	}
	return false
}

func menu(conn *telnet.Connection, db *sql.DB) {
	for {
		util.UpdateUserStatus(util.Cuname, "menu")
		fmt.Fprint(conn, "\r\n[a]NSI ON/OFF\r\n")
		fmt.Fprint(conn, "[w]ho is online\r\n")
		fmt.Fprint(conn, "[i]nfo\r\n")
		fmt.Fprint(conn, "[g]oodbye\r\n")
		status := util.GetUserStatus(util.Cuname)
		showprompt := util.ShowPrompt(status, menuprompt, hasAnsi)
		fmt.Fprint(conn, showprompt)
		text, err := readKey(conn)
		if err != nil {
			returnError := convertError(err)
			fmt.Println("error reading input: ", returnError)
			logout(conn)
			return
		}
		switch text {
		case "a":
			hasAnsi = toggleAnsi(conn)
		case "w":
			listUsers(conn)
		case "i":
			showTextFile(conn, textfiles+"info.txt")
		case "g":
			answer, err := askYesNo(conn, "\r\nAre you sure you want to logout?")
			if err != nil {
				returnError := convertError(err)
				fmt.Println("error reading input: ", returnError)
				logout(conn)
				return
			}
			if answer {
				logout(conn)
				return
			} else {
				continue
			}
		default:
			fmt.Fprintln(conn, "\r\nInvalid option. Please try again.")
		}
	}
}

// This is used to print a text file to the user
// called with:
//
//     filePath := "./example.txt"
//     err := showTextFile(conn, filePath)

func showTextFile(conn *telnet.Connection, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	pageSize := 20 // number of lines to show per page
	lineCount := 0 // current line count
	for scanner.Scan() {
		// word wrap text
		words := strings.Split(scanner.Text(), " ")
		var line string
		for _, word := range words {
			if len(line)+len(word) > 80 {
				fmt.Fprintln(conn, line)
				lineCount++
				if lineCount%pageSize == 0 {
					fmt.Fprint(conn, "\r\nPress any key to continue or 'q' to quit...\r\n")
					text, _ := readKey(conn)
					if text == "q" {
						return nil
					}
					lineCount = 0
				}
				line = ""
			}
			line += word + " "
		}
		fmt.Fprintln(conn, line+"\r")
		lineCount++
		if lineCount%pageSize == 0 {
			fmt.Fprint(conn, "\r\nPress any key to continue or 'q' to quit...\r\n")
			text, _ := readKey(conn)
			if text == "q" {
				return nil
			}
			lineCount = 0
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// check for ANSI negotiation
func checkForANSI(conn *telnet.Connection) bool {
	doesANSI, _ := askYesNo(conn, "Does your terminal support ANSI?  ")
	return doesANSI
}

//
// Three primary functions here: handleConnection, should only be called by func main(), checkDissconnectedUsers, and func main.
//

func handleConnection(conn *telnet.Connection, TheUser map[string]CurrentUser) {
	// Open the database
	db, err := sql.Open("sqlite3", datapath+"bbs.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	sendline := []byte(welcomestring + "\r\nVersion " + TELEVISION_VERSION + "\r\n")
	writeLine(conn, sendline)
	login(conn, db, TheUser)
	menu(conn, db)
	util.LockMaps()
	delete(util.Connections, util.Cuname)
	delete(util.UserStatus, util.Cuname)
	util.UnlockMaps()
}

func checkDisconnectedClients() {
	for {
		time.Sleep(5 * time.Second)
		util.LockMaps()
		for cuname, conn := range util.Connections {
			if conn.RemoteAddr() == nil {
				delete(util.Connections, cuname)
				delete(util.UserStatus, cuname)
			}
		}
		util.UnlockMaps()
	}
}

func main() {
	getConfig()
	getStrings()

	// start the routine to check for disconnected clients
	go checkDisconnectedClients()
	svr := telnet.NewServer(listenaddr+":"+port, telnet.HandleFunc(func(conn *telnet.Connection) {
		handleConnection(conn, TheUser)
	}), options.NAWSOption)
	err := svr.ListenAndServe()
	if err != nil {
		log.Fatalln("Error starting telnet server:", err)
	}
}
