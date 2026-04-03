package scanner

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"

	"github.com/btbytes/bleat/types"
)

func GetProcessTree(pid int) []types.ProcessTreeNode {
	processMap := getAllProcesses()
	var tree []types.ProcessTreeNode

	current := pid
	visited := make(map[int]bool)
	for current > 0 {
		if visited[current] {
			break
		}
		visited[current] = true

		proc, ok := processMap[current]
		if !ok {
			break
		}

		tree = append([]types.ProcessTreeNode{proc}, tree...)
		current = proc.PPID
	}

	return tree
}

func getAllProcesses() map[int]types.ProcessTreeNode {
	result := make(map[int]types.ProcessTreeNode)

	cmd := exec.Command("ps", "-eo", "pid=,ppid=,comm=")
	output, err := cmd.Output()
	if err != nil {
		return result
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		pid, _ := strconv.Atoi(fields[0])
		ppid, _ := strconv.Atoi(fields[1])
		name := fields[2]

		result[pid] = types.ProcessTreeNode{
			PID:  pid,
			PPID: ppid,
			Name: name,
		}
	}

	return result
}
