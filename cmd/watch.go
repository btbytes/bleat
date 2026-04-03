package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/btbytes/bleat/display"
	"github.com/btbytes/bleat/scanner"
)

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func RunWatch() {
	clearScreen()
	display.PrintWatchHeader()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	knownPorts := make(map[int]bool)
	var events []display.WatchEvent

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println()
			if len(events) == 0 {
				fmt.Println(display.S.Dim.Render("No port changes detected."))
			} else {
				display.PrintWatchEvents(events)
			}
			return
		case <-ticker.C:
			ports, err := scanner.ScanPorts()
			if err != nil {
				continue
			}

			currentPorts := make(map[int]bool)
			for _, p := range ports {
				currentPorts[p.Port] = true
			}

			for port := range currentPorts {
				if !knownPorts[port] {
					for _, p := range ports {
						if p.Port == port {
							events = append(events, display.WatchEvent{Kind: "new", Info: p})
							break
						}
					}
				}
			}

			for port := range knownPorts {
				if !currentPorts[port] {
					events = append(events, display.WatchEvent{Kind: "closed", Port: port})
				}
			}

			clearScreen()
			display.PrintWatchHeader()
			display.PrintWatchEvents(events)

			knownPorts = currentPorts
		}
	}
}
