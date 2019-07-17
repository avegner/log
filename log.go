package log

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/avegner/log/out"
)

type Logger interface {
	Printf(level Level, format string, args ...interface{}) error
	SetMask(mask Level)
	Flush() error
	Child(name string) Logger
	Dump(level Level, data []byte) error
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
	SHORT_TIMESTAMP = Flag(0x1)
	LONG_TIMESTAMP  = Flag(0x2)
)

const (
	STD_FLAGS = SHORT_TIMESTAMP
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
	pr = level.String() + ":" + pr

	ts := ""
	if l.common.getFlags()&SHORT_TIMESTAMP != 0x0 {
		ts = t.Format("15:04:05.000")
	} else if l.common.getFlags()&LONG_TIMESTAMP != 0x0 {
		ts = t.Format("2006-01-02 15:04:05.000 -07")
	}

	rec := ts + " " + pr + " " + fmt.Sprintf(format, args...) + "\n"
	errss := []string{}

	for i, o := range l.common.getOuts() {
		if _, err := o.Write([]byte(rec)); err != nil {
			errss = append(errss, fmt.Sprintf("write out %d: '%v'", i, err))
		}
	}

	if len(errss) > 0 {
		return errors.New(strings.Join(errss, ", "))
	}
	return nil
}

func (l *logger) SetMask(mask Level) {
	l.common.setMask(mask)
}

func (l *logger) Flush() error {
	errss := []string{}

	for i, o := range l.common.getOuts() {
		if err := o.Flush(); err != nil {
			errss = append(errss, fmt.Sprintf("flush out %d: '%v'", i, err))
		}
	}

	if len(errss) > 0 {
		return errors.New(strings.Join(errss, ", "))
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

func (l *logger) Dump(level Level, data []byte) error {
	lss := make([]string, 0, 16)
	pl := func(ss []string) error {
		return l.Printf(level, "%s", strings.Join(ss, " "))
	}

	for i, _ := range data {
		lss = append(lss, fmt.Sprintf("%02X", data[i]))
		if len(lss) < cap(lss) {
			continue
		}
		if err := pl(lss); err != nil {
			return err
		}
		lss = lss[:0]
	}

	if len(lss) == 0 {
		return nil
	}
	return pl(lss)
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
