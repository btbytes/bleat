package cmd

import (
	"bufio"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/btbytes/bleat/display"
	"github.com/btbytes/bleat/scanner"
	"github.com/btbytes/bleat/types"
)

func RunPs(all bool) {
	cmd := exec.Command("ps", "-eo", "pid=,pcpu=,pmem=,rss=,lstart=,command=")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println(color.RedString("Error getting process list: ") + err.Error())
		return
	}

	var processes []types.ProcessInfo
	psScanner := bufio.NewScanner(strings.NewReader(string(output)))
	for psScanner.Scan() {
		line := strings.TrimSpace(psScanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}

		pid, _ := strconv.Atoi(fields[0])
		cpu, _ := strconv.ParseFloat(fields[1], 64)
		rss, _ := strconv.Atoi(fields[3])

		startTimeStr := strings.Join(fields[4:9], " ")
		startTime, _ := time.Parse("Mon Jan 2 15:04:05 2006", startTimeStr)

		command := strings.TrimSpace(strings.Join(fields[9:], " "))

		processName := fields[len(fields)-1]
		if idx := strings.LastIndex(processName, "/"); idx != -1 {
			processName = processName[idx+1:]
		}
		if len(processName) > 15 {
			processName = processName[:15]
		}

		if !all && scanner.IsSystemApp(processName) {
			continue
		}

		cwd := getCWDForPID(pid)
		projectName := extractProjectNameFromPath(cwd)
		framework := detectFrameworkForProcess(command, cwd, processName)

		description := command
		if len(description) > 60 {
			description = description[:57] + "..."
		}

		processes = append(processes, types.ProcessInfo{
			PID:         pid,
			ProcessName: processName,
			Command:     command,
			Description: description,
			CPU:         cpu,
			Memory:      formatMem(rss),
			CWD:         cwd,
			ProjectName: projectName,
			Framework:   framework,
			Uptime:      display.FormatUptime(startTime),
		})
	}

	sort.Slice(processes, func(i, j int) bool {
		return processes[i].CPU > processes[j].CPU
	})

	collapseDocker(&processes)

	display.PrintPsTable(processes)
}

func getCWDForPID(pid int) string {
	cmd := exec.Command("lsof", "-a", "-d", "cwd", "-p", strconv.Itoa(pid))
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	header := true
	for scanner.Scan() {
		line := scanner.Text()
		if header {
			header = false
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 9 {
			return fields[8]
		}
	}
	return ""
}

func extractProjectNameFromPath(cwd string) string {
	if cwd == "" || cwd == "/" {
		return ""
	}
	parts := strings.Split(strings.TrimSuffix(cwd, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func detectFrameworkForProcess(command, cwd, processName string) string {
	info := types.PortInfo{
		Command:     command,
		CWD:         cwd,
		ProcessName: processName,
	}
	return scanner.DetectFramework(info)
}

func formatMem(rssKB int) string {
	if rssKB < 1024 {
		return fmt.Sprintf("%d KB", rssKB)
	}
	mb := float64(rssKB) / 1024.0
	if mb < 1024 {
		return fmt.Sprintf("%.0f MB", mb)
	}
	gb := mb / 1024.0
	return fmt.Sprintf("%.1f GB", gb)
}

func collapseDocker(processes *[]types.ProcessInfo) {
	var dockerProcs []types.ProcessInfo
	var others []types.ProcessInfo

	for _, p := range *processes {
		if scanner.IsDockerProcess(p.ProcessName) || strings.Contains(strings.ToLower(p.Command), "docker") {
			dockerProcs = append(dockerProcs, p)
		} else {
			others = append(others, p)
		}
	}

	if len(dockerProcs) > 1 {
		totalCPU := 0.0
		totalRSS := 0
		for _, p := range dockerProcs {
			totalCPU += p.CPU
			rss := parseMemToKB(p.Memory)
			totalRSS += rss
		}

		others = append(others, types.ProcessInfo{
			PID:         dockerProcs[0].PID,
			ProcessName: "docker",
			Description: fmt.Sprintf("%d container processes", len(dockerProcs)),
			CPU:         totalCPU,
			Memory:      formatMem(totalRSS),
			ProjectName: "Docker",
			Framework:   "Docker",
			Uptime:      dockerProcs[0].Uptime,
		})
	} else if len(dockerProcs) == 1 {
		others = append(others, dockerProcs[0])
	}

	*processes = others
}

func parseMemToKB(mem string) int {
	parts := strings.Fields(mem)
	if len(parts) < 2 {
		return 0
	}
	val, _ := strconv.ParseFloat(parts[0], 64)
	unit := parts[1]
	switch unit {
	case "KB":
		return int(val)
	case "MB":
		return int(val * 1024)
	case "GB":
		return int(val * 1024 * 1024)
	}
	return 0
}
