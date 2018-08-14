package listener

import "net"

type ListenerConfig struct {
	Name          string
	LocalAddrress net.Addr
	LogPath       string
	LogLevel      uint8
}
