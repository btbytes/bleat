# bleat

A beautiful CLI tool to visualize and manage processes listening on network ports. Built for developers who run multiple dev servers, databases, and Docker containers locally.

```
  bleat                  # Show dev server ports (default)
  bleat --all            # Show all listening ports
  bleat <port>           # Detailed view of specific port
  bleat ps               # Show all running dev processes
  bleat clean            # Kill orphaned/zombie dev servers
  bleat watch            # Monitor port changes in real-time
```

![Bleat Logo](bleat.png)

## Features

- **Port Visualization** – See all listening ports in a clean, color-coded table
- **Framework Detection** – Automatically detects Next.js, Vite, React, PostgreSQL, Redis, Docker, and 30+ more frameworks
- **Process Management** – Kill processes directly from the detailed port view
- **Clean Orphaned Processes** – Find and kill zombie/orphaned dev servers in one command
- **Real-time Watch Mode** – Monitor port changes as they happen
- **Process List** – `bleat ps` shows all running dev processes sorted by CPU usage
- **Smart Filtering** – Excludes system apps (Spotify, Chrome, Slack, etc.) by default
- **Git Branch Detection** – Shows current git branch in detailed port view
- **Process Tree** – Visualizes the parent process chain for any port
- **Docker Collapse** – Combines Docker internal processes into a single row

## Installation

### Homebrew

```bash
brew install btbytes/brew/bleat
```

### From Source

```bash
git clone https://github.com/btbytes/bleat.git
cd bleat
go build -o bleat .
```

### Install Globally

```bash
go install github.com/btbytes/bleat@latest
```

Or copy the binary to your PATH:

```bash
go build -o bleat .
sudo mv bleat /usr/local/bin/
```

## Usage

### Show Dev Server Ports

```bash
bleat
```

Shows all listening ports for dev processes, excluding system apps.

```
  PORT     PROCESS        PID      PROJECT            FRAMEWORK        UPTIME     STATUS
  ────────────────────────────────────────────────────────────────────────────────────────────
  :3000    node           12345    my-app             Next.js          2h 15m     ● healthy
  :5432    postgres       6789     postgresql@16      PostgreSQL       3d 4h      ● healthy
  :6379    redis-server   1111     redis              Redis            3d 4h      ● healthy
  :8080    python3        2222     api-server         FastAPI          45m        ● healthy

4 port(s) · bleat --all to show all · bleat ps for processes
```

### Show All Listening Ports

```bash
bleat --all
# or shorthand
bleat -a
```

Includes system services and all listening processes.

### Detailed Port View

```bash
bleat 3000
```

Shows detailed information about a specific port:

```
  :3000

  Process
    Name:      node
    PID:       12345
    Status:    ● healthy
    Framework: Next.js
    Memory:    256 MB
    Uptime:    2h 15m
    Started:   2026-04-03 10:30:00

  Location
    Directory: /Users/dev/my-app
    Project:   my-app
    Branch:    feature/new-ui

  Process Tree
    /sbin/launchd (PID: 1)
    └─ node (PID: 12345)

  Command
    node .next/server/index.js

Kill process? [y/N]
```

Type `y` to kill the process, or press Enter to exit.

### Process List

```bash
bleat ps
```

Shows all running dev processes, sorted by CPU usage:

```
  PID      PROCESS        CPU%   MEM        PROJECT            FRAMEWORK        UPTIME     WHAT
  ───────────────────────────────────────────────────────────────────────────────────────────────
  12345    node           45.2   256 MB     my-app             Next.js          2h 15m     node .next/server...
  6789     postgres       2.1    128 MB     postgresql@16      PostgreSQL       3d 4h      postgres -D /opt/...
  2222     python3        1.5    64 MB      api-server         FastAPI          45m        uvicorn main:app...
```

Show all processes (including system):

```bash
bleat ps --all
```

### Clean Orphaned Processes

```bash
bleat clean
```

Finds orphaned (ppid=1) and zombie dev server processes, then prompts to kill them:

```
Found orphaned/zombie processes:

  ● orphaned  PID: 54321   node            Next.js          stale-app
  ● zombie    PID: 54322   python3         Django           old-api

Kill all? [y/N]
```

### Watch Mode

```bash
bleat watch
```

Monitors port changes in real-time, polling every 2 seconds:

```
  bleat watch

Watching for port changes... (Ctrl+C to exit)

  ▲ NEW  :3000    node            my-app             Next.js
  ▲ NEW  :5432   postgres        postgresql@16      PostgreSQL
  ▼ CLOSED  :8080
```

