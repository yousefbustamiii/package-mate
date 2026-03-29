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
- Single command access: `mate` opens the dashboard, `mate <tool>` opens an interactive action menu.

## Installation

### Prerequisites

Package Mate requires Go 1.22 or later to build from source. The binary itself has no runtime dependencies other than the tools it manages.

### From Source

```bash
git clone https://github.com/yousefbustamiii/package-mate.git
cd package-mate
go build -o mate .
sudo cp mate /usr/local/bin/   # or any directory in your PATH
```

Alternatively, you can install using Go:

```bash
go install github.com/yousefbustamiii/package-mate@latest
```

### First Run
When you run `mate` for the first time, if Homebrew is not present, the tool will display a prompt to install it. Homebrew is required for most installation operations.

## Usage

### Dashboard

```bash
mate
```

Opens an interactive terminal UI showing all available tools grouped by category. Status indicators show:

- `[✓]` Installed and up to date (managed)
- `[↻]` Installed but an update is available
- `[⚙]` Installed but not managed by Package Mate (unmanaged)
- `[?]` Multiple installations detected

Use the arrow keys or `j`/`k` to navigate results. Press `Enter` to select a tool and open the action menu.

### Direct Tool Access

Opens an action menu for the specified tool directly, bypassing the search dashboard. For example:

```bash
mate node
mate redis
mate docker
```

The action menu provides three options:
1. Install or Update
2. Uninstall
3. Information & Versions

### Information

![Package Mate Info](https://i.postimg.cc/dQy1N4Mc/Screenshot-2026-03-29-at-12-58-29.png)

### Examples:

```bash
# Open the dashboard
mate

# Install Node.js
mate node
# then select option 1

# Check versions of PostgreSQL
mate postgresql
# then select option 3

# Remove an unmanaged binary
mate redis
# select uninstall, then choose the unmanaged version from the list
```

## Catalog Overview
Package Mate includes a curated catalog of over 80 developer tools organized into sections:

- **Homebrew Setup**: Homebrew itself and update command.
- **Databases**: MySQL, PostgreSQL, MariaDB, CockroachDB, TimescaleDB, MongoDB, SQLite, ClickHouse, ElasticSearch, Neo4j, Firebase CLI, Supabase CLI.
- **Caching & Messaging**: Redis, Memcached, NATS, Kafka, RabbitMQ, ActiveMQ, ZeroMQ, kcat.
- **Containers & DevOps**: Docker, Docker Compose, Colima, Kubernetes CLI, k9s, Helm, Terraform, Ansible, AWS CLI, Google Cloud SDK, Azure CLI.
- **Backend & Runtime**: Node, Bun, Deno, Rust, Python, Java, Go, PHP, Ruby, .NET SDK, Elixir, Erlang.
- **DB GUI & Dev Tools**: TablePlus, DBeaver, RedisInsight, Postman, Insomnia, HTTPie, Lens, Sequel Ace.
- **Coding CLIs & AI**: Claude Code, Gemini CLI, Aider, Continue CLI, GitHub CLI.
- **Dev Essentials**: Git, Git LFS, LazyGit, fzf, ripgrep, bat, eza, zoxide, Starship, direnv, NVM.
- **Package Managers**: npm, Yarn, pnpm, uv, Poetry, pipx, Composer, Corepack.
- **Testing & Utilities**: Playwright, Cypress, Vitest, pre-commit, tox, Ruff, Black, ESLint.
- **System Tools**: htop, btop, NeoVim, Vim, jq, yq, fd, ncdu, tree, entr.

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

- `cmd/` – Command implementations for install, uninstall, and info.
- `internal/components/` – Static catalog of tools, sections, and resolution logic.
- `internal/installer/` – Core installation logic, Homebrew interaction, detection, and version management.
- `internal/installer/specials/` – Custom installers for non-Homebrew tools (npm, pipx, rustup, nvm).
- `internal/installer/versions/` – Version detection for special tools.
- `internal/sys/` – System helpers (architecture, shell detection, protected paths, app-running checks).
- `internal/ui/` – Terminal user interface components (colored output, spinners, interactive search, prompts).

![Package Mate Search](https://i.postimg.cc/x1xn1zr9/Screenshot-2026-03-29-at-13-00-24.png)

## License

This project is licensed under the Apache License 2.0. See the `LICENSE` file for details.

## Contributing
Contributions are welcome. Please open an issue or pull request on GitHub.
