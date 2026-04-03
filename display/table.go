package display

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/btbytes/bleat/config"
	"github.com/btbytes/bleat/theme"
	"github.com/btbytes/bleat/types"
	"github.com/charmbracelet/lipgloss"
	gogh "github.com/willyv3/gogh-themes/lipgloss"
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

var S theme.Styles

func init() {
	cfg, err := config.Load()
	if err != nil {
		S, _ = theme.BuildStyles("")
		return
	}
	name := cfg.Theme
	if name == "" {
		name = theme.DefaultThemeName()
	}
	S, _ = theme.BuildStyles(name)
}

func statusStyled(status string) string {
	switch status {
	case "healthy":
		return S.StatusHealthy.Render("● healthy")
	case "orphaned":
		return S.StatusOrphaned.Render("● orphaned")
	case "zombie":
		return S.StatusZombie.Render("● zombie")
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
		return S.Dim.Render("-")
	}
	c, ok := S.FrameworkColors[framework]
	if !ok {
		return framework
	}
	return lipgloss.NewStyle().Foreground(c).Render(framework)
}

func colorizeFrameworkOrDefault(framework string) string {
	if framework == "" {
		return S.Dim.Render("-")
	}
	return colorizeFramework(framework)
}

func terminalWidth() int {
	return 0
}

func renderTable(headers []string, rows [][]string, colWidths []int) string {
	headerCells := make([]string, len(headers))
	for i, h := range headers {
		headerCells[i] = S.Header.Render(padRight(h, colWidths[i]))
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

	content := headerRow + "\n" + S.Dim.Render(strings.Repeat("─", visibleLen(headerRow))) + "\n"
	for _, row := range rowCells {
		content += row + "\n"
	}

	return S.TableBorder.Render(strings.TrimSuffix(content, "\n"))
}

func PrintPortsTable(ports []types.PortInfo) {
	if len(ports) == 0 {
		fmt.Println(S.Dim.Render("No dev ports found."))
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
			S.Port.Render(portStr),
			S.Process.Render(p.ProcessName),
			S.PID.Render(pidStr),
			S.Project.Render(p.ProjectName),
			colorizeFramework(frameworkStr),
			S.Uptime.Render(uptimeStr),
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
	fmt.Println(S.Dim.Render(fmt.Sprintf("%d port(s) · bleat --all to show all · bleat ps for processes", len(ports))))
}

func PrintPortDetail(info types.PortInfo, tree []types.ProcessTreeNode) {
	fmt.Println()
	fmt.Println(S.DetailTitle.Render(fmt.Sprintf("  :%d  ", info.Port)))
	fmt.Println()

	fmt.Println(S.DetailSection.Render("PROCESS"))
	fmt.Printf("  %s  %s\n", S.DetailLabel.Render("Name:"), info.ProcessName)
	fmt.Printf("  %s  %s\n", S.DetailLabel.Render("PID:"), S.PID.Render(fmt.Sprintf("%d", info.PID)))
	fmt.Printf("  %s  %s\n", S.DetailLabel.Render("Status:"), statusStyled(info.Status))
	fmt.Printf("  %s  %s\n", S.DetailLabel.Render("Framework:"), colorizeFrameworkOrDefault(info.Framework))
	fmt.Printf("  %s  %s\n", S.DetailLabel.Render("Memory:"), info.Memory)
	fmt.Printf("  %s  %s\n", S.DetailLabel.Render("Uptime:"), S.Uptime.Render(info.Uptime))
	if !info.StartTime.IsZero() {
		fmt.Printf("  %s  %s\n", S.DetailLabel.Render("Started:"), info.StartTime.Format("2006-01-02 15:04:05"))
	}

	fmt.Println()
	fmt.Println(S.DetailSection.Render("LOCATION"))
	fmt.Printf("  %s  %s\n", S.DetailLabel.Render("Directory:"), info.CWD)
	fmt.Printf("  %s  %s\n", S.DetailLabel.Render("Project:"), S.Project.Render(info.ProjectName))
	if info.GitBranch != "" {
		fmt.Printf("  %s  %s\n", S.DetailLabel.Render("Branch:"), lipgloss.NewStyle().Foreground(S.StatusHealthy.GetForeground()).Render(info.GitBranch))
	}

	if len(tree) > 1 {
		fmt.Println()
		fmt.Println(S.DetailSection.Render("PROCESS TREE"))
		for i, node := range tree {
			if i == 0 {
				fmt.Printf("  %s %s\n", S.PID.Render(fmt.Sprintf("%d", node.PID)), node.Name)
			} else {
				fmt.Printf("  %s %s %s\n", S.Dim.Render("└─"), node.Name, S.PID.Render(fmt.Sprintf("(PID: %d)", node.PID)))
			}
		}
	}

	if info.Command != "" {
		fmt.Println()
		fmt.Println(S.DetailSection.Render("COMMAND"))
		fmt.Printf("  %s\n", S.Dim.Render(info.Command))
	}

	fmt.Println()
}

func PrintPsTable(processes []types.ProcessInfo) {
	if len(processes) == 0 {
		fmt.Println(S.Dim.Render("No dev processes found."))
		return
	}

	headers := []string{"PID", "PROCESS", "CPU%", "MEM", "PROJECT", "FRAMEWORK", "UPTIME", "WHAT"}
	colWidths := []int{8, 14, 6, 10, 18, 16, 10, 30}

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
			S.PID.Render(pidStr),
			S.Process.Render(p.ProcessName),
			cpuStr,
			p.Memory,
			S.Project.Render(p.ProjectName),
			colorizeFramework(frameworkStr),
			S.Uptime.Render(uptimeStr),
			S.Dim.Render(p.Description),
		}

		rows[i] = vals
	}

	fmt.Println(renderTable(headers, rows, colWidths))
	fmt.Println()
	fmt.Println(S.Dim.Render(fmt.Sprintf("%d process(es)", len(processes))))
}

