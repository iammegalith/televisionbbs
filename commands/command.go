package commands

import (
	"thetelevision/system"

	"github.com/ebarkie/telnet"
)

func Obbs(conn *telnet.Ctx) {
	system.ShowTextFile(conn, "ascii/obbs.asc")
}