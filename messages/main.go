package messages

import (
	"bufio"
	"database/sql"
	"fmt"
	"net"
	"os/exec"
	"strconv"
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

func scanForNewMessages(db *sql.DB, lastRead time.Time) ([]util.Message, error) {
	rows, err := db.Query("SELECT id, basename, subject, author, date, message, postto FROM messages WHERE date > ?", lastRead)
	if err != nil {
		return nil, fmt.Errorf("error querying messages: %v", err)
	}
	defer rows.Close()

	var messages []util.Message
	for rows.Next() {
		var m util.Message
		if err := rows.Scan(&m.ID, &m.Basename, &m.Subject, &m.Author, &m.Date, &m.Message, &m.Postto); err != nil {
			return nil, fmt.Errorf("error scanning message: %v", err)
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error fetching messages: %v", err)
	}
	return messages, nil
}
func fullScreenEditor(conn net.Conn, db *sql.DB, header, messageBase string) error {
	var input string
	var subject, author, postto, message string
	conn.Write([]byte(header + "\n"))
	conn.Write([]byte("Message Base: " + messageBase + "\n"))
	conn.Write([]byte("Ctl-Q to Quit | Ctl-S to save and post\n"))
	// get terminal size
	cmd := exec.Command("stty", "size")
	cmd.Stdin = conn
	out, _ := cmd.Output()
	size := strings.Split(string(out), " ")
	rows, _ := strconv.Atoi(size[0])

	reader := bufio.NewReader(conn)
	for {
		line, _, _ := reader.ReadLine()
		input += string(line) + "\n"
		conn.Write([]byte("\033[1A")) // move cursor up one line
		conn.Write([]byte("\033[K"))  // clear current line
		if len(input) > 0 && input[0] == '.' {
			switch input {
			case ".s":
				// Save and post message
				if err := postMessage(db, messageBase, subject, author, postto, message); err != nil {
					return fmt.Errorf("error posting message: %v", err)
				}
				return nil
			case ".q":
				// Quit editor without saving
				conn.Write([]byte("Exiting editor.\n"))
				return nil
			default:
				message += input
			}
			input = ""
		}
		if len(strings.Split(input, "\n")) >= rows-2 { // subtract 2 to account for the header and bottom bar
			conn.Write([]byte("\033[2J")) // clear screen
			conn.Write([]byte("\033[H"))  // move cursor to top left
			conn.Write([]byte(header + "\n"))
			conn.Write([]byte("Message Base: " + messageBase + "\n"))
			conn.Write([]byte("Ctl-Q to Quit | Ctl-S to save and post\n"))
			conn.Write([]byte(input))
		}
	}
}
