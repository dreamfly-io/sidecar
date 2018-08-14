package log

import "sync"

var (
	StartLogger          Logger

	remoteSyslogPrefixes = map[string]string{
		"syslog+tcp://": "tcp",
		"syslog+udp://": "udp",
		"syslog://":     "udp",
	}

	loggers []*logger
)

func init() {
	//use console  as start logger
	l := &logger{
		Output:  "",
		Level:   DEBUG,
		Roller:  DefaultRoller(),
		fileMux: new(sync.RWMutex),
	}
	l.Start()

	StartLogger= l
}