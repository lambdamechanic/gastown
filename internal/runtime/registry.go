package runtime

import (
	"fmt"
	"sync"

	"github.com/steveyegge/gastown/internal/tmux"
)

var (
	registryMu sync.RWMutex
	registry   = make(map[string]func(*tmux.Tmux) AgentRuntime)
)

// Register adds a runtime adapter by name.
func Register(name string, factory func(*tmux.Tmux) AgentRuntime) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = factory
}

// Get returns a registered runtime adapter by name.
func Get(name string, t *tmux.Tmux) (AgentRuntime, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	if factory, ok := registry[name]; ok {
		return factory(t), nil
	}
	return nil, fmt.Errorf("runtime not registered: %s", name)
}

// Names returns the list of registered runtime names.
func Names() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()
	out := make([]string, 0, len(registry))
	for name := range registry {
		out = append(out, name)
	}
	return out
}
