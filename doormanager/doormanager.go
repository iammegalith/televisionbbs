// doormanager.go
package doormanager

import (
	"errors"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"strings"
	"thetelevision/system"

	"github.com/go-ini/ini"
	"thetelevision/channel" // Replace with the import path of your BBS application
)

type Door interface {
	Name() string
	Description() string
	Play(conn net.Conn, cs *channel.ChannelServer, user system.UserInfo)
}

var (
	ErrNoSymbol       = errors.New("no door symbol specified")
	ErrInvalidSymbol  = errors.New("invalid door symbol format")
	ErrSymbolNotDoor  = errors.New("symbol does not implement Door interface")
	ErrInvalidDoorCfg = errors.New("invalid door configuration")
)

func LoadDoor(name string) (Door, error) {
	doorConfigFile := filepath.Join("modules", name, name+".config")
	doorConfig, err := ini.Load(doorConfigFile)
	if err != nil {
		return nil, err
	}

	doorPath := doorConfig.Section("Door").Key("path").String()
	doorArgsStr := doorConfig.Section("Door").Key("args").String()
	doorArgs := strings.Fields(doorArgsStr)

	doorSymbol := doorConfig.Section("Door").Key("symbol").String()

	cmd := exec.Command(doorPath, doorArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Printf("Opening Door: %s", doorPath)
	doorPlugin, err := plugin.Open(doorPath)
	if err != nil {
		return nil, err
	}

	sym, err := doorPlugin.Lookup(doorSymbol)
	if err != nil {
		return nil, err
	}

	door, ok := sym.(Door)
	if !ok {
		return nil, ErrSymbolNotDoor
	}

	return door, nil
}

func RunDoor(conn net.Conn, user system.UserInfo, cs *channel.ChannelServer, doorName string, lastMenu string) {
	conn.Write([]byte("\r\nRunning door...\r\n"))

	door, err := LoadDoor(doorName)
	if err != nil {
		log.Printf("Error loading door %s: %v", doorName, err)
		conn.Write([]byte("\r\nError loading door\r\n"))
		return
	}

	door.Play(conn, cs, user)

	// Return the user to the last menu they were on
	conn.Write([]byte("\r\nPress any key to return to the menu...\r\n"))
	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != nil {
		log.Printf("Error reading input after door: %v", err)
	}
}

