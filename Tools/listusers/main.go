package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Parse command-line arguments
	dbPath := flag.String("d", "userdata.db", "path to the userdata database")
	flag.Parse()
	fmt.Println(*dbPath)
	// Open the database
	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Query the users table
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Print the column names
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(columns)

	// Print the rows
	for rows.Next() {
		var id int
		var username string
		var password string
		var level int
		var active int
		var created string
		var last_login sql.NullString

		err = rows.Scan(&id, &username, &password, &level, &active, &created, &last_login)
		if err != nil {
			log.Fatal(err)
		}
		if last_login.Valid {
			fmt.Printf("%d %s %s %d %d %s %s\n", id, username, password, level, active, created, last_login.String)
		} else {
			fmt.Printf("%d %s %s %d %d %s NULL\n", id, username, password, level, active, created)
		}
	}
}
