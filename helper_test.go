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

	if u.String() != u2.String() {
		t.Errorf("NamedUUIDs should be the same")
	}

	u3 := h.NamedUUID("something else")

	if u.String() == u3.String() {
		t.Errorf("Diferent named UUIDS should result in different values")
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
