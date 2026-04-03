package display

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/btbytes/bleat/types"
	"github.com/charmbracelet/lipgloss"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripAnsi(s string) string {
	return ansiRegexp.ReplaceAllString(s, "")
}

func visibleLen(s string) int {
	return len(stripAnsi(s))
}

func padRight(s string, width int) string {
	vl := visibleLen(s)
	if vl >= width {
		return s
	}
	return s + strings.Repeat(" ", width-vl)
}

var frameworkColors = map[string]lipgloss.Color{
	"Next.js":       lipgloss.Color("15"),
	"Vite":          lipgloss.Color("220"),
	"React":         lipgloss.Color("51"),
	"Vue":           lipgloss.Color("46"),
	"Angular":       lipgloss.Color("196"),
	"Svelte":        lipgloss.Color("227"),
	"SvelteKit":     lipgloss.Color("227"),
	"Express":       lipgloss.Color("244"),
	"Fastify":       lipgloss.Color("15"),
	"NestJS":        lipgloss.Color("196"),
	"Nuxt":          lipgloss.Color("46"),
	"Remix":         lipgloss.Color("33"),
	"Astro":         lipgloss.Color("163"),
	"Django":        lipgloss.Color("46"),
	"Flask":         lipgloss.Color("15"),
	"FastAPI":       lipgloss.Color("51"),
	"Rails":         lipgloss.Color("196"),
	"Gatsby":        lipgloss.Color("163"),
	"Go":            lipgloss.Color("51"),
	"Rust":          lipgloss.Color("227"),
	"Ruby":          lipgloss.Color("196"),
	"Python":        lipgloss.Color("220"),
	"Node.js":       lipgloss.Color("46"),
	"Java":          lipgloss.Color("196"),
	"Docker":        lipgloss.Color("33"),
	"PostgreSQL":    lipgloss.Color("33"),
	"Redis":         lipgloss.Color("196"),
	"MySQL":         lipgloss.Color("33"),
	"MongoDB":       lipgloss.Color("46"),
	"nginx":         lipgloss.Color("46"),
	"LocalStack":    lipgloss.Color("15"),
	"RabbitMQ":      lipgloss.Color("227"),
	"Kafka":         lipgloss.Color("15"),
	"Elasticsearch": lipgloss.Color("220"),
	"MinIO":         lipgloss.Color("196"),
	"Bun":           lipgloss.Color("220"),
	"Deno":          lipgloss.Color("15"),
	"PHP":           lipgloss.Color("163"),
	"Elixir":        lipgloss.Color("163"),
	".NET":          lipgloss.Color("33"),
}

var (
	styleTableBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(0, 1)

	styleHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("63")).
			Padding(0, 1)

	stylePort = lipgloss.NewStyle().
			Foreground(lipgloss.Color("51")).
			Bold(true)

	stylePID = lipgloss.NewStyle().
			Foreground(lipgloss.Color("165"))

	styleProcess = lipgloss.NewStyle().
			Foreground(lipgloss.Color("189"))

	styleProject = lipgloss.NewStyle().
			Foreground(lipgloss.Color("223"))

	styleUptime = lipgloss.NewStyle().
			Foreground(lipgloss.Color("248"))

	styleDim = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	styleStatusHealthy = lipgloss.NewStyle().
				Foreground(lipgloss.Color("46"))

	styleStatusOrphaned = lipgloss.NewStyle().
				Foreground(lipgloss.Color("220"))

	styleStatusZombie = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196"))

	styleDetailTitle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("15")).
				Background(lipgloss.Color("63")).
				Padding(0, 1)

	styleDetailSection = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("165")).
				MarginTop(1).
				MarginBottom(0)

	styleDetailLabel = lipgloss.NewStyle().
				Foreground(lipgloss.Color("248")).
				Width(12).
				Align(lipgloss.Right)

	styleCleanBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("220")).
				Padding(0, 1)

	styleWatchNew = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	styleWatchClosed = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196")).
				Bold(true)

	styleHelpTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.Color("63")).
			Padding(0, 1)

	styleHelpCmd = lipgloss.NewStyle().
			Foreground(lipgloss.Color("51"))

	styleHelpDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("248"))
)

