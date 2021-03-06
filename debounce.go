// Copyright (c) 2016 Jaco Esterhuizen <jaco@litriv.com>
// All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package debounce implements a signal and IO "debouncer".
package debounce // import "litriv.com/debounce"
import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Signals debounces the input signal, using duration d.  To stop listening, close the input channel; all goroutines spawned by Signals will termitae and the output channel will be closed automatically.
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

// IO debounces tokens (according to sf) received from in.  Debouncing stops at EOF or with closed reader.  All goroutines spawned by IO will terminate.
func IO(in io.Reader, out io.Writer, d time.Duration, sf bufio.SplitFunc) {
	var (
		mu sync.Mutex
		p  []byte
	)

	cin, cout := Signals(d)

	go func() {
		for range cout {
			mu.Lock()
			out.Write(p)
			out.Write([]byte("\n"))
			mu.Unlock()
		}
	}()

	s := bufio.NewScanner(in)
	s.Split(sf)

	go func() {
		defer close(cin)
		for s.Scan() {
			mu.Lock()
			p = s.Bytes()
			mu.Unlock()
			cin <- struct{}{}
		}
		if s.Err() != nil {
			printErr(s.Err())
		}
	}()
}

func printErr(err error) {
	fmt.Fprintln(os.Stderr, "debounce:", err)
}
