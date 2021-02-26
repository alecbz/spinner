package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	flag "github.com/spf13/pflag"
)

type spinner struct {
	s    []string
	freq time.Duration
	pos  int
}

func (s *spinner) next() {
	s.pos = (s.pos + 1) % len(s.s)
}

func (s spinner) curr() string {
	return s.s[s.pos]
}

func spin(freq time.Duration, parts ...string) spinner {
	return spinner{
		s:    parts,
		freq: freq,
		pos:  0,
	}
}

var (
	asciiSpinner   = spin(100*time.Millisecond, `|`, `/`, `-`, `\`)
	unicodeSpinner = spin(50*time.Millisecond, `⠋`, `⠙`, `⠹`, `⠸`, `⠼`, `⠴`, `⠦`, `⠧`, `⠇`, `⠏`)
)

var (
	echo           *bool   = flag.BoolP("echo", "e", false, "echo lines read from stdin")
	finalEcho      *bool   = flag.Bool("final-echo", false, "don't clear the last line displayed before exiting")
	initialMessage *string = flag.StringP("message", "m", "", "initial text to display")
)

func main() {
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	ch := make(chan string, 10)

	go func() {
		for scanner.Scan() {
			if *echo {
				ch <- scanner.Text()
			}
		}
		if scanner.Err() != nil {
			log.Fatal(scanner.Err())
		}
		close(ch)
	}()

	var (
		s    = unicodeSpinner
		line = *initialMessage
		n    = 0
	)

	clear := func() {
		fmt.Printf("\r%-[2]*[1]s\r", "", n)
	}

	display := func() {
		clear()
		n, _ = fmt.Printf("%s %s", s.curr(), line)
	}

L:
	for {
		select {
		case <-time.Tick(s.freq):
			s.next()
			display()
		case l, ok := <-ch:
			if !ok {
				break L
			}
			line = l
			display()
		}
	}

	clear()

	if *finalEcho {
		fmt.Printf("  %s\n", line)
	}
}
