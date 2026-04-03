package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/btbytes/bleat/display"
	"github.com/btbytes/bleat/scanner"
	"github.com/btbytes/bleat/types"
)

func RunClean() {
	ports, err := scanner.ScanPorts()
	if err != nil {
		fmt.Println(color.RedString("Error scanning ports: ") + err.Error())
		return
	}

	var orphaned []types.PortInfo
	for _, p := range ports {
		if p.Status == "orphaned" || p.Status == "zombie" {
			orphaned = append(orphaned, p)
		}
	}

	display.PrintCleanList(orphaned)

	if len(orphaned) == 0 {
		return
	}

	display.PrintCleanPrompt()
	var response string
	fmt.Scanln(&response)
	response = strings.ToLower(strings.TrimSpace(response))
	if response != "y" && response != "yes" {
		fmt.Println(color.HiBlackString("Aborted."))
		return
	}

	success := 0
	failed := 0
	for _, p := range orphaned {
		cmd := exec.Command("kill", fmt.Sprintf("%d", p.PID))
		err := cmd.Run()
		if err != nil {
			cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", p.PID))
			err = cmd.Run()
			if err != nil {
				failed++
				continue
			}
		}
		success++
	}

	display.PrintCleanSummary(success, failed)
}
