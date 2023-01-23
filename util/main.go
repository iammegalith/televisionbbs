package util

import (
	"net"
	"unicode/utf8"
)

// Structures
type Message struct {
	ID       int
	Basename string
	Subject  string
	Author   string
	Date     string
	Message  string
	Postto   string
}


// Global Variables
var LoggedInUsers = make(map[net.Conn]string)

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

// Global Functions
func TrimFirstChar(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
