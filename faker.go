package facto

import (
	"github.com/brianvoe/gofakeit/v6"
)

// faker is in charge of generating fake data. It contains methods
// for each of the fake data types it can fake.
type faker struct{}

// FirstName faked.
func (f faker) FirstName() string {
	return gofakeit.FirstName()
}

// LastName faked.
func (f faker) LastName() string {
	return gofakeit.LastName()
}

// Email faked.
func (f faker) Email() string {
	return gofakeit.Email()
}

// Company faked.
func (f faker) Company() string {
	return gofakeit.Company()
}

// Address faked.
func (f faker) Address() string {
	return gofakeit.Address().Address
}
