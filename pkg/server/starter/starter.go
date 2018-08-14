package starter

import (
	"sync"
	"github.com/dreamfly-io/sidecar/pkg/util/log"
	"github.com/dreamfly-io/sidecar/pkg/server/config"
	"github.com/dreamfly-io/sidecar/pkg/server/server"
)

func Start(config *config.ServerConfig) {
	log.StartLogger.Info("start with server config : %+v", config)

	wg := sync.WaitGroup{}
	wg.Add(1)

	s := newServer(config)
	go s.Start()

	wg.Wait()
}

func newServer(c *config.ServerConfig) server.Server {
	log.InitDefaultLogger(c.DefaultLogPath, log.DEBUG)

	return server.NewServer(c)
}