func statusStyled(status string) string {
	switch status {
	case "healthy":
		return styleStatusHealthy.Render("● healthy")
	case "orphaned":
		return styleStatusOrphaned.Render("● orphaned")
	case "zombie":
		return styleStatusZombie.Render("● zombie")
	default:
		return "● " + status
	}
}

func statusPlain(status string) string {
	switch status {
	case "healthy":
		return "● healthy"
	case "orphaned":
		return "● orphaned"
	case "zombie":
		return "● zombie"
	default:
		return "● " + status
	}
}

func colorizeFramework(framework string) string {
	if framework == "-" {
		return styleDim.Render("-")
	}
	c, ok := frameworkColors[framework]
	if !ok {
		return framework
	}
	return lipgloss.NewStyle().Foreground(c).Render(framework)
}

func colorizeFrameworkOrDefault(framework string) string {
	if framework == "" {
		return styleDim.Render("-")
	}
	return colorizeFramework(framework)
}

func terminalWidth() int {
	return 0
}

func renderTable(headers []string, rows [][]string, colWidths []int) string {
	headerCells := make([]string, len(headers))
	for i, h := range headers {
		headerCells[i] = styleHeader.Render(padRight(h, colWidths[i]))
	}
	headerRow := strings.Join(headerCells, " ")

	rowCells := make([]string, len(rows))
	for ri, row := range rows {
		cells := make([]string, len(row))
		for i, cell := range row {
			cells[i] = padRight(cell, colWidths[i])
		}
		rowCells[ri] = strings.Join(cells, " ")
	}

	content := headerRow + "\n" + styleDim.Render(strings.Repeat("─", visibleLen(headerRow))) + "\n"
	for _, row := range rowCells {
		content += row + "\n"
	}

	return styleTableBorder.Render(strings.TrimSuffix(content, "\n"))
}

func PrintPortsTable(ports []types.PortInfo) {
	if len(ports) == 0 {
		fmt.Println(styleDim.Render("No dev ports found."))
		return
	}

	headers := []string{"PORT", "PROCESS", "PID", "PROJECT", "FRAMEWORK", "UPTIME", "STATUS"}
	colWidths := []int{8, 14, 8, 18, 16, 10, 10}

	rows := make([][]string, len(ports))
	for i, p := range ports {
		portStr := fmt.Sprintf(":%d", p.Port)
		pidStr := fmt.Sprintf("%d", p.PID)
		frameworkStr := p.Framework
		if frameworkStr == "" {
			frameworkStr = "-"
		}
		uptimeStr := p.Uptime
		if uptimeStr == "" {
			uptimeStr = "-"
		}

		vals := []string{
			stylePort.Render(portStr),
			styleProcess.Render(p.ProcessName),
			stylePID.Render(pidStr),
			styleProject.Render(p.ProjectName),
			colorizeFramework(frameworkStr),
			styleUptime.Render(uptimeStr),
			statusStyled(p.Status),
		}

		plainVals := []string{portStr, p.ProcessName, pidStr, p.ProjectName, frameworkStr, uptimeStr, statusPlain(p.Status)}
		for j, v := range plainVals {
			if len(v) > colWidths[j] {
				colWidths[j] = len(v)
			}
		}

		rows[i] = vals
	}

	fmt.Println(renderTable(headers, rows, colWidths))
	fmt.Println()
	fmt.Println(styleDim.Render(fmt.Sprintf("%d port(s) · bleat --all to show all · bleat ps for processes", len(ports))))
}

