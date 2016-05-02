package debounce

import (
	"bufio"
	"bytes"
	"io"
	"testing"
	"time"
)

var d = 300 * time.Millisecond

func TestSignals(t *testing.T) {
	in, out := Signals(d)

	go func() {
		for i := 0; i < 3; i++ {
			in <- struct{}{}
		}
		time.Sleep(d * 2)
		close(in)
	}()

	got := 0

	for range out {
		got++
	}

	if got != 1 {
		t.Errorf("Wanted 1, but got %d", got)
	}
}

func TestRunes(t *testing.T) {
	pIn, in := io.Pipe()
	out := new(bytes.Buffer)

	IO(pIn, out, d, bufio.ScanRunes)

	_, err := in.Write([]byte("abc"))
	if err != nil {
		t.Error(err)
	}

	time.Sleep(d * 2)

	got, _, err := out.ReadRune()
	if err != nil {
		t.Error(err)
	}

	if err = pIn.Close(); err != nil {
		t.Error(err)
	}

	if got != 'c' {
		t.Errorf("Wanted 'c', but got '%c'", got)
	}
}

func TestLines(t *testing.T) {
	pIn, in := io.Pipe()
	out := new(bytes.Buffer)

	IO(pIn, out, d, bufio.ScanLines)

	_, err := in.Write([]byte("abc\ndef\nghi\n"))
	if err != nil {
		t.Error(err)
	}

	time.Sleep(d * 2)

	got := string(out.Bytes()[:3])
	if err != nil {
		t.Error(err)
	}

	if err = pIn.Close(); err != nil {
		t.Error(err)
	}

	if got != "ghi" {
		t.Errorf("Wanted 'ghi', but got '%s'", got)
	}
}
