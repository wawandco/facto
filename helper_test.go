package facto_test

import (
	"strings"
	"testing"

	"github.com/wawandco/facto"
)

func TestHelperNamedUUID(t *testing.T) {
	h := &facto.Helper{Index: 1}

	u := h.NamedUUID("something")
	u2 := h.NamedUUID("something")

	// Convert to UUID string format for comparison
	uuidStr := uuidBytesToString(u)
	uuidStr2 := uuidBytesToString(u2)

	if uuidStr != uuidStr2 {
		t.Errorf("NamedUUIDs should be the same")
	}

	u3 := h.NamedUUID("something else")
	uuidStr3 := uuidBytesToString(u3)

	if uuidStr == uuidStr3 {
		t.Errorf("Different named UUIDs should result in different values")
	}
}

func TestHelperOneOf(t *testing.T) {
	h := &facto.Helper{Index: 1}

	t.Run("Strings", func(t *testing.T) {
		options := []string{"a", "b", "c"}
		s, ok := h.OneOf("a", "b", "c").(string)
		if !ok {
			t.Errorf("OneOf should return a string")
		}

		if strings.Contains(strings.Join(options, "|"), s) == false {
			t.Errorf("OneOf should return one of the options")
		}
	})

	t.Run("CustomType", func(t *testing.T) {
		type Status string
		var (
			StatusOpen   Status = "open"
			StatusClosed Status = "closed"
		)

		s, ok := h.OneOf(StatusOpen, StatusClosed).(Status)
		if !ok {
			t.Errorf("OneOf should return a Status")
		}

		if s != StatusOpen && s != StatusClosed {
			t.Errorf("OneOf should return one of the Status passed")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		s := h.OneOf()
		if s != nil {
			t.Errorf("OneOf should return Nil")
		}
	})

}