func PrintPortDetail(info types.PortInfo, tree []types.ProcessTreeNode) {
	fmt.Println()
	fmt.Println(styleDetailTitle.Render(fmt.Sprintf("  :%d  ", info.Port)))
	fmt.Println()

	fmt.Println(styleDetailSection.Render("PROCESS"))
	fmt.Printf("  %s  %s\n", styleDetailLabel.Render("Name:"), info.ProcessName)
	fmt.Printf("  %s  %s\n", styleDetailLabel.Render("PID:"), stylePID.Render(fmt.Sprintf("%d", info.PID)))
	fmt.Printf("  %s  %s\n", styleDetailLabel.Render("Status:"), statusStyled(info.Status))
	fmt.Printf("  %s  %s\n", styleDetailLabel.Render("Framework:"), colorizeFrameworkOrDefault(info.Framework))
	fmt.Printf("  %s  %s\n", styleDetailLabel.Render("Memory:"), info.Memory)
	fmt.Printf("  %s  %s\n", styleDetailLabel.Render("Uptime:"), styleUptime.Render(info.Uptime))
	if !info.StartTime.IsZero() {
		fmt.Printf("  %s  %s\n", styleDetailLabel.Render("Started:"), info.StartTime.Format("2006-01-02 15:04:05"))
	}

	fmt.Println()
	fmt.Println(styleDetailSection.Render("LOCATION"))
	fmt.Printf("  %s  %s\n", styleDetailLabel.Render("Directory:"), info.CWD)
	fmt.Printf("  %s  %s\n", styleDetailLabel.Render("Project:"), styleProject.Render(info.ProjectName))
	if info.GitBranch != "" {
		fmt.Printf("  %s  %s\n", styleDetailLabel.Render("Branch:"), lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(info.GitBranch))
	}

	if len(tree) > 1 {
		fmt.Println()
		fmt.Println(styleDetailSection.Render("PROCESS TREE"))
		for i, node := range tree {
			if i == 0 {
				fmt.Printf("  %s %s\n", stylePID.Render(fmt.Sprintf("%d", node.PID)), node.Name)
			} else {
				fmt.Printf("  %s %s %s\n", styleDim.Render("└─"), node.Name, stylePID.Render(fmt.Sprintf("(PID: %d)", node.PID)))
			}
		}
	}

	if info.Command != "" {
		fmt.Println()
		fmt.Println(styleDetailSection.Render("COMMAND"))
		fmt.Printf("  %s\n", styleDim.Render(info.Command))
	}

	fmt.Println()
}

func PrintPsTable(processes []types.ProcessInfo) {
	if len(processes) == 0 {
		fmt.Println(styleDim.Render("No dev processes found."))
		return
	}

	headers := []string{"PID", "PROCESS", "CPU%", "MEM", "PROJECT", "FRAMEWORK", "UPTIME", "WHAT"}
	colWidths := []int{8, 14, 6, 10, 18, 16, 10}

	rows := make([][]string, len(processes))
	for i, p := range processes {
		pidStr := fmt.Sprintf("%d", p.PID)
		frameworkStr := p.Framework
		if frameworkStr == "" {
			frameworkStr = "-"
		}
		uptimeStr := p.Uptime
		if uptimeStr == "" {
			uptimeStr = "-"
		}
		cpuStr := fmt.Sprintf("%.1f", p.CPU)

		plainVals := []string{pidStr, p.ProcessName, cpuStr, p.Memory, p.ProjectName, frameworkStr, uptimeStr}
		for j, v := range plainVals {
			if len(v) > colWidths[j] {
				colWidths[j] = len(v)
			}
		}

		vals := []string{
			stylePID.Render(pidStr),
			styleProcess.Render(p.ProcessName),
			cpuStr,
			p.Memory,
			styleProject.Render(p.ProjectName),
			colorizeFramework(frameworkStr),
			styleUptime.Render(uptimeStr),
			styleDim.Render(p.Description),
		}

		rows[i] = vals
	}

	fmt.Println(renderTable(headers, rows, colWidths))
	fmt.Println()
	fmt.Println(styleDim.Render(fmt.Sprintf("%d process(es)", len(processes))))
}

