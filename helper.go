package facto

import (
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
}

// Build a factory.
func (h Helper) Build(name string) Product {
	return Build(name)
}

// NamedUUID is a helper to create a UUID and keep it in the
// Facto context for later use.
func (h Helper) NamedUUID(name string) uuid.UUID {
	name = "uuid_" + name
	if v := defaultRegistry.getVariable(name); v != nil {
		return v.(uuid.UUID)
	}

	uuid := uuid.Must(uuid.NewV4())
	defaultRegistry.setVariable(name, uuid)

	return uuid
}