Press `Ctrl+C` to exit cleanly.

### Help

```bash
bleat help
```

Shows all available commands and their descriptions.

## Framework Detection

bleat automatically detects frameworks using multiple strategies:

### From package.json (Node.js projects)

| Framework | Dependency Key |
|-----------|---------------|
| Next.js | `next` |
| Nuxt | `nuxt`, `nuxt3` |
| SvelteKit | `@sveltejs/kit` |
| Svelte | `svelte` |
| Remix | `@remix-run/react`, `remix` |
| Astro | `astro` |
| Vite | `vite` |
| Angular | `@angular/core` |
| Vue | `vue` |
| React | `react` |
| Express | `express` |
| Fastify | `fastify` |
| Hono | `hono` |
| Koa | `koa` |
| NestJS | `@nestjs/core` |
| Gatsby | `gatsby` |
| Webpack | `webpack-dev-server` |
| esbuild | `esbuild` |
| Parcel | `parcel` |

### From file markers

| File | Framework |
|------|-----------|
| `Cargo.toml` | Rust |
| `go.mod` | Go |
| `manage.py` | Django |
| `Gemfile` | Ruby |
| `pyproject.toml` | Python |
| `requirements.txt` | Python |

### From process name

| Process | Framework |
|---------|-----------|
| `node`, `nodejs` | Node.js |
| `python`, `python3` | Python |
| `ruby` | Ruby |
| `java` | Java |
| `go` | Go |
| `bun` | Bun |
| `deno` | Deno |
| `cargo` | Rust |
| `mix`, `elixir` | Elixir |
| `php` | PHP |
| `dotnet` | .NET |

### From Docker images

| Image Pattern | Framework |
|--------------|-----------|
| `postgres*` | PostgreSQL |
| `redis` | Redis |
| `mysql`, `mariadb` | MySQL |
| `mongo` | MongoDB |
| `nginx` | nginx |
| `localstack` | LocalStack |
| `rabbitmq` | RabbitMQ |
| `kafka` | Kafka |
| `elasticsearch`, `opensearch` | Elasticsearch |
| `minio` | MinIO |

## Status Indicators

| Status | Indicator | Meaning |
|--------|-----------|---------|
| Healthy | `●` green | Normal running process |
| Orphaned | `●` yellow | Parent process exited (ppid=1) |
| Zombie | `●` red | Process is in zombie state |

## Technical Details

### Shell Commands Used

bleat uses these macOS/Unix commands under the hood:

- `lsof -iTCP -sTCP:LISTEN -P -n` – Find listening ports
- `ps -p <pids> -o pid=,ppid=,stat=,rss=,lstart=,command=` – Batch process info
- `lsof -a -d cwd -p <pids>` – Batch working directory lookup
- `docker ps --format "{{.Ports}}\t{{.Names}}\t{{.Image}}"` – Docker container info
- `git -C <dir> rev-parse --abbrev-ref HEAD` – Git branch detection
- `ps -eo pid=,ppid=,comm=` – Process tree construction

### Performance

- All PIDs are batched into single shell calls instead of N individual calls
- Git branch lookups are lazy-loaded (only in detailed port view)
- Framework detection caches results per scan

## Project Structure

```
bleat/
├── main.go              # Entry point, argument parsing, command routing
├── cmd/
│   ├── ports.go         # Main port listing command
│   ├── ps.go            # Process list command
│   ├── clean.go         # Clean orphaned processes
│   └── watch.go         # Watch mode
├── scanner/
│   ├── scanner.go       # Core scanning logic (lsof parsing, batch ps, cwd)
│   ├── detection.go     # Framework detection
│   ├── filtering.go     # Process filtering (system apps, dev processes)
│   └── process.go       # Process tree utilities
├── display/
│   └── table.go         # Table rendering, colors, all output formatting
├── types/
│   └── types.go         # Data structures (PortInfo, ProcessInfo, etc.)
├── go.mod
└── README.md
```

## Platform Support

| Platform | Support |
|----------|---------|
| macOS | Full support (primary target) |
| Linux | Should work with minimal changes |
| Windows | Not supported |

## Dependencies

- Standard library: `os/exec`, `os/user`, `path/filepath`, `io/ioutil`, `time`, `strings`, `strconv`, `regexp`, `sync`, `bufio`, `encoding/json`, `os`, `sort`
- `github.com/fatih/color` – Colored terminal output

## License

MIT
