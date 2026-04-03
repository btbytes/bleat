package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/btbytes/bleat/display"
	"github.com/btbytes/bleat/scanner"
)

func RunWatch() {
	display.PrintWatchHeader()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	knownPorts := make(map[int]bool)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-sigChan:
			fmt.Println()
			fmt.Println(color.HiBlackString("Stopped watching."))
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
							display.PrintNewPort(p)
							break
						}
					}
				}
			}

			for port := range knownPorts {
				if !currentPorts[port] {
					display.PrintClosedPort(port)
				}
			}

			knownPorts = currentPorts
		}
	}
}