func PrintCleanList(ports []types.PortInfo) {
	if len(ports) == 0 {
		fmt.Println(S.StatusHealthy.Render("No orphaned or zombie processes found."))
		return
	}

	fmt.Println(S.StatusOrphaned.Bold(true).Render("Found orphaned/zombie processes:"))
	fmt.Println()

	var lines []string
	for _, p := range ports {
		frameworkStr := p.Framework
		if frameworkStr == "" {
			frameworkStr = "-"
		}
		line := fmt.Sprintf("%s  PID: %s  %s  %s  %s",
			statusStyled(p.Status),
			S.PID.Render(fmt.Sprintf("%d", p.PID)),
			S.Process.Render(p.ProcessName),
			colorizeFramework(frameworkStr),
			S.Project.Render(p.ProjectName),
		)
		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	fmt.Println(S.CleanBorder.Render(content))
	fmt.Println()
}

func PrintWatchHeader() {
	fmt.Println(S.DetailTitle.Render("  bleat watch  "))
	fmt.Println()
	fmt.Println(S.Dim.Render("Watching for port changes... (Ctrl+C to exit)"))
	fmt.Println()
}

type WatchEvent struct {
	Kind string
	Info types.PortInfo
	Port int
}

func PrintWatchEvents(events []WatchEvent) {
	if len(events) == 0 {
		fmt.Println(S.Dim.Render("Watching for port changes..."))
		return
	}

	headers := []string{"EVENT", "PORT", "PROCESS", "PID", "PROJECT", "FRAMEWORK", "UPTIME"}
	colWidths := []int{8, 8, 14, 8, 18, 16, 10}

	rows := make([][]string, len(events))
	for i, e := range events {
		var eventStr, portStr, processStr, pidStr, projectStr, frameworkStr, uptimeStr string

		if e.Kind == "new" {
			eventStr = S.WatchNew.Render("▲ NEW")
			p := e.Info
			portStr = S.Port.Render(fmt.Sprintf(":%d", p.Port))
			processStr = S.Process.Render(p.ProcessName)
			pidStr = S.PID.Render(fmt.Sprintf("%d", p.PID))
			projectStr = S.Project.Render(p.ProjectName)
			fw := p.Framework
			if fw == "" {
				fw = "-"
			}
			frameworkStr = colorizeFramework(fw)
			uptimeStr = S.Uptime.Render(p.Uptime)

			plainVals := []string{"▲ NEW", fmt.Sprintf(":%d", p.Port), p.ProcessName, fmt.Sprintf("%d", p.PID), p.ProjectName, fw, p.Uptime}
			for j, v := range plainVals {
				if len(v) > colWidths[j] {
					colWidths[j] = len(v)
				}
			}
		} else {
			eventStr = S.WatchClosed.Render("▼ CLOSED")
			portStr = S.Port.Render(fmt.Sprintf(":%d", e.Port))
			processStr = S.Dim.Render("-")
			pidStr = S.Dim.Render("-")
			projectStr = S.Dim.Render("-")
			frameworkStr = S.Dim.Render("-")
			uptimeStr = S.Dim.Render("-")

			plainVals := []string{"▼ CLOSED", fmt.Sprintf(":%d", e.Port), "-", "-", "-", "-", "-"}
			for j, v := range plainVals {
				if len(v) > colWidths[j] {
					colWidths[j] = len(v)
				}
			}
		}

		rows[i] = []string{eventStr, portStr, processStr, pidStr, projectStr, frameworkStr, uptimeStr}
	}

	fmt.Println(renderTable(headers, rows, colWidths))
	fmt.Println()
	fmt.Println(S.Dim.Render(fmt.Sprintf("%d event(s)", len(events))))
}

func PrintNewPort(info types.PortInfo) {
	PrintWatchEvents([]WatchEvent{{Kind: "new", Info: info}})
}

func PrintClosedPort(port int) {
	PrintWatchEvents([]WatchEvent{{Kind: "closed", Port: port}})
}

func PrintKillPrompt() {
	fmt.Print(lipgloss.NewStyle().Bold(true).Render("Kill process? [y/N] "))
}

func PrintKillSuccess() {
	fmt.Println(S.StatusHealthy.Render("✓ Process killed successfully."))
}

func PrintKillFailed(err error) {
	fmt.Println(lipgloss.NewStyle().Foreground(S.StatusZombie.GetForeground()).Render("✗ Failed to kill process: ") + err.Error())
	fmt.Println(S.Dim.Render("Try running with sudo."))
}

func PrintCleanPrompt() {
	fmt.Print(lipgloss.NewStyle().Bold(true).Render("Kill all? [y/N] "))
}

func PrintCleanSummary(success, failed int) {
	fmt.Println()
	if failed == 0 {
		fmt.Println(S.StatusHealthy.Render(fmt.Sprintf("✓ Successfully killed %d process(es).", success)))
	} else {
		fmt.Println(S.StatusOrphaned.Render(fmt.Sprintf("Killed %d process(es), %d failed.", success, failed)))
	}
}

func PrintHelp() {
	fmt.Println()
	fmt.Println(S.HelpTitle.Render("  bleat  "))
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
		{"bleat theme", "Show current theme"},
		{"bleat theme list", "List available themes"},
		{"bleat theme set <name>", "Set the active theme"},
		{"bleat help", "Show this help message"},
	}

	cmdWidth := 0
	for _, c := range commands {
		if len(c[0]) > cmdWidth {
			cmdWidth = len(c[0])
		}
	}

	fmt.Println(S.DetailSection.Render("USAGE"))
	fmt.Println()
	for _, c := range commands {
		fmt.Printf("  %s  %s\n",
			S.HelpCmd.Render(padRight(c[0], cmdWidth)),
			S.HelpDesc.Render(c[1]))
	}
	fmt.Println()
}

