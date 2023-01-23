package messages

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	"televisionbbs/util"
)

func postMessage(db *sql.DB, basename, subject, author, postto, message string) error {
	stmt, err := db.Prepare("INSERT INTO messages(basename, subject, author, date, message, postto) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(basename, subject, author, time.Now(), message, postto)
	if err != nil {
		return fmt.Errorf("error inserting message: %v", err)
	}
	fmt.Println("Message posted successfully.")
	return nil
}

func messageEditor(conn net.Conn, db *sql.DB) error {
	reader := bufio.NewReader(conn)
	conn.Write([]byte("Enter message subject: "))
	subject, _ := reader.ReadString('\n')
	subject = strings.TrimSpace(subject)

	conn.Write([]byte("Enter message contents: "))
	contents, _ := reader.ReadString('\n')
	contents = strings.TrimSpace(contents)

	if err := postMessage(db, "general", subject, "John Doe", "general", contents); err != nil {
		conn.Write([]byte("Error posting message: " + err.Error()))
		return err
	}
	conn.Write([]byte("Message posted successfully."))
	return nil
}

func scanForNewMessages(db *sql.DB, lastRead time.Time) ([]Message, error) {
	rows, err := db.Query("SELECT id, basename, subject, author, date, message, postto FROM messages WHERE date > ?", lastRead)
	if err != nil {
		return nil, fmt.Errorf("Error querying messages: %v", err)
	}
	defer rows.Close()

	var messages util.[]Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.Basename, &m.Subject, &m.Author, &m.Date, &m.Message, &m.Postto); err != nil {
			return nil, fmt.Errorf("Error scanning message: %v", err)
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error fetching messages: %v", err)
	}
	return messages, nil
}