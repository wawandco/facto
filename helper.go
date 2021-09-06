package facto

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofrs/uuid"
)

// Helper gets injected into the factory and provides convenience methods
// for the fixtures.
type Helper struct {
	// Faker instance to provide the ability to create fake data to the helper.
	Faker *gofakeit.Faker

	// Index is useful when creating N elements, this
	// Allow to differentiate the elements by the index.
	Index int
}

// NamedUUID is a helper to create a UUID and keep it in the
// Facto context for later use this could come handy for database relations.
func (h Helper) NamedUUID(name string) uuid.UUID {
	name = "uuid_" + name
	if v := defaultRegistry.getVariable(name); v != nil {
		return v.(uuid.UUID)
	}

	uuid := uuid.Must(uuid.NewV4())
	defaultRegistry.setVariable(name, uuid)

	return uuid
}

// One of the passed elements, this method is useful when you
// have enum values and you want the falue of a field to be
// one of the possible values.
// e.g.
// ...
// u := User{
// // here the value of the field is one of the passed elements.
//	Status: OneOf(UserStatusActive, UserStatusInactive).(UserStatus)
// }
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