func PrintThemeList(themes []string, currentTheme string) {
	fmt.Println()
	fmt.Println(S.HelpTitle.Render("  bleat themes  "))
	fmt.Println()
	fmt.Println(S.DetailLabel.Render("Current: ") + S.StatusHealthy.Render(currentTheme))
	fmt.Println()

	fmt.Println(S.DetailSection.Render("AVAILABLE THEMES"))
	fmt.Println()

	for _, name := range themes {
		if name == currentTheme {
			fmt.Printf("  %s %s\n", S.StatusHealthy.Render("●"), S.HelpCmd.Render(name))
		} else {
			gt, ok := gogh.Get(name)
			if ok {
				fmt.Printf("    %s %s%s%s%s%s%s%s%s\n",
					name,
					lipgloss.NewStyle().Foreground(gt.Red).Render("██"),
					lipgloss.NewStyle().Foreground(gt.Green).Render("██"),
					lipgloss.NewStyle().Foreground(gt.Yellow).Render("██"),
					lipgloss.NewStyle().Foreground(gt.Blue).Render("██"),
					lipgloss.NewStyle().Foreground(gt.Magenta).Render("██"),
					lipgloss.NewStyle().Foreground(gt.Cyan).Render("██"),
					lipgloss.NewStyle().Foreground(gt.BrightWhite).Render("██"),
					lipgloss.NewStyle().Foreground(gt.Foreground).Render("██"),
				)
			} else {
				fmt.Printf("    %s\n", S.Dim.Render(name))
			}
		}
	}
	fmt.Println()
}

func PrintThemeNotFound(name string) {
	fmt.Println()
	fmt.Println(S.StatusZombie.Render(fmt.Sprintf("Theme %q not found.", name)))
	fmt.Println(S.Dim.Render("Run 'bleat theme list' to see available themes."))
	fmt.Println()
}

func PrintThemeSetSuccess(name string) {
	fmt.Println()
	fmt.Println(S.StatusHealthy.Render(fmt.Sprintf("✓ Theme set to %q.", name)))
	fmt.Println(S.Dim.Render("Run any bleat command to see the new theme."))
	fmt.Println()
}

func PrintThemeShow(currentTheme string) {
	fmt.Println()
	fmt.Println(S.HelpTitle.Render("  bleat theme  "))
	fmt.Println()
	fmt.Printf("  %s  %s\n", S.DetailLabel.Render("Current:"), S.StatusHealthy.Render(currentTheme))
	fmt.Println()
	fmt.Println(S.Dim.Render("Run 'bleat theme list' to see available themes."))
	fmt.Println(S.Dim.Render("Run 'bleat theme set <name>' to change theme."))
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

func ListThemesSorted() []string {
	names := gogh.Names()
	sort.Strings(names)
	return names
}
