package main

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func TestSignals(t *testing.T) {
	d := 300 * time.Millisecond
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

	d := 300 * time.Millisecond

	pIn, in := io.Pipe()
	out := new(bytes.Buffer)

	x := Runes(pIn, out, d)

	_, err := in.Write([]byte{'a', 'b', 'c'})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(d * 2)

	got, _, err := out.ReadRune()
	if err != nil {
		t.Error(err)
	}

	x()

	if got != 'c' {
		t.Errorf("Wanted 'c', but got '%c'", got)
	}
}
