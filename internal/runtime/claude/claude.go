// Package claude implements the AgentRuntime for Claude Code.
package claude

import (
	"context"
	"errors"

	"github.com/steveyegge/gastown/internal/runtime"
)

// Runtime is the Claude Code runtime adapter.
type Runtime struct {
	Command        string
	Args           []string
	ReadinessStyle string
}

var errNotImplemented = errors.New("claude runtime adapter not wired yet")

// Start starts a Claude session.
func (r *Runtime) Start(ctx context.Context, opts runtime.StartOptions) (runtime.SessionHandle, error) {
	return runtime.SessionHandle{}, errNotImplemented
}

// Resume resumes a Claude session.
func (r *Runtime) Resume(ctx context.Context, handle runtime.SessionHandle) error {
	return errNotImplemented
}

// SendMessage sends a message to a Claude session.
func (r *Runtime) SendMessage(ctx context.Context, handle runtime.SessionHandle, msg runtime.Message) error {
	return errNotImplemented
}

// Stop stops a Claude session.
func (r *Runtime) Stop(ctx context.Context, handle runtime.SessionHandle, reason string) error {
	return errNotImplemented
}

// IsReady checks if Claude is ready to receive input.
func (r *Runtime) IsReady(ctx context.Context, handle runtime.SessionHandle) (bool, error) {
	return false, errNotImplemented
}

// DetectRunning checks if Claude is running for a session.
func (r *Runtime) DetectRunning(ctx context.Context, handle runtime.SessionHandle) (bool, error) {
	return false, errNotImplemented
}

// ListSessions lists Claude sessions.
func (r *Runtime) ListSessions(ctx context.Context, filter runtime.SessionFilter) ([]runtime.SessionHandle, error) {
	return nil, errNotImplemented
}

func init() {
	runtime.Register("claude", &Runtime{})
}
