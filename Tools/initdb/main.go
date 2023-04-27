package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s database_name schema.sql\n", os.Args[0])
		os.Exit(1)
	}
	dbName := os.Args[1]
	schemaFile := os.Args[2]

	// Read the schema from the file
	schemaBytes, err := os.ReadFile(schemaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading schema file: %s\n", err)
		os.Exit(1)
	}

	// Connect to the database
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()

	// Execute the schema SQL to create the tables
	_, err = db.Exec(string(schemaBytes))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing schema SQL: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Database initialized successfully!\n")
}
