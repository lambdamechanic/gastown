package runtime

import (
	"fmt"
	"sync"
)

var (
	registryMu sync.RWMutex
	registry   = make(map[string]AgentRuntime)
)

// Register adds a runtime adapter by name.
func Register(name string, rt AgentRuntime) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = rt
}

// Get returns a registered runtime adapter by name.
func Get(name string) (AgentRuntime, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	if rt, ok := registry[name]; ok {
		return rt, nil
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
