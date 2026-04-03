package scanner

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/btbytes/bleat/types"
)

func DetectFramework(info types.PortInfo) string {
	if info.CWD == "" {
		return detectFromProcess(info.ProcessName, info.Command)
	}

	if framework := detectFromPackageJSON(info.CWD); framework != "" {
		return framework
	}

	if framework := detectFromFileMarkers(info.CWD); framework != "" {
		return framework
	}

	if framework := detectFromDocker(info.Command); framework != "" {
		return framework
	}

	return detectFromProcess(info.ProcessName, info.Command)
}

func detectFromPackageJSON(cwd string) string {
	packagePath := filepath.Join(cwd, "package.json")
	data, err := os.ReadFile(packagePath)
	if err != nil {
		return ""
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return ""
	}

	allDeps := make(map[string]string)
	for k, v := range pkg.Dependencies {
		allDeps[k] = v
	}
	for k, v := range pkg.DevDependencies {
		allDeps[k] = v
	}

	frameworkOrder := []struct {
		name string
		keys []string
	}{
		{"Next.js", []string{"next"}},
		{"Nuxt", []string{"nuxt", "nuxt3"}},
		{"SvelteKit", []string{"@sveltejs/kit"}},
		{"Svelte", []string{"svelte"}},
		{"Remix", []string{"@remix-run/react", "remix"}},
		{"Astro", []string{"astro"}},
		{"Angular", []string{"@angular/core"}},
		{"NestJS", []string{"@nestjs/core", "nestjs"}},
		{"Gatsby", []string{"gatsby"}},
		{"Vite", []string{"vite"}},
		{"Vue", []string{"vue"}},
		{"React", []string{"react"}},
		{"Express", []string{"express"}},
		{"Fastify", []string{"fastify"}},
		{"Hono", []string{"hono"}},
		{"Koa", []string{"koa"}},
		{"Webpack", []string{"webpack-dev-server"}},
		{"esbuild", []string{"esbuild"}},
		{"Parcel", []string{"parcel"}},
	}

	for _, f := range frameworkOrder {
		for _, key := range f.keys {
			if _, ok := allDeps[key]; ok {
				return f.name
			}
		}
	}

	return ""
}

func detectFromFileMarkers(cwd string) string {
	markers := map[string]string{
		"Cargo.toml":       "Rust",
		"go.mod":           "Go",
		"manage.py":        "Django",
		"Gemfile":          "Ruby",
		"pyproject.toml":   "Python",
		"requirements.txt": "Python",
	}

	for marker, framework := range markers {
		if _, err := os.Stat(filepath.Join(cwd, marker)); err == nil {
			return framework
		}
	}

	return ""
}

func detectFromProcess(processName, command string) string {
	lower := strings.ToLower(processName)

	if strings.Contains(command, "uvicorn") || strings.Contains(command, "gunicorn") {
		if strings.Contains(command, "fastapi") {
			return "FastAPI"
		}
		return "Flask"
	}

	if strings.Contains(command, "rails") {
		return "Rails"
	}

	switch lower {
	case "node", "nodejs":
		return "Node.js"
	case "python", "python3":
		return "Python"
	case "ruby":
		return "Ruby"
	case "java":
		return "Java"
	case "go":
		return "Go"
	case "bun":
		return "Bun"
	case "deno":
		return "Deno"
	case "php":
		return "PHP"
	case "cargo":
		return "Rust"
	case "mix":
		return "Elixir"
	case "elixir":
		return "Elixir"
	case "dotnet":
		return ".NET"
	}

	return ""
}

func detectFromDocker(command string) string {
	lower := strings.ToLower(command)

	dockerImages := map[string]string{
		"postgres":      "PostgreSQL",
		"redis":         "Redis",
		"mysql":         "MySQL",
		"mariadb":       "MySQL",
		"mongo":         "MongoDB",
		"nginx":         "nginx",
		"localstack":    "LocalStack",
		"rabbitmq":      "RabbitMQ",
		"kafka":         "Kafka",
		"elasticsearch": "Elasticsearch",
		"opensearch":    "Elasticsearch",
		"minio":         "MinIO",
	}

	for image, framework := range dockerImages {
		if strings.Contains(lower, image) {
			return framework
		}
	}

	if strings.Contains(lower, "docker") {
		return "Docker"
	}

	return ""
}
