package log_writer

import (
	"fmt"
	"os"

	"github.com/flywithbug/log4go"
)

type colorRecord log4go.Record

var (
	LEVEL_FLAGS = [...]string{"DEBUG", "INFO", "INFO", "WARN", "ERROR", "FATAL"}
)

type Writer interface {
	Init() error
	Write(*log4go.Record) error
}

const (
	DEBUG = iota
	INFO
	NEVERShow
	WARNING
	ERROR
	FATAL
)

func (r *colorRecord) String() string {
	switch r.Level {
	case DEBUG:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[34m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.Time, LEVEL_FLAGS[r.Level], r.Code, r.Info)
	case INFO, NEVERShow:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[32m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.Time, LEVEL_FLAGS[r.Level], r.Code, r.Info)
	case WARNING:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[33m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.Time, LEVEL_FLAGS[r.Level], r.Code, r.Info)
	case ERROR:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[31m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.Time, LEVEL_FLAGS[r.Level], r.Code, r.Info)
	case FATAL:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[35m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.Time, LEVEL_FLAGS[r.Level], r.Code, r.Info)
	}
	return ""
}

type ConsoleWriter struct {
	color bool
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

func (w *ConsoleWriter) Write(r *log4go.Record) error {
	if r.Level == NEVERShow {
		return nil
	}
	if w.color {
		fmt.Fprint(os.Stdout, ((*colorRecord)(r)).String())
	} else {
		fmt.Fprint(os.Stdout, r.String())
	}
	return nil
}

func (w *ConsoleWriter) Init() error {
	return nil
}

func (w *ConsoleWriter) SetColor(c bool) {
	w.color = c
}
