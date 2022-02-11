package facto

import (
	"fmt"
	"reflect"
)

var (
	// creator shared instance to be used by Create and CreateN
	creator Creator

	ErrNoCreatorDefined = fmt.Errorf("no creator set")
)

// SetCreator that CreateN and Create will use when invoked.
func SetCreator(c Creator) {
	creator = c
}

// Creator allows to persist the record in the database
// a creator is anything that has a Create method that receives
// an interface and returns an error in case it fails.
type Creator interface {
	Create(interface{}) error
}

// MemoryCreator is a creator that serves as a way to test
// the creation of records. Its not intended to be used in
// production as it is not thread safe. It is also intended
// to show how to implement a Creator in case don't have one.
type MemoryCreator struct {
	created []interface{}
}

func (mc *MemoryCreator) Create(i interface{}) error {
	mc.created = append(mc.created, i)

	return nil
}

func (mc *MemoryCreator) Contains(i interface{}) bool {
	if len(mc.created) == 0 {
		return false
	}

	for _, c := range mc.created {
		if reflect.DeepEqual(c, i) {
			return true
		}
	}

	return false
}
