package log

type Logger interface {
	Printf(level Level, format string, args ...interface{})
	Dump(level Level, addr uintptr, size uint)
	Flush()
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
