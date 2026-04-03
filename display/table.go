package display

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/btbytes/bleat/types"
)

var frameworkColors = map[string]*color.Color{
	"Next.js":       color.New(color.FgWhite, color.BgBlack),
	"Vite":          color.New(color.FgYellow),
	"React":         color.New(color.FgCyan),
	"Vue":           color.New(color.FgGreen),
	"Angular":       color.New(color.FgRed),
	"Svelte":        color.New(color.FgHiYellow),
	"SvelteKit":     color.New(color.FgHiYellow),
	"Express":       color.New(color.FgHiBlack),
	"Fastify":       color.New(color.FgWhite),
	"NestJS":        color.New(color.FgRed),
	"Nuxt":          color.New(color.FgGreen),
	"Remix":         color.New(color.FgBlue),
	"Astro":         color.New(color.FgMagenta),
	"Django":        color.New(color.FgGreen),
	"Flask":         color.New(color.FgWhite),
	"FastAPI":       color.New(color.FgCyan),
	"Rails":         color.New(color.FgRed),
	"Gatsby":        color.New(color.FgMagenta),
	"Go":            color.New(color.FgCyan),
	"Rust":          color.New(color.FgHiYellow),
	"Ruby":          color.New(color.FgRed),
	"Python":        color.New(color.FgYellow),
	"Node.js":       color.New(color.FgGreen),
	"Java":          color.New(color.FgRed),
	"Docker":        color.New(color.FgBlue),
	"PostgreSQL":    color.New(color.FgBlue),
	"Redis":         color.New(color.FgRed),
	"MySQL":         color.New(color.FgBlue),
	"MongoDB":       color.New(color.FgGreen),
	"nginx":         color.New(color.FgGreen),
	"LocalStack":    color.New(color.FgWhite),
	"RabbitMQ":      color.New(color.FgHiYellow),
	"Kafka":         color.New(color.FgWhite),
	"Elasticsearch": color.New(color.FgYellow),
	"MinIO":         color.New(color.FgRed),
	"Bun":           color.New(color.FgYellow),
	"Deno":          color.New(color.FgWhite),
	"PHP":           color.New(color.FgMagenta),
	"Elixir":        color.New(color.FgMagenta),
	".NET":          color.New(color.FgBlue),
}

func PrintPortsTable(ports []types.PortInfo) {
	if len(ports) == 0 {
		fmt.Println(color.YellowString("No dev ports found."))
		return
	}

	header := fmt.Sprintf("  %-8s %-14s %-8s %-18s %-16s %-10s %-10s",
		"PORT", "PROCESS", "PID", "PROJECT", "FRAMEWORK", "UPTIME", "STATUS")

	fmt.Println(color.New(color.Bold).Sprint(header))
	fmt.Println(strings.Repeat("─", len(header)))

	for _, p := range ports {
		portStr := fmt.Sprintf(":%d", p.Port)
		processStr := p.ProcessName
		pidStr := fmt.Sprintf("%d", p.PID)
		projectStr := p.ProjectName
		frameworkStr := p.Framework
		if frameworkStr == "" {
			frameworkStr = "-"
		}
		uptimeStr := p.Uptime
		if uptimeStr == "" {
			uptimeStr = "-"
		}

		statusIcon := statusIndicator(p.Status)
		frameworkColored := colorizeFramework(frameworkStr)

		line := fmt.Sprintf("  %-8s %-14s %-8s %-18s %-16s %-10s %s",
			color.CyanString(portStr),
			processStr,
			color.MagentaString(pidStr),
			projectStr,
			frameworkColored,
			uptimeStr,
			statusIcon,
		)
		fmt.Println(line)
	}

	fmt.Println()
	fmt.Println(color.HiBlackString(fmt.Sprintf("%d port(s) · bleat --all to show all · bleat ps for processes", len(ports))))
}

func statusIndicator(status string) string {
	switch status {
	case "healthy":
		return color.GreenString("● healthy")
	case "orphaned":
		return color.YellowString("● orphaned")
	case "zombie":
		return color.RedString("● zombie")
	default:
		return "● " + status
	}
}

func colorizeFramework(framework string) string {
	if framework == "-" {
		return color.HiBlackString(framework)
	}
	c, ok := frameworkColors[framework]
	if !ok {
		return framework
	}
	return c.Sprint(framework)
}

