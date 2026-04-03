package scanner

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/btbytes/bleat/types"
)

var lsofLineRegex = regexp.MustCompile(`^(\S+)\s+(\d+)\s+\S+\s+\S+\s+\d+\s+\S+\s+\S+\s+\S+\s+(\S+)`)

func ScanPorts() ([]types.PortInfo, error) {
	cmd := exec.Command("lsof", "-iTCP", "-sTCP:LISTEN", "-P", "-n")
	output, err := cmd.Output()
	if err != nil {
		return nil, nil
	}

	ports := make(map[int]types.PortInfo)
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	header := true
	for scanner.Scan() {
		line := scanner.Text()
		if header {
			header = false
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		rawName := fields[0]
		pidStr := fields[1]
		pid, _ := strconv.Atoi(pidStr)

		nameField := fields[8]
		port := extractPort(nameField)
		if port == 0 {
			continue
		}

		processName := normalizeProcessName(fields[0])

		if existing, ok := ports[port]; ok {
			if existing.PID == 0 || pid > 0 {
				ports[port] = types.PortInfo{
					Port:        port,
					PID:         pid,
					ProcessName: processName,
					RawName:     rawName,
				}
			}
		} else {
			ports[port] = types.PortInfo{
				Port:        port,
				PID:         pid,
				ProcessName: processName,
				RawName:     rawName,
			}
		}
	}

	var portList []types.PortInfo
	for _, p := range ports {
		portList = append(portList, p)
	}

	enrichPortInfo(portList)

	return portList, nil
}

func extractPort(field string) int {
	parts := strings.Split(field, ":")
	if len(parts) < 2 {
		return 0
	}
	portStr := parts[len(parts)-1]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return 0
	}
	return port
}

func normalizeProcessName(raw string) string {
	base := raw
	if idx := strings.LastIndex(base, "/"); idx != -1 {
		base = base[idx+1:]
	}
	if len(base) > 15 {
		base = base[:15]
	}
	return base
}

func enrichPortInfo(ports []types.PortInfo) {
	if len(ports) == 0 {
		return
	}

	var pids []string
	pidMap := make(map[string]int)
	for i, p := range ports {
		if p.PID > 0 {
			pids = append(pids, strconv.Itoa(p.PID))
			pidMap[strconv.Itoa(p.PID)] = i
		}
	}

	if len(pids) == 0 {
		return
	}

	psCmd := exec.Command("ps", "-p", strings.Join(pids, ","), "-o", "pid=,ppid=,stat=,rss=,lstart=,command=")
	psOutput, err := psCmd.Output()
	if err != nil {
		return
	}

	psScanner := bufio.NewScanner(strings.NewReader(string(psOutput)))
	for psScanner.Scan() {
		line := psScanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		pid, _ := strconv.Atoi(fields[0])
		ppid, _ := strconv.Atoi(fields[1])
		stat := fields[2]
		rss, _ := strconv.Atoi(fields[3])

		startTimeStr := strings.Join(fields[4:9], " ")
		startTime, _ := time.Parse("Mon Jan 2 15:04:05 2006", startTimeStr)

		command := strings.TrimSpace(strings.Join(fields[9:], " "))

		idx, ok := pidMap[strconv.Itoa(pid)]
		if !ok {
			continue
		}

		ports[idx].PID = pid
		ports[idx].Command = command
		ports[idx].Memory = formatMemory(rss)
		ports[idx].StartTime = startTime
		ports[idx].Uptime = formatUptime(startTime)

		if strings.Contains(stat, "Z") {
			ports[idx].Status = "zombie"
		} else if ppid == 1 && isDevProcess(ports[idx].ProcessName) {
			ports[idx].Status = "orphaned"
		} else {
			ports[idx].Status = "healthy"
		}
	}

	cwdMap := getCWDs(pids)
	for i, p := range ports {
		if cwd, ok := cwdMap[strconv.Itoa(p.PID)]; ok {
			ports[i].CWD = cwd
			ports[i].ProjectName = extractProjectName(cwd)
		}
	}

	for i := range ports {
		ports[i].Framework = DetectFramework(ports[i])
	}
}

func getCWDs(pids []string) map[string]string {
	result := make(map[string]string)
	if len(pids) == 0 {
		return result
	}

	cmd := exec.Command("lsof", "-a", "-d", "cwd", "-p", strings.Join(pids, ","))
	output, err := cmd.Output()
	if err != nil {
		return result
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
		if len(fields) < 9 {
			continue
		}
		pid := fields[1]
		cwd := fields[8]
		if strings.HasPrefix(cwd, "/") {
			result[pid] = cwd
		}
	}

	return result
}

func formatMemory(rssKB int) string {
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

func formatUptime(start time.Time) string {
	if start.IsZero() {
		return ""
	}
	duration := time.Since(start)
	days := int(duration.Hours()) / 24
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

func extractProjectName(cwd string) string {
	if cwd == "" || cwd == "/" {
		return ""
	}
	parts := strings.Split(strings.TrimSuffix(cwd, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
