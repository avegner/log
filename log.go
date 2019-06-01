package log

import (
	"fmt"
	"sync"
	"time"

	"github.com/avegner/log/out"
)

type Logger interface {
	Printf(level Level, format string, args ...interface{}) error
	SetMask(mask Level)
	Flush()
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

type Flag int

const (
	SHORT_TIME_PREFIX = Flag(0x1)
	LONG_TIME_PREFIX  = Flag(0x2)
)

const (
	STD_FLAGS = SHORT_TIME_PREFIX
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

func New(name string, mask Level, flags Flag, outs []out.Outputter) Logger {
	return &logger{
		name: name,
		common: &common{
			mask:  mask,
			flags: flags,
			outs:  outs,
		},
	}
}

func (l *logger) Printf(level Level, format string, args ...interface{}) error {
	if level&l.common.getMask() == 0x0 {
		return nil
	}

	t := time.Now()
	pr := ""

	for cl := l; cl != nil; cl = cl.parent {
		pr = "[" + cl.name + "]" + pr
	}
	pr = level.String() + " " + pr

	ts := ""
	if l.common.getFlags()&SHORT_TIME_PREFIX != 0x0 {
		ts = t.Format("15:04:05.000")
	} else if l.common.getFlags()&LONG_TIME_PREFIX != 0x0 {
		ts = t.Format("2006-01-02 15:04:05.000 -07")
	}

	rec := ts + " " + pr + " " + fmt.Sprintf(format, args...) + "\n"

	for _, o := range l.common.getOuts() {
		if _, err := o.Write([]byte(rec)); err != nil {
			return err
		}
	}

	return nil
}

func (l *logger) SetMask(mask Level) {
	l.common.setMask(mask)
}

func (l *logger) Flush() {
	for _, o := range l.common.getOuts() {
		o.Flush()
	}
}

func (l *logger) Child(name string) Logger {
	return &logger{
		parent: l,
		name:   name,
		common: l.common,
	}
}

type logger struct {
	parent *logger
	name   string
	*common
}

type common struct {
	mu    sync.RWMutex
	mask  Level
	flags Flag
	outs  []out.Outputter
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

func (c *common) getFlags() Flag {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.flags
}

func (c *common) getOuts() []out.Outputter {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.outs
}
