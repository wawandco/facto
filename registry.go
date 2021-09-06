package facto

import "sync"

var (
	defaultRegistry = &registry{}
)

// Registry holds factories and variables in memory for
// its use by other components.
type registry struct {
	factories map[string]Factory
	variables map[string]interface{}

	mu sync.RWMutex
}

func (r *registry) setVariable(key string, value interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.variables == nil {
		r.variables = make(map[string]interface{})
	}

	r.variables[key] = value
}

func (r *registry) getVariable(key string) interface{} {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.variables[key]
}
