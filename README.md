![Package Mate Terminal](https://i.postimg.cc/7hyBkmmQ/Screenshot-2026-03-29-at-12-48-05.png)

# Package Mate

Package Mate is a command line tool for macOS that manages the installation and removal of developer tools. It integrates with Homebrew, npm, pipx, and several other package managers to provide a unified interface for setting up a development environment.

The tool performs a system scan to detect existing installations, distinguishes between managed (Homebrew) and unmanaged binaries, handles updates, and safely resolves conflicts. It also protects core system paths from accidental modification.

## Features

- Interactive search dashboard with real-time filtering and keyboard navigation.
- Installation, update, and uninstallation of tools via Homebrew formulae and casks.
- Support for special installers: Rust (rustup), NVM, npm global packages, pipx, and direct binary checks.
- Detection of multiple installation sources (managed, unmanaged, outdated, or different Homebrew versions).
- Conflict resolution with prompts for overriding unmanaged binaries or removing dependency-linked formulae.
- Protection against modifying macOS system protected paths (SIP).
- Detailed version information and installation date retrieval.
- Background installation system with desktop notifications.
- Interactive search dashboard: type to filter, press Enter to open tool actions.

## Installation

Install Package Mate with a single command:

```bash
curl -fsSL https://package-mate.com/install.sh | bash
```

### First Run

When you run `mate` for the first time, if Homebrew is not present, the tool will display a prompt to install it. Homebrew is required for most installation operations.

## Usage

### Dashboard

```bash
mate
```

Opens an interactive terminal UI showing all available tools grouped by category. Type to search and filter results in real-time. Status indicators show:

- `[✓]` Installed and up to date (managed)
- `[↻]` Installed but an update is available
- `[⚙]` Installed but not managed by Package Mate (unmanaged)
- `[?]` Multiple installations detected

Use the arrow keys or `j`/`k` to navigate results. Press `Enter` to select a tool and open the action menu.

### Direct Tool Access

Select a tool from the dashboard to open its action menu with three options:

1. **Install or Update** - Installs the tool or updates if already installed
2. **Uninstall** - Removes the tool from your system
3. **Information & Versions** - Shows detailed version and installation info

Use the search dashboard to find and select tools. Type to filter results in real-time, then press `Enter` to open the action menu for any tool.

### Information

![Package Mate Info](https://i.postimg.cc/dQy1N4Mc/Screenshot-2026-03-29-at-12-58-29.png)

### Examples:

```bash
# Open the dashboard
mate

# Type to search for Node.js, then press Enter to open its action menu

# Select option 1 to install or update

# Search for PostgreSQL, press Enter, then select option 3 for version info

# Search for redis, press Enter, select uninstall to remove an unmanaged binary

# View background installations
mate bg

# Clean up Homebrew cache (free disk space)
mate cleanup

# Clear finished background jobs
mate cleanup bg
```

## Catalog Overview

Package Mate includes a curated catalog of **over 400 developer tools** organized into sections:

- **Homebrew Setup**: Homebrew itself and update command.
- **Backend & Runtime**: Node, Bun, Deno, Rust, Python, Java, Go, PHP, Ruby, .NET SDK, Elixir, Erlang, Delve, Air, golang-migrate.
- **Languages & Runtimes**: Zig, Haskell, OCaml, C++ (Clang), Rustup, Lua, Perl, Dart, Kotlin, Scala, Fortran (GCC), Julia, Flutter.
- **Coding CLIs & AI**: Claude Code, Gemini CLI, Aider, Continue CLI, GitHub CLI, Ollama, Swag, Mockery, Templ, Antigravity.
- **Dev Essentials**: Git, Git LFS, LazyGit, fzf, ripgrep, bat, eza, zoxide, Starship, direnv, NVM.
- **Package Managers**: npm, Yarn, pnpm, uv, Poetry, pipx, Composer, Corepack.
- **Testing & Utilities**: Playwright, Cypress, Vitest, pre-commit, tox, Ruff, Black, ESLint, Hurl, Vegeta, Prism.
- **Editors & IDEs**: Cursor, Zed, Sublime Text, Android Studio, Nova.
- **Performance & Profiling**: gperftools, Valgrind, Hyperfine, bpftrace, k6, wrk, Locust.
- **Dev Utilities**: Watchman, Mise, Lefthook, Act, JSON Export (jtbl), golangci-lint, GoReleaser, Biome, pyenv, mkcert, Task.
- **Databases**: MySQL, PostgreSQL, MariaDB, CockroachDB, TimescaleDB, MongoDB, SQLite, ClickHouse, ElasticSearch, Neo4j, Firebase CLI, Supabase CLI, pgcli, sqlc, Atlas, dbmate, Litestream, pgBadger.
- **Caching & Messaging**: Redis, Memcached, NATS, Kafka, RabbitMQ, ActiveMQ, ZeroMQ, kcat.
- **Containers & DevOps**: Docker, Docker Compose, Colima, Kubernetes CLI, k9s, Helm, Terraform, Ansible, AWS CLI, Google Cloud SDK, Azure CLI, Lazydocker.
- **DB GUI & Dev Tools**: TablePlus, DBeaver, RedisInsight, Postman, Insomnia, Lens, Sequel Ace, pgweb, usql.
- **Infrastructure & Cloud**: Vercel CLI, Netlify CLI, Heroku CLI, Pulumi, OpenTofu, AWS CDK, DigitalOcean CLI, Prometheus, Grafana, Vault, Boundary.
- **Virtualization**: UTM, VMware Fusion, VirtualBox, Vagrant, vagrant-completion.
- **Cloud & Kubernetes**: Krew, Kustomize, Minikube, LocalStack CLI, Tfenv, Terragrunt, Kwid.
- **System Tools**: htop, btop, NeoVim, Vim, jq, yq, fd, ncdu, Tree, entr, Lnav, Ouch.
- **Terminal Glow-Up & Shell**: Tmux, Fastfetch, Glow, Gum, Oh My Zsh, Zsh Autosuggestions, Zsh Syntax Highlighting, Aria2, VHS, Mods, Atuin.
- **macOS Essentials**: Raycast, Rectangle, AltTab, Latest, HiddenBar, Shottr, AppCleaner, IINA, Xcodes CLI, Syncthing, Maccy, CleanShot X.
- **Modern CLI Replacements**: Dust, Bottom, Tldr, TheFuck, Rclone, Dog, Choose, Gdu, Zellij, Just, OrbStack, fnm, scc, Watchexec, Dasel.
- **Terminal Emulators**: Warp, WezTerm, Alacritty, Kitty, Ghostty.
- **Media & Graphics**: FFmpeg, ImageMagick, ExifTool, OptiPNG, Inkscape, GIMP, Blender, Spotify.
- **Low-Level & Embedded**: avr-gcc, arm-none-eabi-gcc, SDCC, avrdude, OpenOCD, dfu-util, minicom, screen, picocom, KiCad, Fritzing.
- **Docs & Static Sites**: Hugo, Jekyll, Eleventy, MkDocs, Pandoc, Doxygen, Graphviz, Spectral, OpenAPI Generator.
- **Data & Analytics**: DuckDB, VisiData, DVC, PostGIS.
- **Web Browsers**: Google Chrome, Firefox Developer Edition, Arc, Brave.
- **Knowledge & Productivity**: Notion, Obsidian, Linear, Todoist.
- **Security & Secrets**: Age, Sops, Cosign, Bitwarden CLI, 1Password CLI, KeepassXC, Wireshark, Nmap, Burp Suite, Snyk CLI, Trivy, Gitleaks, Ngrok, Tailscale.
- **Reverse Engineering**: Ghidra, Radare2, Hopper Disassembler, Charles, mitmproxy, Proxyman, sqlmap, John the Ripper, Hashcat, Sleuth Kit, Binwalk.
- **Networking & API**: Doggo, MTR, Step CLI, HTTPie, Cloudflared, Termius, Stripe CLI, Websocat, Bruno, xh, gRPCurl, Mole, Trippy.
- **Communications**: Discord, Slack, Telegram, Zoom.
- **Security Network**: Additional security and networking tools.

Each tool defines its primary binary name, description, and an optional Homebrew formula or cask. Special tools use custom installation logic (for example, Rust via rustup, Claude Code via npm).

## How It Works

### Detection System

When you run `mate`, the tool performs a single-pass system scan:

1. Scans `PATH` for binaries and resolves symlinks.
2. Queries Homebrew for installed formulae and casks, including versioned formulae (e.g., `postgresql@16`).
3. Checks for outdated packages using `brew outdated`.
4. Detects unmanaged `.app` bundles in `/Applications` and `~/Applications`.

The detection results drive the dashboard status icons and inform the interactive prompts.

### Installation Flow

For a given tool, the installation process follows this logic:

1. Check if the tool is already installed using the detection system.
2. If installed and up to date, exit with a notice.
3. If installed but outdated, prompt the user to update.
4. If an unmanaged binary exists at a different path, offer to override or install alongside.
5. If a different Homebrew version exists, offer to update or override.
6. For casks, verify the application is not running before proceeding.
7. Execute the appropriate install command (`brew install`, `brew install --cask`, `npm install -g`, `pipx install`, etc.).
8. If an override was performed, remove the old binary or Homebrew package after the new installation completes, handling permission errors and dependency conflicts with `sudo` or force removal as needed.

### Background Installation System

Installations run in the background as detached processes:

- Jobs are persisted to `~/.package-mate/jobs/` as JSON files
- Desktop notifications when installations complete or fail
- Real-time status tracking with `mate bg`
- Automatic cleanup of jobs older than 24 hours
- Manual cleanup with `mate cleanup bg`

### Homebrew Cleanup

Free up disk space by removing old Homebrew formula versions and cached files:

```bash
mate cleanup
```

Choose between:

- **Standard Cleanup** - Removes old versions, keeps current (safe)
- **Deep Cleanup** - Removes ALL cached files including unused dependencies

Shows disk space freed after completion.

### Safety Protections

Package Mate refuses to modify files in system protected directories:

- `/System`
- `/Library/Apple`
- `/usr/bin`, `/usr/sbin`, `/bin`, `/sbin`
- `/var/root`
- `/opt/apple`

If an unmanaged binary resides in one of these paths, the tool displays a safety alert and aborts the operation. This prevents accidental damage to macOS core components.

### Uninstall Handling

Uninstallation supports removing specific versions when multiple are present. For managed versions, it uses `brew uninstall`. For unmanaged binaries, it uses `rm` with `sudo` when necessary, after confirming the path is not protected. For dependency conflicts (e.g., a formula required by others), the tool prompts the user to force uninstall or skip.

## Building from Source

```bash
git clone https://github.com/yousefbustamiii/package-mate.git
cd package-mate
go mod download
go build -o mate .
```

Run tests (none currently, but the codebase is structured for future addition):

```bash
go test ./...
```

## Project Structure

- `cmd/` – Command implementations for install, uninstall, info, and background jobs.
- `cmd/bg/` – Interactive background jobs viewer.
- `internal/components/` – Static catalog of tools, sections, and resolution logic.
- `internal/installer/` – Core installation logic, Homebrew interaction, detection, and version management.
- `internal/installer/specials/` – Custom installers for non-Homebrew tools (npm, pipx, rustup, nvm).
- `internal/installer/versions/` – Version detection for special tools.
- `internal/background/` – Background job management, persistence, and notifications.
- `internal/sys/` – System helpers (architecture, shell detection, protected paths, app-running checks).
- `internal/ui/` – Terminal user interface components (colored output, spinners, interactive search, prompts).

![Package Mate Search](https://i.postimg.cc/ZRpM07KX/Screenshot-2026-03-31-at-21-36-24.png)

## License

This project is licensed under the Apache License 2.0. See the `LICENSE` file for details.

## Contributing

Contributions are welcome. Please open an issue or pull request on GitHub.
