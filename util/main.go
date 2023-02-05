package util

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/PatrickRudolph/telnet"
)

// User Database Structure
type UserStruct struct {
	Id          int64
	Username    string
	Password    string
	Level       int
	Linefeeds   int
	Hotkeys     bool
	Active      bool
	Clearscreen bool
}


// Global Variables
var Connections map[string]*telnet.Connection
var Mutex sync.Mutex
var UserStatus map[string]string

var (
	// User Database
	User          UserStruct
	Cuid          int64
	Cuname        string
	Culevel       int
	Culinefeeds   int
	Cuhotkeys     bool
	Cuactive      bool
	Cuclearscreen bool
)

const (
	ANSI_RESET             = "\x1b[0m"
	ANSI_BOLD              = "\x1b[1m"
	ANSI_BLACK             = "\x1b[30m"
	ANSI_RED               = "\x1b[31m"
	ANSI_GREEN             = "\x1b[32m"
	ANSI_YELLOW            = "\x1b[33m"
	ANSI_BLUE              = "\x1b[34m"
	ANSI_MAGENTA           = "\x1b[35m"
	ANSI_CYAN              = "\x1b[36m"
	ANSI_WHITE             = "\x1b[37m"
	ANSI_BLACK_BG          = "\x1b[40m"
	ANSI_RED_BG            = "\x1b[41m"
	ANSI_GREEN_BG          = "\x1b[42m"
	ANSI_YELLOW_BG         = "\x1b[43m"
	ANSI_BLUE_BG           = "\x1b[44m"
	ANSI_MAGENTA_BG        = "\x1b[45m"
	ANSI_CYAN_BG           = "\x1b[46m"
	ANSI_WHITE_BG          = "\x1b[47m"
	ANSI_CLR               = "\x1b[2J"
	ANSI_HOME              = "\x1b[H"
	ANSI_UP                = "\x1b[1A"
	ANSI_DOWN              = "\x1b[1B"
	ANSI_RIGHT             = "\x1b[1C"
	ANSI_LEFT              = "\x1b[1D"
	ANSI_CLEAR_LINE        = "\x1b[2K"
	ANSI_CLEAR_LINE_UP     = "\x1b[1K"
	ANSI_CLEAR_LINE_DOWN   = "\x1b[0K"
	ANSI_CLEAR_SCREEN      = "\x1b[2J"
	ANSI_CLEAR_SCREEN_UP   = "\x1b[1J"
	ANSI_CLEAR_SCREEN_DOWN = "\x1b[0J"
	ANSI_SAVE_CURSOR       = "\x1b[s"
	ANSI_RESTORE_CURSOR    = "\x1b[u"
	ANSI_CURSOR_HOME       = "\x1b[H"
	CR_LF                  = "\r\n"
)

func UpdateUserStatus(username, status string) {
	Mutex.Lock()
	UserStatus[Cuname] = status
	Mutex.Unlock()
}

func GetUserStatus(username string) string {
	Mutex.Lock()
	status := UserStatus[Cuname]
	Mutex.Unlock()
	return status
}

func LockMaps() {
	Mutex.Lock()
}

func UnlockMaps() {
	Mutex.Unlock()
}

func GenRandomString(length int) string {
	chars := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	result := make([]rune, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func ShowPrompt(status, menuprompt string, Cuansi bool) (showprompt string) {
	LockMaps()
	defer UnlockMaps()
	fmt.Println("Cuansi: ", Cuansi)
	if Cuansi {
		fmt.Println("ANSI ENABLED")
		showprompt := (ANSI_BOLD + ANSI_WHITE + status + ANSI_RESET + ":[" + ANSI_BOLD + ANSI_GREEN + menuprompt + ANSI_RESET + " ]: ")
		return showprompt
	} else {
		showprompt := (status + ":[" + menuprompt + " ]: ")
		return showprompt
	}
}