func PrintPortDetail(info types.PortInfo, tree []types.ProcessTreeNode) {
	fmt.Println()
	fmt.Println(color.New(color.Bold, color.BgCyan, color.FgBlack).Sprint(fmt.Sprintf("  :%d  ", info.Port)))
	fmt.Println()

	fmt.Println(color.New(color.Bold).Sprint("  Process"))
	fmt.Printf("    Name:      %s\n", info.ProcessName)
	fmt.Printf("    PID:       %d\n", info.PID)
	fmt.Printf("    Status:    %s\n", statusIndicator(info.Status))
	fmt.Printf("    Framework: %s\n", colorizeFrameworkOrDefault(info.Framework))
	fmt.Printf("    Memory:    %s\n", info.Memory)
	fmt.Printf("    Uptime:    %s\n", info.Uptime)
	if !info.StartTime.IsZero() {
		fmt.Printf("    Started:   %s\n", info.StartTime.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()

	fmt.Println(color.New(color.Bold).Sprint("  Location"))
	fmt.Printf("    Directory: %s\n", info.CWD)
	fmt.Printf("    Project:   %s\n", info.ProjectName)
	if info.GitBranch != "" {
		fmt.Printf("    Branch:    %s\n", color.GreenString(info.GitBranch))
	}
	fmt.Println()

	if len(tree) > 1 {
		fmt.Println(color.New(color.Bold).Sprint("  Process Tree"))
		for i, node := range tree {
			prefix := "    "
			if i > 0 {
				prefix = "    └─ "
			}
			fmt.Printf("%s%s (PID: %d)\n", prefix, node.Name, node.PID)
		}
		fmt.Println()
	}

	if info.Command != "" {
		fmt.Println(color.New(color.Bold).Sprint("  Command"))
		fmt.Printf("    %s\n", color.HiBlackString(info.Command))
		fmt.Println()
	}
}

func colorizeFrameworkOrDefault(framework string) string {
	if framework == "" {
		return color.HiBlackString("-")
	}
	return colorizeFramework(framework)
}

func PrintPsTable(processes []types.ProcessInfo) {
	if len(processes) == 0 {
		fmt.Println(color.YellowString("No dev processes found."))
		return
	}

	header := fmt.Sprintf("  %-8s %-14s %-6s %-10s %-18s %-16s %-10s %s",
		"PID", "PROCESS", "CPU%", "MEM", "PROJECT", "FRAMEWORK", "UPTIME", "WHAT")

	fmt.Println(color.New(color.Bold).Sprint(header))
	fmt.Println(strings.Repeat("─", len(header)))

	for _, p := range processes {
		frameworkStr := p.Framework
		if frameworkStr == "" {
			frameworkStr = "-"
		}
		uptimeStr := p.Uptime
		if uptimeStr == "" {
			uptimeStr = "-"
		}

		line := fmt.Sprintf("  %-8s %-14s %-6s %-10s %-18s %-16s %-10s %s",
			color.MagentaString(fmt.Sprintf("%d", p.PID)),
			p.ProcessName,
			fmt.Sprintf("%.1f", p.CPU),
			p.Memory,
			p.ProjectName,
			colorizeFramework(frameworkStr),
			uptimeStr,
			color.HiBlackString(p.Description),
		)
		fmt.Println(line)
	}

	fmt.Println()
	fmt.Println(color.HiBlackString(fmt.Sprintf("%d process(es)", len(processes))))
}

func PrintCleanList(ports []types.PortInfo) {
	if len(ports) == 0 {
		fmt.Println(color.GreenString("No orphaned or zombie processes found."))
		return
	}

	fmt.Println(color.New(color.Bold, color.FgYellow).Sprint("Found orphaned/zombie processes:"))
	fmt.Println()

	for _, p := range ports {
		frameworkStr := p.Framework
		if frameworkStr == "" {
			frameworkStr = "-"
		}
		fmt.Printf("  %s  PID: %-6d  %-14s  %-16s  %s\n",
			statusIndicator(p.Status),
			p.PID,
			p.ProcessName,
			colorizeFramework(frameworkStr),
			p.ProjectName,
		)
	}
	fmt.Println()
}

func PrintWatchHeader() {
	fmt.Println(color.New(color.Bold, color.BgCyan, color.FgBlack).Sprint("  bleat watch  "))
	fmt.Println()
	fmt.Println(color.HiBlackString("Watching for port changes... (Ctrl+C to exit)"))
	fmt.Println()
}

func PrintNewPort(info types.PortInfo) {
	frameworkStr := info.Framework
	if frameworkStr == "" {
		frameworkStr = "-"
	}
	fmt.Printf("  %s  :%-6d  %-14s  %-18s  %s\n",
		color.GreenString("▲ NEW"),
		info.Port,
		info.ProcessName,
		info.ProjectName,
		colorizeFramework(frameworkStr),
	)
}

func PrintClosedPort(port int) {
	fmt.Printf("  %s  :%d\n",
		color.RedString("▼ CLOSED"),
		port,
	)
}

func PrintKillPrompt() {
	fmt.Print(color.New(color.Bold).Sprint("Kill process? [y/N] "))
}

func PrintKillSuccess() {
	fmt.Println(color.GreenString("✓ Process killed successfully."))
}

func PrintKillFailed(err error) {
	fmt.Println(color.RedString("✗ Failed to kill process: ") + err.Error())
	fmt.Println(color.HiBlackString("Try running with sudo."))
}

func PrintCleanPrompt() {
	fmt.Print(color.New(color.Bold).Sprint("Kill all? [y/N] "))
}

func PrintCleanSummary(success, failed int) {
	fmt.Println()
	if failed == 0 {
		fmt.Println(color.GreenString(fmt.Sprintf("✓ Successfully killed %d process(es).", success)))
	} else {
		fmt.Println(color.YellowString(fmt.Sprintf("Killed %d process(es), %d failed.", success, failed)))
	}
}

func PrintHelp() {
	fmt.Println()
	fmt.Println(color.New(color.Bold, color.BgCyan, color.FgBlack).Sprint("  bleat  "))
	fmt.Println()
	fmt.Println(color.New(color.Bold).Sprint("USAGE"))
	fmt.Println("  bleat                  Show dev server ports (default)")
	fmt.Println("  bleat --all            Show all listening ports")
	fmt.Println("  bleat -a               Shorthand for --all")
	fmt.Println("  bleat <port>           Detailed view of specific port")
	fmt.Println("  bleat ps               Show all running dev processes")
	fmt.Println("  bleat ps --all         Show all processes")
	fmt.Println("  bleat clean            Kill orphaned/zombie dev servers")
	fmt.Println("  bleat watch            Monitor port changes in real-time")
	fmt.Println("  bleat help             Show this help message")
	fmt.Println()
}

func FormatUptime(startTime time.Time) string {
	if startTime.IsZero() {
		return ""
	}
	duration := time.Since(startTime)
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
