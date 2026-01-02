// Package codex implements the AgentRuntime for Codex.
package codex

import (
	"context"
	"errors"

	"github.com/steveyegge/gastown/internal/runtime"
	"github.com/steveyegge/gastown/internal/tmux"
)

// Runtime is the Codex runtime adapter.
type Runtime struct {
	tmux          *tmux.Tmux
	Command       string
	Args          []string
	ReadinessMode string
}

// New returns a Codex runtime adapter bound to a tmux instance.
func New(t *tmux.Tmux) *Runtime {
	return &Runtime{
		tmux:          t,
		ReadinessMode: runtime.ReadinessWarmup,
	}
}

var errNotImplemented = errors.New("codex runtime adapter not wired yet")

// Start starts a Codex session.
func (r *Runtime) Start(ctx context.Context, opts runtime.StartOptions) (runtime.SessionHandle, error) {
	return runtime.SessionHandle{}, errNotImplemented
}

// Resume resumes a Codex session.
func (r *Runtime) Resume(ctx context.Context, handle runtime.SessionHandle) error {
	return errNotImplemented
}

// SendMessage sends a message to a Codex session.
func (r *Runtime) SendMessage(ctx context.Context, handle runtime.SessionHandle, msg runtime.Message) error {
	return errNotImplemented
}

// Stop stops a Codex session.
func (r *Runtime) Stop(ctx context.Context, handle runtime.SessionHandle, reason string) error {
	return errNotImplemented
}

// IsReady checks if Codex is ready to receive input.
func (r *Runtime) IsReady(ctx context.Context, handle runtime.SessionHandle) (bool, error) {
	return false, errNotImplemented
}

// DetectRunning checks if Codex is running for a session.
func (r *Runtime) DetectRunning(ctx context.Context, handle runtime.SessionHandle) (bool, error) {
	return false, errNotImplemented
}

// ListSessions lists Codex sessions.
func (r *Runtime) ListSessions(ctx context.Context, filter runtime.SessionFilter) ([]runtime.SessionHandle, error) {
	return nil, errNotImplemented
}
