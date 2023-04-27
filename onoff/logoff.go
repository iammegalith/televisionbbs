package onoff

import (
	"net"

	"github.com/ebarkie/telnet"
)

func Logoff(ctx *telnet.Ctx, conn net.Conn) {
	ctx.Write([]byte("\r\nGoodbye!\r\n"))
	conn.Close()
}
