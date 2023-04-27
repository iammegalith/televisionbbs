package bbsmenu

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"thetelevision/channel"
	"thetelevision/commands"
	"thetelevision/doormanager"
	"thetelevision/system"

	"github.com/ebarkie/telnet"
	"github.com/go-ini/ini"
)

func GetMenu(ctx *telnet.Ctx, conn net.Conn, cs *channel.ChannelServer, user system.UserInfo, menuName string) {
	currentMenu := menuName
	var err error
	var mcfg *ini.File
	mcfg, err = ini.Load(fmt.Sprintf(system.BBSConfig.MenusPath+"%s.config", menuName))
	log.Printf("Path: %s", system.BBSConfig.MenusPath+menuName)
	if err != nil {
		ctx.Write([]byte(fmt.Sprintf("\r\nError: %v\r\n", err)))
		return
	}

	options := mcfg.Sections()

	system.ShowTextFile(ctx, fmt.Sprintf(system.BBSConfig.AsciiPath+"%s.asc", menuName))
	ctx.Write([]byte("\r\nEnter your selection: "))

	buf := make([]byte, 1024)
	n, err := ctx.Read(buf)
	if err != nil && err != io.EOF {
		ctx.Write([]byte(fmt.Sprintf("\r\nError reading input: %v\r\n", err)))
		return
	}
	user_input := strings.TrimSpace(string(buf[:n]))

	var selection *ini.Section
	for _, option := range options {
		if option.Key("fast").String() == user_input {
			selection = option
			break
		}
	}

	if selection == nil {
		ctx.Write([]byte("\r\nInvalid selection\r\n"))
		GetMenu(ctx, conn, cs, user, currentMenu)
		return
	}

	typ := selection.Key("type").String()
	args := selection.Key("arguments").String()
	lvl := selection.Key("level").MustInt(0)

	if typ == "menu" {
		GetMenu(ctx, conn, cs, user, args)
	} else {
		handleSelection(ctx, conn, cs, user, selection, lvl, currentMenu)
	}
}

func handleSelection(ctx *telnet.Ctx, conn net.Conn, cs *channel.ChannelServer, user system.UserInfo, selection *ini.Section, mylvl int, currentMenu string) {
	typ := selection.Key("type").String()
	args := selection.Key("arguments").String()
	lvl := selection.Key("level").MustInt(0)
	ctx.Write([]byte(fmt.Sprintf("\r\nLevel: %d\r\n", user.Level)))
	ctx.Write([]byte(fmt.Sprintf("Command: %s\r\n", args)))
	switch typ {
    case "menu":
        GetMenu(ctx, conn, cs, user, args)
	case "command":
		if user.Level >= lvl {
			// Execute the command here
			conn.Write([]byte(fmt.Sprintf("\r\nExecuting command: %s\r\n", args)))
			if args == "obbs" {
				commands.Obbs(ctx)
			}
			if args == "goodbye" {
				ctx.Write([]byte("\r\nGoodbye!\r\n"))
				conn.Close()
			}
			GetMenu(ctx, conn, cs, user, currentMenu) // Updated
		} else {
			conn.Write([]byte("\r\nInsufficient access level\r\n"))
		}

	case "door":
		if user.Level >= lvl {
			// Run the external application here
			doorName := args // assuming args contains the door name
			doormanager.RunDoor(conn, user, cs, doorName, currentMenu)
		} else {
			conn.Write([]byte("\r\nInsufficient access level\r\n"))
		}
	case "text":
		system.ShowTextFile(ctx, fmt.Sprintf("ascii/%s.asc", args))
	default:
		conn.Write([]byte("\r\nInvalid selection\r\n"))
		GetMenu(ctx, conn, cs, user, currentMenu)
	}
}