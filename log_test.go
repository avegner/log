package log

import (
	"testing"

	"github.com/avegner/log/out"
)

var (
	mask = ALL_LEVELS
)

func TestRootLog(t *testing.T) {
	l := createRootLogger(t, "main", mask)

	l.Printf(INFO, "Hello, world! My name is %s", "Aleksei Vegner")
	l.Printf(CRITICAL, "got some crtitical error")
}

func TestChildLog(t *testing.T) {
	l := createRootLogger(t, "main", mask)
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

func createRootLogger(t *testing.T, name string, mask Level) Logger {
	outs := createOuts(t)
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
