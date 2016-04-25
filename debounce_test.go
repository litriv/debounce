package main

import (
	"bytes"
	"io"
	"testing"
	"time"
)

func TestDebounce(t *testing.T) {

	d := 300 * time.Millisecond

	pIn, in := io.Pipe()
	out := new(bytes.Buffer)

	go loop(pIn, out, d)

	_, err := in.Write([]byte{'a', 'b', 'c'})
	if err != nil {
		t.Error(err)
	}

	time.Sleep(d * 2)

	got, _, err := out.ReadRune()
	if err != nil {
		t.Error(err)
	}
	if got != 'c' {
		t.Errorf("Wanted 'c', but got '%c'", got)
	}
}
