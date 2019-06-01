package log

import (
	"testing"

	"github.com/avegner/log/out"
)

var (
	mask = ALL_LEVELS
)

func TestRootLog(t *testing.T) {
	outs := createOuts(t)
	defer closeOuts(t, outs)
	l := createRootLogger(t, "main", mask, outs)

	l.Printf(INFO, "Hello, world! My name is %s", "Aleksei Vegner")
	l.Printf(CRITICAL, "got some crtitical error")
}

func TestChildLog(t *testing.T) {
	outs := createOuts(t)
	defer closeOuts(t, outs)
	l := createRootLogger(t, "main", mask, outs)
	cl := createChildLogger(t, l, "mod")

	cl.Printf(INFO, "Hello, world! My name is %s", "Aleksei Vegner")
	cl.Printf(CRITICAL, "got some crtitical error")
}

func createOuts(t *testing.T) []out.Outputter {
	stderr, err := out.NewStderrOut()
	if err != nil {
		t.Fatal(err)
	}
	return []out.Outputter{stderr}
}

func closeOuts(t *testing.T, outs []out.Outputter) {
	for _, o := range outs {
		if err := o.Close(); err != nil {
			t.Fatalf("Close(): got '%v' error, want no error", err)
		}
	}
}

func createRootLogger(t *testing.T, name string, mask Level, outs []out.Outputter) Logger {
	l := New(name, mask, STD_FLAGS, outs)
	if l == nil {
		t.Fatalf("New(): got nil, want no nil")
	}
	return l
}

func createChildLogger(t *testing.T, l Logger, name string) Logger {
	cl := l.Child(name)
	if cl == nil {
		t.Fatalf("Child(): got nil, want no nil")
	}
	return cl
}
