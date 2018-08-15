package listener

import (
	"net"
	"context"
	"github.com/dreamfly-io/sidecar/pkg/util/log"
	"runtime/debug"
	"time"
	"github.com/dreamfly-io/sidecar/pkg/server/config"
)

// Listener is a wrapper of tcp listener
type Listener interface {
	Name() string

	Address() net.Addr

	SetEventListener(eventListener EventListener)

	GetEventListener() EventListener

	Start(context context.Context)

	Stop()
}

// implement listener.Listener
type listener struct {
	name          string
	address       net.Addr
	addressString string
	rawListener   *net.TCPListener
	eventListener EventListener
	stopChan      chan struct{}
	logger        log.Logger
}

func NewListener(config *config.ListenerConfig, stopChan chan struct{}, logger log.Logger) Listener {
	l := &listener{
		name:          config.Name,
		addressString: config.Address,
		stopChan:      stopChan,
		logger:        logger,
	}

	return l
}

func (l *listener) Start(context context.Context) {
	if l.rawListener == nil {
		if err := l.listen(context); err != nil {
			log.StartLogger.Fatal("Fail to start listener: address=%s, error=%-v", l.addressString, err)
			return
		}
		log.StartLogger.Info("Success to start listener: address=%s", l.addressString)
	}

	for {
		if err := l.accept(context); err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				l.logger.Info("listener %s stop accepting connections by deadline", l.name)
				return
			} else if ope, ok := err.(*net.OpError); ok {
				// not timeout error and not temporary, which means the error is non-recoverable
				// stop accepting loop and log the event
				if !(ope.Timeout() && ope.Temporary()) {
					// accept error raised by sockets closing
					if ope.Op == "accept" {
						l.logger.Info("listener %s %s closed", l.name, l.address)
					} else {
						l.logger.Error("listener %s occurs non-recoverable error, stop listening and accepting:%s", l.name, err.Error())
					}
					return
				}
			} else {
				l.logger.Error("listener %s occurs unknown error while accepting:%s", l.name, err.Error())
			}
		}
	}
}

func (l *listener) listen(context context.Context) error {
	addr, err := net.ResolveTCPAddr("tcp", l.addressString)
	if err != nil {
		return err
	}
	l.address = addr

	rawListener, err := net.ListenTCP("tcp", l.address.(*net.TCPAddr))
	if err != nil {
		return err
	}

	l.rawListener = rawListener
	return nil
}

func (l *listener) accept(lctx context.Context) error {
	rawConnection, err := l.rawListener.Accept()

	if err != nil {
		return err
	}
	l.logger.Trace("Accepted a connection on listener %s", l.addressString)

	go func() {
		defer func() {
			if p := recover(); p != nil {
				l.logger.Error("panic %v", p)

				debug.PrintStack()
			}
		}()

		l.eventListener.OnAccept(rawConnection)
	}()

	return nil
}

func (l *listener) SetEventListener(eventListener EventListener) {
	l.eventListener = eventListener
}

func (l *listener) GetEventListener() EventListener {
	return l.eventListener
}

func (l *listener) Stop() {
	l.rawListener.SetDeadline(time.Now())
}

func (l *listener) Name() string {
	return l.name
}

func (l *listener) Address() net.Addr {
	return l.address
}
