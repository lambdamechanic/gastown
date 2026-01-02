// Package claude implements the AgentRuntime for Claude Code.
package claude

import (
	"context"
	"errors"
	"time"

	"github.com/steveyegge/gastown/internal/constants"
	"github.com/steveyegge/gastown/internal/runtime"
	"github.com/steveyegge/gastown/internal/tmux"
)

// Runtime is the Claude Code runtime adapter.
type Runtime struct {
	tmux          *tmux.Tmux
	Command       string
	Args          []string
	ReadinessMode string
}

// New returns a Claude runtime adapter bound to a tmux instance.
func New(t *tmux.Tmux) *Runtime {
	return &Runtime{
		tmux:          t,
		ReadinessMode: runtime.ReadinessPrompt,
	}
}

var errNotImplemented = errors.New("claude runtime adapter not wired for this operation")

// Start starts a Claude session in an existing tmux session.
func (r *Runtime) Start(ctx context.Context, opts runtime.StartOptions) (runtime.SessionHandle, error) {
	if r.tmux == nil {
		return runtime.SessionHandle{}, errors.New("claude runtime requires tmux")
	}
	if opts.SessionID == "" {
		return runtime.SessionHandle{}, errors.New("claude runtime requires session id")
	}
	if opts.Command == "" {
		return runtime.SessionHandle{}, errors.New("claude runtime requires command")
	}

	if err := r.tmux.SendKeys(opts.SessionID, opts.Command); err != nil {
		return runtime.SessionHandle{}, err
	}

	// Non-fatal: Claude might still be starting.
	_ = r.tmux.WaitForCommand(opts.SessionID, constants.SupportedShells, constants.ClaudeStartTimeout)

	// Conservative warmup to avoid prompt detection false positives.
	time.Sleep(10 * time.Second)

	return runtime.SessionHandle{
		Runtime:   "claude",
		SessionID: opts.SessionID,
		WorkDir:   opts.WorkDir,
		StartedAt: time.Now(),
	}, nil
}

// Resume resumes a Claude session.
func (r *Runtime) Resume(ctx context.Context, handle runtime.SessionHandle) error {
	return errNotImplemented
}

// SendMessage sends a message to a Claude session.
func (r *Runtime) SendMessage(ctx context.Context, handle runtime.SessionHandle, msg runtime.Message) error {
	if r.tmux == nil {
		return errors.New("claude runtime requires tmux")
	}
	if msg.Delivery != "" && msg.Delivery != runtime.DeliveryTmux {
		return errors.New("claude runtime only supports tmux delivery")
	}
	return r.tmux.NudgeSession(handle.SessionID, msg.Text)
}

// Stop stops a Claude session.
func (r *Runtime) Stop(ctx context.Context, handle runtime.SessionHandle, reason string) error {
	if r.tmux == nil {
		return errors.New("claude runtime requires tmux")
	}
	return r.tmux.KillSession(handle.SessionID)
}

// IsReady checks if Claude is ready to receive input.
func (r *Runtime) IsReady(ctx context.Context, handle runtime.SessionHandle) (bool, error) {
	if r.tmux == nil {
		return false, errors.New("claude runtime requires tmux")
	}
	if err := r.tmux.WaitForClaudeReady(handle.SessionID, 2*time.Second); err != nil {
		return false, nil
	}
	return true, nil
}

// DetectRunning checks if Claude is running for a session.
func (r *Runtime) DetectRunning(ctx context.Context, handle runtime.SessionHandle) (bool, error) {
	if r.tmux == nil {
		return false, errors.New("claude runtime requires tmux")
	}
	return r.tmux.IsClaudeRunning(handle.SessionID), nil
}

// ListSessions lists Claude sessions.
func (r *Runtime) ListSessions(ctx context.Context, filter runtime.SessionFilter) ([]runtime.SessionHandle, error) {
	if r.tmux == nil {
		return nil, errors.New("claude runtime requires tmux")
	}

	var sessions []string
	var err error
	if filter.WorkDir != "" {
		sessions, err = r.tmux.FindSessionByWorkDir(filter.WorkDir, true)
		if err != nil {
			return nil, err
		}
	} else {
		sessions, err = r.tmux.ListSessions()
		if err != nil {
			return nil, err
		}
	}

	handles := make([]runtime.SessionHandle, 0, len(sessions))
	for _, session := range sessions {
		if session == "" {
			continue
		}
		handles = append(handles, runtime.SessionHandle{
			Runtime:   "claude",
			SessionID: session,
		})
	}
	return handles, nil
}
