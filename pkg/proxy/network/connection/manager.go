package connection

import (
	"net"
	"github.com/dreamfly-io/sidecar/pkg/proxy/network/listener"
	"context"
	"github.com/dreamfly-io/sidecar/pkg/util/log"
)

type Manager struct {
	rawConnection         net.Conn
	listenerEventListener listener.EventListener
	acceptedFilters       []listener.Filter
	acceptedFilterIndex   int
	stopChan              chan struct{}
	logger                log.Logger
}

func NewManager(rawConnection net.Conn, listenerEventListener listener.EventListener, stopChan chan struct{},
	logger log.Logger) *Manager {
	return &Manager{
		rawConnection:         rawConnection,
		listenerEventListener: listenerEventListener,
		stopChan:              stopChan,
		logger:                logger,
	}
}

func (m *Manager) ContinueFilterChain(ctx context.Context) {
	// 1. call listener filters to check if we should continue to serve this connection
	for ; m.acceptedFilterIndex < len(m.acceptedFilters); m.acceptedFilterIndex++ {
		s := m.acceptedFilters[m.acceptedFilterIndex].OnAccept()
		if s == listener.StopIteration {
			return
		}
	}

	// 2. if true, create a Connection
	c := NewServerConnection(m.rawConnection, m.stopChan, m.logger)

	// 3. and fire the event of OnNewConnection
	m.listenerEventListener.OnNewConnection(ctx, c)
}
