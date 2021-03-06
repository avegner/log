package out

import (
	"compress/zlib"
	"fmt"
	"testing"
)

const (
	logsNumber = 1000
)

func BenchmarkStderrOut(b *testing.B) {
	o, err := NewStderrOut()
	if err != nil {
		b.Fatalf("NewStderrOut(): got '%v' error, want no error", err)
	}
	defer o.Close()

	benchOut(b, o)
}

func BenchmarkFileOut(b *testing.B) {
	o, err := NewFileOut(fileName, 0644, false, zlib.NoCompression)
	if err != nil {
		b.Fatalf("NewFileOut(): got '%v' error, want no error", err)
	}
	defer o.Close()

	benchOut(b, o)
}

func BenchmarkNetOut(b *testing.B) {
	o, err := NewNetOut(network, address)
	if err != nil {
		b.Fatalf("NewNetOut(): got '%v' error, want no error", err)
	}
	defer o.Close()

	benchOut(b, o)
}

func benchOut(b *testing.B, o Outputter) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < logsNumber; i++ {
			_, _ = o.Write([]byte(fmt.Sprintf("benchmark record number %d\n", i)))
		}
	}
	b.StopTimer()
}
