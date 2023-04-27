package system

import (
	"errors"
	"time"

	"github.com/ebarkie/telnet"
)

// Globals
var ErrUserNotFound = errors.New("user not found")

// Structs
type UserMap struct {
	Username string
	Ctx      *telnet.Ctx
	Status   string
	// Add other fields as needed
}

type UserInfo struct {
	ID        int
	Username  string
	Password  string
	Level     int
	Active    bool
	Created   time.Time
	LastLogin *time.Time
}

// Maps
var ExUsermap = make(map[string]*UserMap)

// Global Variables

var (
	BBS_VERSION = "0.0.1"
)
