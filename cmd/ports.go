package cmd

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/btbytes/bleat/display"
	"github.com/btbytes/bleat/scanner"
	"github.com/btbytes/bleat/types"
)

func RunPorts(all bool, specificPort string) {
	ports, err := scanner.ScanPorts()
	if err != nil {
		fmt.Println(color.RedString("Error scanning ports: ") + err.Error())
		return
	}

	if specificPort != "" {
		portNum, err := strconv.Atoi(specificPort)
		if err != nil {
			fmt.Println(color.RedString("Invalid port number: ") + specificPort)
			return
		}

		var found *types.PortInfo
		for i, p := range ports {
			if p.Port == portNum {
				found = &ports[i]
				break
			}
		}

		if found == nil {
			fmt.Println(color.YellowString(fmt.Sprintf("No process found on port :%d", portNum)))
			return
		}

		enrichPortDetail(found)

		tree := scanner.GetProcessTree(found.PID)
		display.PrintPortDetail(*found, tree)

		display.PrintKillPrompt()
		var response string
		fmt.Scanln(&response)
		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			killProcess(found.PID)
		}
		return
	}

	if !all {
		ports = scanner.FilterDevPorts(ports)
	}

	display.PrintPortsTable(ports)
}

func enrichPortDetail(info *types.PortInfo) {
	if info.CWD != "" {
		cmd := exec.Command("git", "-C", info.CWD, "rev-parse", "--abbrev-ref", "HEAD")
		output, err := cmd.Output()
		if err == nil {
			info.GitBranch = strings.TrimSpace(string(output))
		}
	}
}

func killProcess(pid int) {
	cmd := exec.Command("kill", fmt.Sprintf("%d", pid))
	err := cmd.Run()
	if err != nil {
		display.PrintKillFailed(err)
	} else {
		display.PrintKillSuccess()
	}
}

func RunHelp() {
	display.PrintHelp()
}
