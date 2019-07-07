# log

This a logger supporting the following outputs:
* stderr
* network (TCP, UDP, Unix networks)
* file (zlib compression)
* custom output

# API

Logger has a very simple interface:

```go
type Logger interface {
	Printf(level Level, format string, args ...interface{}) error
	SetMask(mask Level)
	Flush() error
	Child(name string) Logger
}
```

* `Printf` - prints a log message with the given level
* `SetMask` - sets enabled log levels
* `Flush` - flushes all log outputs
* `Child` - creates a child logger with the same parameters (mask, outputs, flags)

Each output supports the following interface:

```go
type Outputter interface {
	io.Writer
	io.Closer
	Flush() error
}
```

There are standard outputs:

```go
func NewNetOut(network, address string) (Outputter, error)

func NewFileOut(name string, perm os.FileMode, append bool, comprLevel int) (Outputter, error)

func NewStderrOut() (Outputter, error)
```

Any custom output implementing `Outputter` interface may be created.

To create a logger and a child logger:

```go
// main scope
o, err := out.NewStderrOut()
if err != nil {
    // ...
}
mlog := log.New("main", log.ALL_LEVELS, log.STD_FLAGS, []out.Outputter{o})

// module scope
chlog := mlog.Child("module")
```
