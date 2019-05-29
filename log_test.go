package log

import (
	"testing"

	"github.com/avegner/log/out"
	"github.com/avegner/log/out/stderr"
)

func TestRootLog(t *testing.T) {
	l := createRootLogger(t, "main")

	l.Printf(INFO, "Hello, world! My name is %s", "Aleksei Vegner")
	l.Printf(CRITICAL, "got some crtitical error")
}

func TestChildLog(t *testing.T) {
	l := createRootLogger(t, "main")
	cl := createChildLogger(t, l, "mod")

	cl.Printf(INFO, "Hello, world! My name is %s", "Aleksei Vegner")
	cl.Printf(CRITICAL, "got some crtitical error")
}

func createOuts(t *testing.T) []out.Outputter {
	stderr, err := stderr.New()

	if err != nil {
		t.Fatal(err)
	}
	return []out.Outputter{stderr}
}

func createRootLogger(t *testing.T, name string) Logger {
	outs := createOuts(t)
	l := New(name, outs)

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
