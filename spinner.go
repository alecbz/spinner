package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	spinner = []string{`⠋`, `⠙`, `⠹`, `⠸`, `⠼`, `⠴`, `⠦`, `⠧`, `⠇`, `⠏`}

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
		spinnerPos = 0
		line       = *initialMessage
		longest    = len(line)
	)

	display := func() {
		fmt.Printf("\r%-[2]*[1]s\r", "", longest) // clear the line
		n, _ := fmt.Printf("%s %s", spinner[spinnerPos], line)
		if n > longest {
			longest = n
		}
	}

L:
	for {
		select {
		case <-time.Tick(50 * time.Millisecond):
			spinnerPos = (spinnerPos + 1) % len(spinner)
			display()
		case l, ok := <-ch:
			if !ok {
				break L
			}
			line = l
			display()
		}
	}

	fmt.Printf("\r%-[2]*[1]s\r", "", longest) // clear the line

	if *finalEcho {
		fmt.Printf("  %s\n", line)
	}
}
