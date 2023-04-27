package system

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ebarkie/telnet"
)

func ShowTextFile(conn *telnet.Ctx, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		conn.Write([]byte(fmt.Sprintf("\r\nError: %v\r\n", err)))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		conn.Write([]byte(scanner.Text() + "\r\n"))
	}

	if err := scanner.Err(); err != nil {
		conn.Write([]byte(fmt.Sprintf("\r\nError: %v\r\n", err)))
		return
	}
}