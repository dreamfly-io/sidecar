package log

import (
	"io"
	"sync"
	"log"
	"os"
	"path/filepath"
	"fmt"
	"github.com/hashicorp/go-syslog"
	"strings"
)

type Level uint8

const (
	FATAL Level = iota
	ERROR
	WARN
	INFO
	DEBUG
	TRACE
)

const (
	InfoPre  string = "[INFO]"
	DebugPre string = "[DEBUG]"
	WarnPre  string = "[WARN]"
	ErrorPre string = "[ERROR]"
	FatalPre string = "[Fatal]"
	TracePre string = "[TRACE]"
)

type Logger interface {
	Trace(format string, args ...interface{})

	Info(format string, args ...interface{})

	Debug(format string, args ...interface{})

	Warn(format string, args ...interface{})

	Error(format string, args ...interface{})

	Fatal(format string, args ...interface{})
}

type logger struct {
	*log.Logger

	Output  string
	Level   Level
	Roller  *Roller
	writer  io.Writer
	fileMux *sync.RWMutex
}

func (l *logger) Info(format string, args ...interface{}) {
	if l.Level >= INFO {
		l.Printf(InfoPre+format, args...)
	}
}

func (l *logger) Debug(format string, args ...interface{}) {
	if l.Level >= DEBUG {
		l.Printf(DebugPre+format, args...)
	}
}

func (l *logger) Warn(format string, args ...interface{}) {
	if l.Level >= WARN {
		l.Printf(WarnPre+format, args...)
	}
}

func (l *logger) Error(format string, args ...interface{}) {
	if l.Level >= ERROR {
		l.Printf(ErrorPre+format, args...)
	}
}

func (l *logger) Trace(format string, args ...interface{}) {
	if l.Level >= TRACE {
		l.Printf(TracePre+format, args...)
	}
}

func (l *logger) Fatal(format string, args ...interface{}) {
	if l.Level >= FATAL {
		l.Printf(FatalPre+format, args...)
	}
}

type syslogAddress struct {
	network string
	address string
}

func (l *logger) Start() error {
	var err error

selectwriter:
	switch l.Output {
	case "", "stderr", "/dev/stderr":
		l.writer = os.Stderr
	case "stdout", "/dev/stdout":
		l.writer = os.Stdout
	case "syslog":
		l.writer, err = gsyslog.NewLogger(gsyslog.LOG_ERR, "LOCAL0", "mosn")
		if err != nil {
			return err
		}
	default:
		if address := parseSyslogAddress(l.Output); address != nil {
			l.writer, err = gsyslog.DialLogger(address.network, address.address, gsyslog.LOG_ERR, "LOCAL0", "mosn")

			if err != nil {
				return err
			}

			break selectwriter
		}

		var file *os.File

		//create parent dir if not exists
		err := os.MkdirAll(filepath.Dir(l.Output), 0755)

		fmt.Println(err)

		file, err = os.OpenFile(l.Output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}

		if l.Roller != nil {
			file.Close()
			l.Roller.Filename = l.Output
			l.writer = l.Roller.GetLogWriter()
		} else {
			l.writer = file
		}
	}

	l.Logger = log.New(l.writer, "", log.LstdFlags)

	return nil
}

func parseSyslogAddress(location string) *syslogAddress {
	for prefix, network := range remoteSyslogPrefixes {
		if strings.HasPrefix(location, prefix) {
			return &syslogAddress{
				network: network,
				address: strings.TrimPrefix(location, prefix),
			}
		}
	}

	return nil
}
