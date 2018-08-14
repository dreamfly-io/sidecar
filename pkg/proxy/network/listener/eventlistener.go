package listener

import (
	"net"
	"context"
	"github.com/dreamfly-io/sidecar/pkg/proxy/network/connection"
)

type EventListener interface {

	OnAccept(rawConnection net.Conn)

	OnNewConnection(context context.Context, conn connection.Connection)

	OnClose()
}
