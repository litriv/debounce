package debounce // import "litriv.com/debounce"

import (
	"bufio"
	"fmt"
	"io"
	"time"
	"unicode/utf8"
)

// Signals debounces the input signal, using duration d.  To stop listening, close the input channel; the output channel will be closed automatically.
func Signals(d time.Duration) (chan<- struct{}, <-chan struct{}) {
	in, out := make(chan struct{}), make(chan struct{})
	t := time.NewTimer(time.Hour)
	t.Stop()

	exit := make(chan struct{})

	go func() {
		defer close(out)

		for x := false; !x; {
			select {
			case <-t.C:
				out <- struct{}{}
			case <-exit:
				x = true
			}
		}
	}()

	go func() {
		for range in {
			t.Reset(d)
		}
		exit <- struct{}{}
	}()

	return in, out
}

// Runes debounces runes read from in. To stop debouncing, call the returned function.
func Runes(in io.Reader, out io.Writer, d time.Duration) func() {
	var c rune

	cin, cout := Signals(d)

	go func() {
		p := make([]byte, 4)
		for range cout {
			n := utf8.EncodeRune(p, c)
			out.Write(p[:n])
		}
	}()

	r := bufio.NewReader(in)
	var err error

	go func() {
		for {
			if c, _, err = r.ReadRune(); err != nil {
				printErr(err)
			}
			cin <- struct{}{}
		}
	}()

	return func() {
		close(cin)
	}
}

func printErr(err error) {
	fmt.Println("debounce:", err)
}
