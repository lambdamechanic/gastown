// Package codex implements the AgentRuntime for Codex.
package codex

import (
	"context"
	"errors"
	"time"

	"github.com/steveyegge/gastown/internal/constants"
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
	if r.tmux == nil {
		return runtime.SessionHandle{}, errors.New("codex runtime requires tmux")
	}
	if opts.SessionID == "" {
		return runtime.SessionHandle{}, errors.New("codex runtime requires session id")
	}
	if opts.Command == "" {
		return runtime.SessionHandle{}, errors.New("codex runtime requires command")
	}

	if err := r.tmux.SendKeys(opts.SessionID, opts.Command); err != nil {
		return runtime.SessionHandle{}, err
	}

	_ = r.tmux.WaitForCommand(opts.SessionID, constants.SupportedShells, constants.ClaudeStartTimeout)
	time.Sleep(5 * time.Second)

	return runtime.SessionHandle{
		Runtime:   "codex",
		SessionID: opts.SessionID,
		WorkDir:   opts.WorkDir,
		StartedAt: time.Now(),
	}, nil
}

// Resume resumes a Codex session.
func (r *Runtime) Resume(ctx context.Context, handle runtime.SessionHandle) error {
	return errNotImplemented
}

// SendMessage sends a message to a Codex session.
func (r *Runtime) SendMessage(ctx context.Context, handle runtime.SessionHandle, msg runtime.Message) error {
	if r.tmux == nil {
		return errors.New("codex runtime requires tmux")
	}
	if msg.Delivery == runtime.DeliveryTmux || msg.Delivery == "" {
		return r.tmux.SendKeys(handle.SessionID, msg.Text)
	}
	return errNotImplemented
}

// Stop stops a Codex session.
func (r *Runtime) Stop(ctx context.Context, handle runtime.SessionHandle, reason string) error {
	if r.tmux == nil {
		return errors.New("codex runtime requires tmux")
	}
	return r.tmux.KillSession(handle.SessionID)
}

// IsReady checks if Codex is ready to receive input.
func (r *Runtime) IsReady(ctx context.Context, handle runtime.SessionHandle) (bool, error) {
	if r.tmux == nil {
		return false, errors.New("codex runtime requires tmux")
	}
	running, err := r.tmux.HasSession(handle.SessionID)
	if err != nil {
		return false, err
	}
	return running, nil
}

// DetectRunning checks if Codex is running for a session.
func (r *Runtime) DetectRunning(ctx context.Context, handle runtime.SessionHandle) (bool, error) {
	if r.tmux == nil {
		return false, errors.New("codex runtime requires tmux")
	}
	return r.tmux.HasSession(handle.SessionID)
}

// ListSessions lists Codex sessions.
func (r *Runtime) ListSessions(ctx context.Context, filter runtime.SessionFilter) ([]runtime.SessionHandle, error) {
	if r.tmux == nil {
		return nil, errors.New("codex runtime requires tmux")
	}

	sessions, err := r.tmux.ListSessions()
	if err != nil {
		return nil, err
	}

	handles := make([]runtime.SessionHandle, 0, len(sessions))
	for _, session := range sessions {
		if session == "" {
			continue
		}
		handles = append(handles, runtime.SessionHandle{
			Runtime:   "codex",
			SessionID: session,
		})
	}
	return handles, nil
}
func init() {
	runtime.Register("codex", func(t *tmux.Tmux) runtime.AgentRuntime {
		return New(t)
	})
}
