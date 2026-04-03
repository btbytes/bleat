package scanner

import (
	"strings"

	"github.com/btbytes/bleat/types"
)

var systemApps = map[string]bool{
	"spotify":       true,
	"raycast":       true,
	"tableplus":     true,
	"postman":       true,
	"linear":        true,
	"cursor":        true,
	"controlce":     true,
	"rapportd":      true,
	"superhuma":     true,
	"setappage":     true,
	"slack":         true,
	"discord":       true,
	"firefox":       true,
	"chrome":        true,
	"google":        true,
	"safari":        true,
	"figma":         true,
	"notion":        true,
	"zoom":          true,
	"teams":         true,
	"code":          true,
	"iterm2":        true,
	"warp":          true,
	"arc":           true,
	"loginwindow":   true,
	"windowserver":  true,
	"systemuise":    true,
	"kernel_task":   true,
	"launchd":       true,
	"mdworker":      true,
	"mds_stores":    true,
	"cfprefsd":      true,
	"coreaudio":     true,
	"corebrightne":  true,
	"airportd":      true,
	"bluetoothd":    true,
	"sharingd":      true,
	"usernoted":     true,
	"notificationc": true,
	"cloudd":        true,
}

var devProcessNames = map[string]bool{
	"node":       true,
	"python":     true,
	"python3":    true,
	"ruby":       true,
	"java":       true,
	"go":         true,
	"cargo":      true,
	"deno":       true,
	"bun":        true,
	"php":        true,
	"uvicorn":    true,
	"gunicorn":   true,
	"flask":      true,
	"rails":      true,
	"npm":        true,
	"npx":        true,
	"yarn":       true,
	"pnpm":       true,
	"tsc":        true,
	"tsx":        true,
	"esbuild":    true,
	"rollup":     true,
	"turbo":      true,
	"nx":         true,
	"jest":       true,
	"vitest":     true,
	"mocha":      true,
	"pytest":     true,
	"cypress":    true,
	"playwright": true,
	"rustc":      true,
	"dotnet":     true,
	"gradle":     true,
	"mvn":        true,
	"mix":        true,
	"elixir":     true,
}

func IsSystemApp(name string) bool {
	lower := strings.ToLower(name)
	return systemApps[lower]
}

func isDevProcess(name string) bool {
	lower := strings.ToLower(name)
	return devProcessNames[lower]
}

func IsDockerProcess(name string) bool {
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "com.docke") {
		return true
	}
	return lower == "docker" || lower == "docker-sandbox"
}

func FilterDevPorts(ports []types.PortInfo) []types.PortInfo {
	var filtered []types.PortInfo
	for _, p := range ports {
		if !IsSystemApp(p.RawName) && !IsDockerProcess(p.RawName) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
