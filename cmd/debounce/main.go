package main

import (
	"bufio"
	"flag"
	"litriv.com/debounce"
	"os"
	"os/signal"
	"time"
)

func main() {
	d := flag.Int64("i", 300, "duration in milliseconds after last action, after which function executes")
	s := flag.String("s", "lines", "'lines' to split on lines, 'runes' to split on runes")
	var sf bufio.SplitFunc
	switch *s {
	case "runes":
		sf = bufio.ScanRunes
	default:
		sf = bufio.ScanLines
	}
	debounce.IO(os.Stdin, os.Stdout, time.Duration(*d)*time.Millisecond, sf)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
