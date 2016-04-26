package main

import (
	"flag"
	"os"
	"time"
	"litriv.com/debounce"
)

func main() {
	d := flag.Int64("i", 300, "duration in milliseconds after last action, after which function executes")
	debounce.Runes(os.Stdin, os.Stdout, time.Duration(*d)*time.Millisecond)
}