package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	spinner = []string{`⠋`, `⠙`, `⠹`, `⠸`, `⠼`, `⠴`, `⠦`, `⠧`, `⠇`, `⠏`}
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	ch := make(chan string, 10)

	go func() {
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		if scanner.Err() != nil {
			log.Fatal(scanner.Err())
		}
		close(ch)
	}()

	var (
		spinnerPos = 0
		line       = ""
		longest    = 0
		ok         bool
	)

	display := func() {
		// fmt.Printf("printing %q\n", line)
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
		case line, ok = <-ch:
			if !ok {
				break L
			}
			display()
		}
	}

	fmt.Printf("\r  %s\n", line)
}
