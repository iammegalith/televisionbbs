package onoff

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"strings"
	"thetelevision/system"
	"thetelevision/users"

	"github.com/ebarkie/telnet"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB, tn *telnet.Ctx, conn net.Conn) (system.UserInfo, error) {
	// Prompt user for username and password
	var u system.UserInfo
	var password string
	for i := 0; i < 3; i++ {
		var username string
		users.PrintToUser(tn, "(NEW for NEW USER)\r\n  -> Enter username: ")
		username, err := users.ReadString(tn, '\r')
		if err != nil {
			return u, err
		}
		username = strings.TrimSpace(username)
		if username == "new" {
			err = RegisterUser(db, tn)
			if err != nil {
				return u, err
			}
			// Exit the Login function and return without an error
			return u, nil
		}
		// Check if the username exists
		exists, err := users.CheckForUsername(db, username)
		if err != nil {
			return u, err
		}

		if !exists {
			users.PrintToUser(tn, "Username does not exist. Please try again.\r\n")
			continue
		}
		fmt.Println("username checked.." + username)
		users.PrintToUser(tn, "\r\nEnter password: ")
		password, err = users.ReadString(tn, '\r')
		if err != nil {
			return u, err
		}
		password = strings.TrimSpace(password)
		// Get the user from the database
		fmt.Println("Password read.." + password)
		var lastLogin sql.NullTime
		row := db.QueryRow(`SELECT id, username, password, level, active, created, last_login
							 FROM users
							 WHERE username = ?`, username)
		err = row.Scan(&u.ID, &u.Username, &u.Password, &u.Level, &u.Active, &u.Created, &lastLogin)
		if err != nil {
			return u, err
		}
		if lastLogin.Valid {
			lastLoginTime := lastLogin.Time
			u.LastLogin = &lastLoginTime
		}
		fmt.Println("User checked: " + u.Username)
		// Check the password
		err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		if err != nil {
			users.PrintToUser(tn, "Invalid password. Please try again.\r\n")
			continue
		}
		fmt.Println("updating lastlogin...")
		// Update the last login time for the user in the database
		err = users.UpdateLastLogin(db, u.ID)
		if err != nil {
			return u, err
		}
		fmt.Println("Login successful.")
		// Login successful, return the user
		fmt.Println("Added user to usermap: ", username)
		users.AddUserToMap(username, tn)
		thisuser := system.ExUsermap[username]
		fmt.Println("User has status :", username, thisuser.Status)

		return u, nil
	}

	// If the user has tried three times and still failed, log them out
	Logoff(tn, conn)

	return u, errors.New("login failed")
}
