# OpenSpec: Agent Runtime Interface (ARI) v0.1

Status: draft

## Summary

ARI defines a runtime-agnostic interface for launching, resuming, messaging, and
monitoring autonomous agent sessions. It abstracts Claude Code specifics so
alternate runtimes (e.g., Codex) can be supported without changing higher-level
Gas Town workflows (mail, beads, convoys, molecules).

## Goals

- Decouple agent orchestration from any single CLI/runtime.
- Preserve current Claude behavior with a Claude adapter.
- Enable a Codex adapter with minimal changes to core logic.
- Keep minimal mode and full stack (tmux) mode parity.

## Non-Goals (v0.1)

- Unified transcript storage or cross-runtime replay.
- Multi-runtime scheduling policies or orchestration heuristics.
- UI changes.

## Terms

- Runtime: an agent CLI process (Claude Code, Codex).
- Session: a single runtime invocation tied to a worktree.
- Adapter: runtime-specific implementation of ARI.

## Interface

### Go-ish interface (conceptual)

```
// AgentRuntime is the runtime abstraction.
type AgentRuntime interface {
    Start(ctx context.Context, opts StartOptions) (SessionHandle, error)
    Resume(ctx context.Context, handle SessionHandle) error
    SendMessage(ctx context.Context, handle SessionHandle, msg Message) error
    Stop(ctx context.Context, handle SessionHandle, reason string) error
    IsReady(ctx context.Context, handle SessionHandle) (bool, error)
    DetectRunning(ctx context.Context, handle SessionHandle) (bool, error)
    ListSessions(ctx context.Context, filter SessionFilter) ([]SessionHandle, error)
}

// StartOptions defines a launch request.
type StartOptions struct {
    WorkDir       string
    RuntimeName   string
    AccountDir    string
    Env           map[string]string
    InitialPrompt string
    Mode          string // "minimal" | "tmux"
}

type SessionHandle struct {
    Runtime   string
    SessionID string
    WorkDir   string
    PID       int
    StartedAt time.Time
    ReadyAt   time.Time
}

// Message defines runtime-agnostic message delivery.
type Message struct {
    Text     string
    Delivery string // "stdin" | "tmux" | "rpc"
    Timeout  time.Duration
}
```

### Readiness policy

- Each adapter must implement a readiness strategy and a retry budget.
- If a runtime lacks a prompt marker, the adapter should:
  - wait a fixed warmup delay
  - attempt delivery
  - retry with backoff on failure

## Environment Contract

Gas Town defines canonical env vars. Adapters may map to native vars.

- `GT_SESSION_ID` (canonical)
- `GT_ROLE`, `GT_RIG`, `GT_POLECAT`, `BD_ACTOR`, `GT_WORKDIR`

Adapters may also set runtime-specific vars (e.g., `CLAUDE_SESSION_ID`) but core
logic should read the canonical names only.

## Hook Events

Canonical events and payload schema used across runtimes.

### Events

- `SessionStart`
- `SessionStop`
- `OnMessage`
- `OnError`

### Payload schema (JSON)

```
{
  "event": "SessionStart",
  "runtime": "claude",
  "session_id": "gt-abc123",
  "workdir": "/path/to/worktree",
  "rig": "greenplace",
  "role": "polecat",
  "bead": "gp-123",
  "actor": "greenplace/polecats/ivy",
  "timestamp": "2025-01-01T12:34:56Z",
  "data": {
    "message": "optional runtime-specific fields"
  }
}
```

Adapters map this schema to runtime hook mechanisms.

## Configuration

### Runtime registry

`~/.gastown/runtimes.json` (example):

```
{
  "default": "claude",
  "runtimes": {
    "claude": {
      "bin": "claude",
      "ready": "prompt",
      "delivery": "tmux"
    },
    "codex": {
      "bin": "codex",
      "ready": "warmup",
      "delivery": "stdin"
    }
  }
}
```

### Per-rig override

`<rig>/.gastown/runtime.json`:

```
{ "runtime": "codex" }
```

## Adapter Behavior

### Claude adapter (baseline)

- Launches `claude` in tmux or minimal mode.
- Uses Claude prompt detection for readiness.
- Maps `GT_SESSION_ID` <-> `CLAUDE_SESSION_ID`.
- Uses Claude settings hooks for events.

### Codex adapter (initial scope)

- Launches `codex` CLI with a bootstrap prompt that triggers `gt prime`.
- Uses warmup delay + retry as readiness strategy.
- Uses wrapper/shim to emit hook events on start/stop.
- Message delivery via stdin or CLI-supported mechanism.

## CLI Surface

- `gt runtime list`
- `gt runtime set default <name>`
- `gt runtime doctor` (optional checks per adapter)

Existing flags gain `--runtime` overrides where relevant:

- `gt start`, `gt sling`, `gt crew`, `gt polecat`, `gt witness`, `gt refinery`

## Migration Plan (v0.1)

1. Introduce `internal/runtime` package with interface definitions.
2. Implement `claude` adapter using current logic.
3. Route `internal/session` and `internal/tmux` to runtime interface.
4. Add `codex` adapter scaffolding (no behavior change to default).
5. Wire CLI flags and config resolution.
6. Update docs and examples.

## Refactor Map (touch points)

- `internal/session/manager.go`: call runtime interface for start/resume/send.
- `internal/tmux/tmux.go`: move Claude-specific readiness into adapter.
- `internal/claude/*`: migrate into `internal/runtime/claude`.
- `internal/cmd/*`: replace direct `claude` calls with runtime selection.
- `internal/config/*`: add runtime config resolution and defaults.
- `docs/INSTALLING.md`: add Codex optional install and runtime selection.
- `README.md`: update to “multi-runtime” wording and examples.

## Open Questions

- Codex prompt/ready detection specifics.
- Codex session resume semantics and transcript access.
- Whether hooks should be runtime-native or wrapper-only.