func PrintCleanList(ports []types.PortInfo) {
	if len(ports) == 0 {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("No orphaned or zombie processes found."))
		return
	}

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220")).Render("Found orphaned/zombie processes:"))
	fmt.Println()

	var lines []string
	for _, p := range ports {
		frameworkStr := p.Framework
		if frameworkStr == "" {
			frameworkStr = "-"
		}
		line := fmt.Sprintf("%s  PID: %s  %s  %s  %s",
			statusStyled(p.Status),
			stylePID.Render(fmt.Sprintf("%d", p.PID)),
			styleProcess.Render(p.ProcessName),
			colorizeFramework(frameworkStr),
			styleProject.Render(p.ProjectName),
		)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	fmt.Println(styleCleanBorder.Render(content))
	fmt.Println()
}

func PrintWatchHeader() {
	fmt.Println(styleDetailTitle.Render("  bleat watch  "))
	fmt.Println()
	fmt.Println(styleDim.Render("Watching for port changes... (Ctrl+C to exit)"))
	fmt.Println()
}

func PrintNewPort(info types.PortInfo) {
	frameworkStr := info.Framework
	if frameworkStr == "" {
		frameworkStr = "-"
	}
	fmt.Printf("  %s  %s  %s  %s  %s\n",
		styleWatchNew.Render("▲ NEW"),
		stylePort.Render(fmt.Sprintf(":%d", info.Port)),
		styleProcess.Render(info.ProcessName),
		styleProject.Render(info.ProjectName),
		colorizeFramework(frameworkStr),
	)
}

func PrintClosedPort(port int) {
	fmt.Printf("  %s  %s\n",
		styleWatchClosed.Render("▼ CLOSED"),
		stylePort.Render(fmt.Sprintf(":%d", port)),
	)
}

func PrintKillPrompt() {
	fmt.Print(lipgloss.NewStyle().Bold(true).Render("Kill process? [y/N] "))
}

func PrintKillSuccess() {
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("✓ Process killed successfully."))
}

func PrintKillFailed(err error) {
	fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render("✗ Failed to kill process: ") + err.Error())
	fmt.Println(styleDim.Render("Try running with sudo."))
}

func PrintCleanPrompt() {
	fmt.Print(lipgloss.NewStyle().Bold(true).Render("Kill all? [y/N] "))
}

func PrintCleanSummary(success, failed int) {
	fmt.Println()
	if failed == 0 {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(fmt.Sprintf("✓ Successfully killed %d process(es).", success)))
	} else {
		fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color("220")).Render(fmt.Sprintf("Killed %d process(es), %d failed.", success, failed)))
	}
}

func PrintHelp() {
	fmt.Println()
	fmt.Println(styleHelpTitle.Render("  bleat  "))
	fmt.Println()

	commands := [][]string{
		{"bleat", "Show dev server ports (default)"},
		{"bleat --all", "Show all listening ports"},
		{"bleat -a", "Shorthand for --all"},
		{"bleat <port>", "Detailed view of specific port"},
		{"bleat ps", "Show all running dev processes"},
		{"bleat ps --all", "Show all processes"},
		{"bleat clean", "Kill orphaned/zombie dev servers"},
		{"bleat watch", "Monitor port changes in real-time"},
		{"bleat help", "Show this help message"},
	}

	cmdWidth := 0
	for _, c := range commands {
		if len(c[0]) > cmdWidth {
			cmdWidth = len(c[0])
		}
	}

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("165")).Render("USAGE"))
	fmt.Println()
	for _, c := range commands {
		fmt.Printf("  %s  %s\n",
			styleHelpCmd.Render(padRight(c[0], cmdWidth)),
			styleHelpDesc.Render(c[1]))
	}
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
