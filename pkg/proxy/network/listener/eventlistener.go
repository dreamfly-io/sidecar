package listener

import (
	"net"
	"context"
	"github.com/dreamfly-io/sidecar/pkg/proxy/network/connection"
	"github.com/dreamfly-io/sidecar/pkg/util/log"
	ctx "github.com/dreamfly-io/sidecar/pkg/server/context"
	"strings"
	"strconv"
)

type EventListener interface {
	OnAccept(rawConnection net.Conn)

	OnNewConnection(context context.Context, conn connection.Connection)

	OnClose()
}

type eventListener struct {
	listener   Listener
	listenIP   string
	listenPort int
	stopChan   chan struct{}
	logger     log.Logger
}

func NewEventListener(listener Listener, stopChan chan struct{}, logger log.Logger) EventListener {
	el := &eventListener{
		listener: listener,
		stopChan: stopChan,
		logger:   logger,
	}

	el.listenIP, el.listenPort = parseIPAndPort(listener)

	return el
}

func parseIPAndPort(listener Listener) (listenIP string, listenPort int) {
	listenPort = 0
	listenIP = listener.Address().String()

	if temps := strings.Split(listenIP, ":"); len(temps) > 0 {
		listenPort, _ = strconv.Atoi(temps[len(temps)-1])
		listenIP = temps[0]
	}
	return
}

func (el *eventListener) OnAccept(rawConnection net.Conn) {
	log.DefaultLogger.Debug("Accepted connection from %s", el.listener.Address())

	m := connection.NewManager(rawConnection, el, el.stopChan, el.logger)
	c := prepareContext(el)
	m.ContinueFilterChain(c)
}

func prepareContext(el *eventListener) context.Context {
	c := context.WithValue(context.Background(), ctx.ContextKeyListenerName, el.listener.Name())
	c = context.WithValue(c, ctx.ContextKeyListenerIp, el.listenIP)
	c = context.WithValue(c, ctx.ContextKeyListenerPort, el.listenPort)
	c = context.WithValue(c, ctx.ContextKeyLogger, el.logger)
	return c
}

func (el *eventListener) OnNewConnection(context context.Context, conn connection.Connection) {

}

func (el *eventListener) OnClose() {

}
