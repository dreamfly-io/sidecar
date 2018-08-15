package connection

import (
	"context"
	"net"
	"github.com/dreamfly-io/sidecar/pkg/util/log"
	"sync/atomic"
	"runtime/debug"
	"sync"
	"time"
	"io"
	"runtime"
)

var idCounter uint64 = 1

// CloseType represent connection close type
type CloseType string

//Connection close types
const (
	// FlushWrite means write buffer to underlying io then close connection
	FlushWrite CloseType = "FlushWrite"
	// NoFlush means close connection without flushing buffer
	NoFlush CloseType = "NoFlush"
)

type Connection interface {
	// ID returns unique connection id
	ID() uint64

	// Start starts connection with context.
	Start(ctx context.Context)

	// Close closes connection with connection type and event type.
	Close(ct CloseType, e Event) error
}

type connection struct {
	id            uint64
	rawConnection net.Conn
	localAddress  net.Addr
	remoteAddress net.Addr
	logger        log.Logger

	status status
}

type status struct {
	internalLoopStarted bool
	startOnce           sync.Once
	stopChan            chan struct{}
	internalStopChan    chan struct{}
	readEnabled         bool
	readEnabledChan     chan bool
}

func NewServerConnection(rawConnection net.Conn, stopChan chan struct{}, logger log.Logger) Connection {
	id := atomic.AddUint64(&idCounter, 1)

	c := &connection{
		id:            id,
		rawConnection: rawConnection,
		localAddress:  rawConnection.LocalAddr(),
		remoteAddress: rawConnection.RemoteAddr(),
		logger:        logger,
		status: status{
			internalStopChan: make(chan struct{}),
			stopChan:         stopChan,
			readEnabled:      true,
			readEnabledChan:  make(chan bool, 1),
		},
	}

	return c
}

func (c *connection) ID() uint64 {
	return c.id
}

func (c *connection) Start(ctx context.Context) {
	c.status.startOnce.Do(func() {
		c.status.internalLoopStarted = true

		go func() {
			defer func() {
				if p := recover(); p != nil {
					c.logger.Error("panic %v", p)
					debug.PrintStack()

					c.startReadLoop()
				}
			}()

			c.startReadLoop()
		}()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					c.logger.Error("panic %v", p)
					debug.PrintStack()

					//c.startWriteLoop()
				}
			}()

			//c.startWriteLoop()
		}()
	})
}

func (c *connection) startReadLoop() {
	for {
		// exit loop asap. one receive & one default block will be optimized by go compiler
		select {
		case <-c.status.internalStopChan:
			return
		default:
		}

		select {
		case <-c.status.internalStopChan:
			return
		case <-c.status.readEnabledChan:
		default:
			if c.status.readEnabled {
				err := c.doRead()

				if err != nil {
					if err == io.EOF {
						c.Close(NoFlush, RemoteClose)
					} else {
						c.Close(NoFlush, OnReadErrClose)
					}

					c.logger.Error("Error on read. Connection = %d, Remote Address = %s, err = %s",
						c.id, c.remoteAddress.String(), err)
					return
				}
			} else {
				select {
				case <-c.status.readEnabledChan:
				case <-time.After(100 * time.Millisecond):
				}
			}

			runtime.Gosched()
		}
	}
}

func (c *connection) doRead() (err error) {
	return nil
}

func (c *connection) Close(ct CloseType, e Event) error {
	return nil
}
