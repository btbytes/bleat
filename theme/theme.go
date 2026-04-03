package theme

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
	gogh "github.com/willyv3/gogh-themes/lipgloss"
)

type Theme struct {
	Name      string
	goghTheme gogh.Theme
	Styles    Styles
}

type Styles struct {
	Header          lipgloss.Style
	TableBorder     lipgloss.Style
	Port            lipgloss.Style
	PID             lipgloss.Style
	Process         lipgloss.Style
	Project         lipgloss.Style
	Uptime          lipgloss.Style
	Dim             lipgloss.Style
	StatusHealthy   lipgloss.Style
	StatusOrphaned  lipgloss.Style
	StatusZombie    lipgloss.Style
	DetailTitle     lipgloss.Style
	DetailSection   lipgloss.Style
	DetailLabel     lipgloss.Style
	CleanBorder     lipgloss.Style
	WatchNew        lipgloss.Style
	WatchClosed     lipgloss.Style
	HelpTitle       lipgloss.Style
	HelpCmd         lipgloss.Style
	HelpDesc        lipgloss.Style
	FrameworkColors map[string]lipgloss.Color
}

var defaultThemeName = "Console"

func BuildStyles(name string) (Styles, error) {
	if name == "" || name == "default" {
		name = defaultThemeName
	}

	gt, ok := gogh.Get(name)
	if !ok {
		return Styles{}, fmt.Errorf("theme %q not found", name)
	}

	frameworkColors := map[string]lipgloss.Color{
		"Next.js":       gt.BrightCyan,
		"Vite":          gt.Yellow,
		"React":         gt.BrightBlue,
		"Vue":           gt.Green,
		"Angular":       gt.Red,
		"Svelte":        gt.Yellow,
		"SvelteKit":     gt.Yellow,
		"Express":       gt.BrightBlack,
		"Fastify":       gt.BrightCyan,
		"NestJS":        gt.Red,
		"Nuxt":          gt.Green,
		"Remix":         gt.Blue,
		"Astro":         gt.Magenta,
		"Django":        gt.Green,
		"Flask":         gt.BrightCyan,
		"FastAPI":       gt.BrightBlue,
		"Rails":         gt.Red,
		"Gatsby":        gt.Magenta,
		"Go":            gt.BrightBlue,
		"Rust":          gt.Yellow,
		"Ruby":          gt.Red,
		"Python":        gt.Yellow,
		"Node.js":       gt.Green,
		"Java":          gt.Red,
		"Docker":        gt.Blue,
		"PostgreSQL":    gt.Blue,
		"Redis":         gt.Red,
		"MySQL":         gt.Blue,
		"MongoDB":       gt.Green,
		"nginx":         gt.Green,
		"LocalStack":    gt.BrightCyan,
		"RabbitMQ":      gt.Yellow,
		"Kafka":         gt.BrightCyan,
		"Elasticsearch": gt.Yellow,
		"MinIO":         gt.Red,
		"Bun":           gt.Yellow,
		"Deno":          gt.BrightCyan,
		"PHP":           gt.Magenta,
		"Elixir":        gt.Magenta,
		".NET":          gt.Blue,
	}

	return Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(gt.Foreground).
			Background(gt.Blue).
			Padding(0, 1),

		TableBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gt.Magenta).
			Padding(0, 1),

		Port: lipgloss.NewStyle().
			Foreground(gt.Cyan).
			Bold(true),

		PID: lipgloss.NewStyle().
			Foreground(gt.Magenta),

		Process: lipgloss.NewStyle().
			Foreground(gt.BrightWhite),

		Project: lipgloss.NewStyle().
			Foreground(gt.BrightYellow),

		Uptime: lipgloss.NewStyle().
			Foreground(gt.BrightBlack),

		Dim: lipgloss.NewStyle().
			Foreground(gt.BrightBlack),

		StatusHealthy: lipgloss.NewStyle().
			Foreground(gt.Green),

		StatusOrphaned: lipgloss.NewStyle().
			Foreground(gt.Yellow),

		StatusZombie: lipgloss.NewStyle().
			Foreground(gt.Red),

		DetailTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(gt.Foreground).
			Background(gt.Blue).
			Padding(0, 1),

		DetailSection: lipgloss.NewStyle().
			Bold(true).
			Foreground(gt.Magenta).
			MarginTop(1).
			MarginBottom(0),

		DetailLabel: lipgloss.NewStyle().
			Foreground(gt.BrightBlack).
			Width(12).
			Align(lipgloss.Right),

		CleanBorder: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gt.Yellow).
			Padding(0, 1),

		WatchNew: lipgloss.NewStyle().
			Foreground(gt.Green).
			Bold(true),

		WatchClosed: lipgloss.NewStyle().
			Foreground(gt.Red).
			Bold(true),

		HelpTitle: lipgloss.NewStyle().
			Bold(true).
			Foreground(gt.Foreground).
			Background(gt.Blue).
			Padding(0, 1),

		HelpCmd: lipgloss.NewStyle().
			Foreground(gt.Cyan),

		HelpDesc: lipgloss.NewStyle().
			Foreground(gt.BrightBlack),

		FrameworkColors: frameworkColors,
	}, nil
}

func ListThemes() []string {
	names := gogh.Names()
	sort.Strings(names)
	return names
}

func HasTheme(name string) bool {
	_, ok := gogh.Get(name)
	return ok
}

func DefaultThemeName() string {
	return defaultThemeName
}
