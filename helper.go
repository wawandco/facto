package facto

import (
	"sync"

	"github.com/gofrs/uuid"
)

// Helper gets injected into the factory and provides convenience methods
// for the fixtures.
type Helper struct {
	// Index is useful when creating N elements, this
	// Allow to differentiate the elements by the index.
	Index int

	// Faker instance to provide the ability to create fake data
	// on the factory.
	Faker faker

	// registry of the variables set in the helper, features like
	// NamedUUID depend on this registry.
	registry   map[string]interface{}
	registryMu sync.Mutex
}

// Build a factory.
func (h *Helper) Build(name string) Product {
	return Build(name)
}

// NamedUUID is a helper to create a UUID and keep it in the
// Facto context for later use.
func (h *Helper) NamedUUID(name string) uuid.UUID {
	name = "uuid_" + name
	if v := h.getVar(name); v != nil {
		return v.(uuid.UUID)
	}

	uuid := uuid.Must(uuid.NewV4())
	h.setVar(name, uuid)

	return uuid
}

func (h *Helper) setVar(key string, value interface{}) {
	h.registryMu.Lock()
	defer h.registryMu.Unlock()

	if h.registry == nil {
		h.registry = make(map[string]interface{})
	}

	h.registry[key] = value
}

func (h *Helper) getVar(key string) interface{} {
	return h.registry[key]
}
