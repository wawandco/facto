package facto

import (
	crand "crypto/rand"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

// Helper gets injected into the factory and provides convenience methods
// for the fixtures.
type Helper struct {
	// Faker instance to provide the ability to create fake data to the helper.
	Faker *gofakeit.Faker

	// Index is useful when creating N elements, this
	// Allow to differentiate the elements by the index.
	Index int

	// UUIDGenerator allows users to provide their own UUID generation function.
	// If nil, a default crypto/rand-based v4 UUID generator is used.
	UUIDGenerator func() []byte
}

// NamedUUID is a helper to create a UUID and keep it in the
// Facto context for later use this could come handy for database relations.
// Uses UUIDGenerator if set, otherwise generates a v4 UUID using crypto/rand.
func (h Helper) NamedUUID(name string) []byte {
	name = "uuid_" + name
	if v := defaultRegistry.getVariable(name); v != nil {
		return v.([]byte)
	}

	var b []byte
	if h.UUIDGenerator != nil {
		b = h.UUIDGenerator()
	} else {
		b = generateDefaultUUID()
	}

	defaultRegistry.setVariable(name, b)
	return b
}

// generateDefaultUUID creates a v4 UUID using crypto/rand
func generateDefaultUUID() []byte {
	b := make([]byte, 16)
	crand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant
	return b
}

// One of the passed elements, this method is useful when you
// have enum values and you want the falue of a field to be
// one of the possible values.
// e.g.
// ...
// u := User{
// // here the value of the field is one of the passed elements.
//
//		Status: OneOf(UserStatusActive, UserStatusInactive).(UserStatus)
//	}
func (h Helper) OneOf(values ...interface{}) interface{} {
	if len(values) == 0 {
		return nil
	}

	rand.Seed(time.Now().Unix())

	return values[rand.Intn(len(values))]
}

func NewHelper() Helper {
	return Helper{
		Faker: gofakeit.New(0),
	}
}
