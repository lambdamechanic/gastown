// Package runtime defines the interface for agent runtimes (Claude, Codex, etc.).
package runtime

import (
	"context"
	"time"
)

// AgentRuntime abstracts lifecycle operations for an agent runtime.
type AgentRuntime interface {
	Start(ctx context.Context, opts StartOptions) (SessionHandle, error)
	Resume(ctx context.Context, handle SessionHandle) error
	SendMessage(ctx context.Context, handle SessionHandle, msg Message) error
	Stop(ctx context.Context, handle SessionHandle, reason string) error
	IsReady(ctx context.Context, handle SessionHandle) (bool, error)
	DetectRunning(ctx context.Context, handle SessionHandle) (bool, error)
	ListSessions(ctx context.Context, filter SessionFilter) ([]SessionHandle, error)
}

// StartOptions describes a new runtime session request.
type StartOptions struct {
	WorkDir       string
	RuntimeName   string
	AccountDir    string
	Env           map[string]string
	InitialPrompt string
	Mode          string // "minimal" | "tmux"
}

// SessionHandle describes a running runtime session.
type SessionHandle struct {
	Runtime   string
	SessionID string
	WorkDir   string
	PID       int
	StartedAt time.Time
	ReadyAt   time.Time
}

// Message represents a runtime-agnostic message delivery request.
type Message struct {
	Text     string
	Delivery string // "stdin" | "tmux" | "rpc"
	Timeout  time.Duration
}

// SessionFilter scopes ListSessions results.
type SessionFilter struct {
	Runtime string
	WorkDir string
}
