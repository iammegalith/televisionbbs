package users

import (
	"bytes"
	"database/sql"
	"fmt"
	"thetelevision/system"
	"time"

	"github.com/ebarkie/telnet"
)

// ShowOnlineUsers displays the list of online users.
func ShowOnlineUsers(tn *telnet.Ctx) {
	tn.Write([]byte("\r\nOnline Users:\r\n"))
	for username, user := range system.ExUsermap {
		if user.Status == "Online" {
			tn.Write([]byte(fmt.Sprintf("\r\n  - %s\r\n", username)))
		}
	}
}

// SendToUser sends a message to the given recipient.
func SendToUser(tn *telnet.Ctx, recipient, message string) {
	if recipient == "" {
		tn.Write([]byte("\r\nInvalid recipient.\r\n"))
		return
	}

	if user, found := system.ExUsermap[recipient]; found {
		var sender *system.UserMap
		for _, u := range system.ExUsermap {
			if u.Ctx == tn {
				sender = u
				break
			}
		}
		if sender != nil {
			conn := user.Ctx
			conn.Write([]byte(fmt.Sprintf("\r\nMessage from %s: %s\r\n", sender.Username, message)))
		} else {
			tn.Write([]byte("\r\nError: Sender not found in users map.\r\n"))
		}
	} else {
		tn.Write([]byte(fmt.Sprintf("\r\nUser %s not found.\r\n", recipient)))
	}
}

func ReadString(ctx *telnet.Ctx, delim byte) (string, error) {
	var buf bytes.Buffer
	b := make([]byte, 1)
	for {
		_, err := ctx.Read(b)
		if err != nil {
			return "", err
		}
		if b[0] == delim {
			break
		}
		buf.Write(b)
	}
	return buf.String(), nil
}

func PrintToUser(ctx *telnet.Ctx, msg string) {
	ctx.Write([]byte(msg))
}

func CheckForUsername(db *sql.DB, username string) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetUser(db *sql.DB, username string) (*system.UserInfo, error) {
	var u system.UserInfo
	row := db.QueryRow(`SELECT id, username, password, level, active, created, last_login
                         FROM users
                         WHERE username = ?`, username)
	err := row.Scan(&u.ID, &u.Username, &u.Password, &u.Level, &u.Active, &u.Created, &u.LastLogin)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a non-nil zero-value of UserInfo if the user is not found
			return &system.UserInfo{}, nil
		}
		return nil, err
	}
	return &u, nil
}

func UpdateLastLogin(db *sql.DB, userID int) error {
	// Execute SQL update statement to set the last login time
	stmt, err := db.Prepare(`UPDATE users SET last_login = ? WHERE id = ?`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(time.Now(), userID)
	if err != nil {
		return err
	}
	return nil
}

func AddUserToMap(username string, tn *telnet.Ctx) {
	system.ExUsermap[username] = &system.UserMap{
		Username: username,
		Ctx:      tn,
		Status:   "Online",
	}
}
