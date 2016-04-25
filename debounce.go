package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
	"unicode/utf8"
)

func main() {
	d := flag.Int64("i", 300, "duration in milliseconds after last action, after which function executes")
	loop(os.Stdin, os.Stdout, time.Duration(*d)*time.Millisecond)
}



func loop(in io.Reader, out io.Writer, d time.Duration) {
	var (
		t = time.NewTimer(d)
		c rune
	)

	go func() {
		p := make([]byte, 4)
		for {
			<-t.C
			n := utf8.EncodeRune(p, c)
			out.Write(p[:n])
		}
	}()

	r := bufio.NewReader(in)
	var err error

	for {
		if c, _, err = r.ReadRune(); err != nil {
			printErr(err)
		}
	
		t.Reset(d)
	}
}

func printErr(err error) {
	fmt.Println("debounce:", err)
}
