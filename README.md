# agent-project-manager

`agent-project-manager` is a small Go service intended to run on a Raspberry Pi and orchestrate agent-based LLM workflows.
It ships as two binaries:

- **`agentd`**: daemon (HTTP API + orchestrator runtime + workers)
- **`agentctl`**: CLI client (runs locally or on the Pi)

---

## Repository layout

```text
llm-orchestrator/
├─ cmd/
│  ├─ agentd/                 # Raspberry Pi daemon: HTTP API + orchestrator runtime
│  │  └─ main.go
│  └─ agentctl/               # CLI client (runs locally or on the Pi)
│     └─ main.go
│
├─ internal/                  # Private application packages (not importable by other modules)
│  ├─ api/                    # HTTP routing, handlers, validation, auth middleware
│  │  ├─ router.go
│  │  └─ handlers_jobs.go
│  │
│  ├─ orchestrator/           # Workflow engine: state machine, retries, dispatch, wiring
│  │  ├─ orchestrator.go
│  │  └─ workflows.go
│  │
│  ├─ queue/                  # In-process job queue + worker pool
│  │  └─ queue.go
│  │
│  ├─ state/                  # SQLite persistence layer + models (DB access, queries)
│  │  ├─ store.go
│  │  └─ models.go
│  │
│  ├─ artifact/               # Workdir/repo management + artifact storage (files on disk)
│  │  └─ store.go
│  │
│  ├─ llm/                    # LLM provider interface + adapters (OpenAI, Ollama, etc.)
│  │  ├─ client.go
│  │  ├─ openai_client.go
│  │  └─ ollama_client.go
│  │
│  ├─ agents/                 # Agent implementations (architect / codegen / review)
│  │  ├─ architect.go
│  │  ├─ codegen.go
│  │  └─ review.go
│  │
│  ├─ tools/                  # Wrappers around external tools (git, go test, linters, etc.)
│  │  ├─ exec.go
│  │  └─ git.go
│  │
│  ├─ obs/                    # Observability: logging/metrics/tracing setup
│  │  └─ obs.go
│  │
│  └─ config/                 # Config loading + merge (defaults, overrides, env)
│     └─ config.go
│
├─ migrations/                # goose SQL migrations for SQLite
│  └─ 0001_init.sql
│
├─ configs/
│  ├─ config.yaml             # Default configuration (committed)
│  └─ config.local.yaml       # Optional local override (gitignored)
│
├─ scripts/
│  └─ dev.sh
│
├─ Makefile
├─ go.mod
└─ go.sum
```

---

## What runs where

- **`cmd/agentd`**: main runtime on the Raspberry Pi (HTTP API + orchestrator + worker execution)
- **`cmd/agentctl`**: CLI used to submit jobs, inspect status, fetch artifacts, etc.

---

## Key packages

| Package | Purpose |
|--------|---------|
| `internal/api` | HTTP routes/handlers, request validation, auth |
| `internal/orchestrator` | Workflow/state machine, retries, dispatch to agents/tools |
| `internal/queue` | In-process job queue + worker pool |
| `internal/state` | SQLite persistence + models |
| `internal/artifact` | Workspace + artifact storage on disk |
| `internal/llm` | Provider interface + adapters (OpenAI/Ollama) |
| `internal/agents` | Architect/codegen/review agent implementations |
| `internal/tools` | Wrappers around external tools (git, go test, golangci-lint, etc.) |
| `internal/obs` | Logging/metrics/tracing |
| `internal/config` | Config loading/merge (defaults + overrides) |

---

## Request flow (high level)

```text
agentctl
  │
  ▼
agentd (HTTP API)
  │
  ▼
orchestrator (workflow/state machine)
  │
  ▼
queue (worker pool)
  │
  ├─► agents (architect/codegen/review)
  └─► tools (git / go test / linters / etc.)
  │
  ▼
state (SQLite) + artifacts (files on disk)
```

---

## Configuration

Configuration defaults live in `configs/config.yaml`.
You can override locally with `configs/config.local.yaml` (ignored by git).

---

## Build

```bash
go build ./cmd/agentd
go build ./cmd/agentctl
```

---

## Run (examples)

```bash
# Run the daemon (example)
./agentd --config configs/config.yaml

# Use the CLI (example)
./agentctl --help
```
