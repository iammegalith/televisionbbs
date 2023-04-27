package onoff

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"thetelevision/system"
	"thetelevision/users"
	"time"

	"github.com/ebarkie/telnet"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB, tn *telnet.Ctx) error {
	// Prompt user for username and password
	var username string
	var password string
	for i := 0; i < 3; i++ {
		users.PrintToUser(tn, "Enter username: ")
		var err error
		username, err = users.ReadString(tn, '\r')
		if err != nil {
			return err
		}
		username = strings.TrimSpace(username)
		// Check if the username already exists
		exists, err := users.CheckForUsername(db, username)
		if err != nil {
			return err
		}

		if !exists {
			break
		}

		users.PrintToUser(tn, "Username already exists. Please try again.\r\n")
	}

	for i := 0; i < 3; i++ {
		users.PrintToUser(tn, "Enter password: ")
		var err error
		password, err = users.ReadString(tn, '\r')
		password = strings.TrimSpace(password)
		if err != nil {
			fmt.Println("Error reading password")
		}

		users.PrintToUser(tn, "Confirm password: ")
		confirmPassword, err := users.ReadString(tn, '\r')
		confirmPassword = strings.TrimSpace(confirmPassword)
		if err != nil {
			return err
		}

		if password == confirmPassword {
			break
		}

		users.PrintToUser(tn, "Passwords do not match. Please try again.\r\n")
	}

	if password == "" {
		return errors.New("password cannot be blank")
	}

	// Hash the password using a secure algorithm like bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Insert the new user into the database
	stmt, err := db.Prepare(`INSERT INTO users (username, password, level, active, created, last_login)
                              VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}

	u := system.UserInfo{
		Username:  username,
		Password:  string(hashedPassword),
		Level:     1,
		Active:    true,
		Created:   time.Now(),
		LastLogin: new(time.Time),
	}

	_, err = stmt.Exec(u.Username, u.Password, u.Level, u.Active, u.Created, u.LastLogin)
	if err != nil {
		return err
	}
	users.AddUserToMap(username, tn)
	return nil
}