package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type UserInfo struct {
	ID        int
	Username  string
	Password  string
	Level     int
	Active    bool
	Created   time.Time
	LastLogin *time.Time
}

func main() {
	var username, level, password, dataDir string

	// parse command line arguments
	flag.StringVar(&username, "u", "", "Username")
	flag.StringVar(&level, "l", "", "Level")
	flag.StringVar(&password, "p", "", "Password")
	flag.StringVar(&dataDir, "d", "", "Path to data directory")
	flag.Parse()

	// Check required arguments
	if username == "" || level == "" {
		fmt.Println("Usage: createuser -d /path/to/data/directory -u username -l level [-p password]")
		os.Exit(1)
	}

	if password == "" {
		password = "pass"
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// Open the database
	dbPath := filepath.Join(dataDir, "userdata.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert the new user into the database
	stmt, err := db.Prepare(`INSERT INTO users (username, password, level, active, created, last_login)
                              VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, hashedPassword, level, true, time.Now(), nil)
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created user %s with level %s. Rows affected: %d\n", username, level, rowsAffected)
}
