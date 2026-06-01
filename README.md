# yank-run-cli

Command-line client for [yank.run](https://yank.run) — a content-addressed
snippet store. See [the backend API reference](../yank-run-backend/docs/BACKEND_API.md)
and the [CLI design plan](../yank-run-backend/docs/CLI_PLAN.md) for the full spec.

## Install

From source (Go 1.26+):

```bash
git clone https://github.com/korotinm/yank-run-cli
cd yank-run-cli
make build              # → ./bin/yank
# or
make install            # → $GOBIN/yank
```

Or directly:

```bash
go install github.com/korotinm/yank-run-cli/cmd/yank@latest
```

## Quick start

```bash
# Push the last shell command as a snippet
history | tail -1 | yank push --tag bash

# Search and run
yank search vault kubectl
yank cat 0 | bash               # 0 = first result of last search

# Inspect first, then run
yank show 2
yank cat 2 | bash

# Copy to clipboard
yank copy <id>

# Open in browser (uses YANK_WEB_URL if set)
yank open <id>
```

## Commands

| Command | Purpose |
|---|---|
| `yank push`   | create a snippet from `-b BODY`, `-f FILE`, or piped stdin |
| `yank cat`    | print snippet body to stdout (pipe-friendly) |
| `yank show`   | colored, human-readable snippet view |
| `yank search` | full-text search; results cached for index references |
| `yank copy`   | put snippet body on the system clipboard |
| `yank open`   | open the snippet URL in the default browser |

`REF` arguments in `cat`/`show`/`copy`/`open` accept either a 64-char
hex id or a 0-based index into the last `yank search` result list.

## Configuration

Environment variables:

| Var | Default | Effect |
|---|---|---|
| `YANK_URL`     | `https://api.yank.run` | backend API base URL (set to `http://localhost:8080` for local backend dev) |
| `YANK_WEB_URL` | unset | frontend URL used by `yank open` |
| `NO_COLOR`     | unset | disable ANSI colors when set |

CLI flags (`--url`, `--timeout`, `--json`, `--no-color`) always override env.

## Exit codes

| Code | Meaning |
|---|---|
| 0 | success |
| 2 | usage / client-side validation error |
| 3 | server 4xx (not 429) or bad REF |
| 4 | rate limited (HTTP 429) |
| 5 | server 5xx, network failure, clipboard error |

## Build

```bash
make build              # local bin/yank
make cross              # dist/ for darwin, linux, windows × amd64/arm64
make vet
make fmt
```

Version is embedded via `-ldflags`; `yank --version` prints it.
