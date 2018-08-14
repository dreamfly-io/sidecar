package config

func Load(path string) *ServerConfig {
	server := &ServerConfig{
		Processor: 1,
		DefaultLogPath:  "",
		DefaultLogLevel: "INFO",
	}

	listeners := make([]*ListenerConfig, 0, 1)
	listener := &ListenerConfig{
		Name: "HTTP",
		Address: "localhost:7001",
		LogPath:  "stdout",
		LogLevel: "INFO",
	}
	listeners = append(listeners, listener)
	server.Listeners = listeners

	return server
}


