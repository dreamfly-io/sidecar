package server

import (
	"sync"
	"github.com/dreamfly-io/sidecar/pkg/util/log"
	"runtime"
	"github.com/dreamfly-io/sidecar/pkg/server/config"
	"github.com/dreamfly-io/sidecar/pkg/proxy/network/listener"
)

type Server interface {
	Start()

	Restart()

	Close()
}

type server struct {
	logger          log.Logger
	listenerConfigs sync.Map
	handler         listener.Handler
}

func NewServer(serverConfig *config.ServerConfig) Server {

	processorNumber := runtime.NumCPU()
	if serverConfig.Processor > 0 {
		processorNumber = serverConfig.Processor
	}
	runtime.GOMAXPROCS(processorNumber)

	s := &server{
		logger:          log.StartLogger,		//TODO: use StartLogger for debug, log.DefaultLogger for production
		listenerConfigs: sync.Map{},
		handler: listener.NewHandler(log.StartLogger),
	}
	for _, listenerConfig := range serverConfig.Listeners {
		s.addListener(listenerConfig)
	}

	return s
}

func (srv *server) Start() {
	srv.handler.StartListeners()
}

func (srv *server) Restart() {

}

func (srv *server) Close() {

}

func (srv *server) addListener(listenerConfig *config.ListenerConfig) {
	if _, ok := srv.listenerConfigs.Load(listenerConfig.Name); ok {
		log.DefaultLogger.Warn("Listen Already Started, Listen = %+v", listenerConfig)
		return
	}

	srv.listenerConfigs.Store(listenerConfig.Name, listenerConfig)
	srv.handler.AddListener(listenerConfig)
}
