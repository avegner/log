package log

import (
	"fmt"

	"github.com/avegner/log/out"
)

type Logger interface {
	Printf(level Level, format string, args ...interface{}) error
	Child(prefix string) Logger
}

type Level int

const (
	ERROR    = Level(0x1)
	INFO     = Level(0x2)
	WARNING  = Level(0x4)
	DEBUG    = Level(0x8)
	CRITICAL = Level(0x10)
)

func (l Level) String() string {
	s := ""

	switch l {
	case ERROR:
		s = "ERR"
	case INFO:
		s = "INF"
	case WARNING:
		s = "WRN"
	case DEBUG:
		s = "DBG"
	case CRITICAL:
		s = "CRI"
	default:
		panic("unknown log level")
	}

	return s
}

func New(name string, outs []out.Outputter) Logger {
	return &logger{
		name: name,
		outs: outs,
	}
}

type logger struct {
	parent Logger
	name   string
	outs   []out.Outputter
}

func (l *logger) Printf(level Level, format string, args ...interface{}) error {
	if l.parent != nil {
		return l.parent.Printf(level, "["+l.name+"] "+format, args...)
	}

	rec := level.String() + " [" + l.name + "] " + fmt.Sprintf(format, args...) + "\n"

	for _, o := range l.outs {
		if _, err := o.Write([]byte(rec)); err != nil {
			return err
		}
		o.Flush()
	}

	return nil
}

func (l *logger) Child(name string) Logger {
	return &logger{
		parent: l,
		name:   name,
	}
}
