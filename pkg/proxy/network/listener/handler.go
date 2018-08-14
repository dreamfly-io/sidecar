package listener

import (
	"github.com/dreamfly-io/sidecar/pkg/server/config"
	"github.com/dreamfly-io/sidecar/pkg/util/log"
)

type Handler interface {

	AddListener(listenerConfig *config.ListenerConfig)

	StartListeners()
}

type handler struct {
	logger         log.Logger
	listeners 	[]Listener
}

func NewHandler(logger log.Logger) Handler {
	ch := &handler{
		logger:         logger,
		listeners: make([]Listener, 0),
	}

	return ch
}

func (h *handler) AddListener(lc *config.ListenerConfig) {
	l := NewListener(lc, h.logger)
	h.listeners = append(h.listeners, l)
}

func (h *handler) StartListeners() {
	for _, l := range h.listeners {
		go l.Start(nil)
	}
}