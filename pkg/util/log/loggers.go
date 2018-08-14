package log

import (
	"sync"
	"os"
)

const (
	BasePath = string(os.PathSeparator) + "tmp" + string(os.PathSeparator) + "sidecar"
	LogBasePath = BasePath + string(os.PathSeparator) + "logs"
	LogDefaultPath = LogBasePath + string(os.PathSeparator) + "sidecar.log"
)

var (
	StartLogger          Logger
	DefaultLogger        Logger

	remoteSyslogPrefixes = map[string]string{
		"syslog+tcp://": "tcp",
		"syslog+udp://": "udp",
		"syslog://":     "udp",
	}

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

func InitDefaultLogger(path string, level Level) {

	var logPath string
	var logLevel Level

	//use default log path
	if path == "" {
		logPath = LogDefaultPath
	} else {
		logPath = path
	}
	logLevel = level


	l := &logger{
		Output:  logPath,
		Level:   logLevel,
		Roller:  DefaultRoller(),
		fileMux: new(sync.RWMutex),
	}

	error :=l.Start()
	if error != nil {
		StartLogger.Fatal("Fail to initialize default logger: %+v", error)
	}

	StartLogger.Debug("Success to initialize default logger: logPath=%+v, logLevel=%+v", logPath, logLevel)
	DefaultLogger = l
}