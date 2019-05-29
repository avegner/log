package log

import (
	"fmt"
	"sync"

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

const (
	ALL_LEVELS = ERROR | INFO | WARNING | DEBUG | CRITICAL
)

func (l Level) String() string {
	var s string

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
		s = "!unk!"
	}

	return s
}

func New(name string, mask Level, outs []out.Outputter) Logger {
	return &logger{
		name: name,
		common: &common{
			mask: mask,
			outs: outs,
		},
	}
}

func (l *logger) Printf(level Level, format string, args ...interface{}) error {
	if level&l.common.getMask() == 0x0 {
		return nil
	}

	if l.parent != nil {
		return l.parent.Printf(level, "["+l.name+"] "+format, args...)
	}

	rec := level.String() + " [" + l.name + "] " + fmt.Sprintf(format, args...) + "\n"

	for _, o := range l.common.getOuts() {
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
		common: l.common,
	}
}

type logger struct {
	parent Logger
	name   string
	*common
}

type common struct {
	mu   sync.RWMutex
	mask Level
	outs []out.Outputter
}

func (c *common) setMask(mask Level) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.mask = mask
}

func (c *common) getMask() Level {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.mask
}

func (c *common) getOuts() []out.Outputter {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.outs
}
