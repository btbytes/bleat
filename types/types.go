package types

import "time"

type ProcessTreeNode struct {
	PID  int
	PPID int
	Name string
}

type PortInfo struct {
	Port        int
	PID         int
	ProcessName string
	RawName     string
	Command     string
	CWD         string
	ProjectName string
	Framework   string
	Uptime      string
	StartTime   time.Time
	Status      string
	Memory      string
	GitBranch   string
	ProcessTree []ProcessTreeNode
}

type ProcessInfo struct {
	PID         int
	ProcessName string
	Command     string
	Description string
	CPU         float64
	Memory      string
	CWD         string
	ProjectName string
	Framework   string
	Uptime      string
}
