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
		t.Fatalf("NamedUUIDs should be the same")
	}

	u3 := h.NamedUUID("something else")

	if u.String() == u3.String() {
		t.Fatalf("Diferent named UUIDS should result in different values")
	}
}

func TestHelperOneOf(t *testing.T) {
	h := &facto.Helper{Index: 1}

	t.Run("Strings", func(t *testing.T) {
		options := []string{"a", "b", "c"}
		s, ok := h.OneOf("a", "b", "c").(string)
		if !ok {
			t.Fatalf("OneOf should return a string")
		}

		if strings.Contains(strings.Join(options, "|"), s) == false {
			t.Fatalf("OneOf should return one of the options")
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
			t.Fatalf("OneOf should return a Status")
		}

		if s != StatusOpen && s != StatusClosed {
			t.Fatalf("OneOf should return one of the Status passed")
		}
	})

	t.Run("Empty", func(t *testing.T) {
		s := h.OneOf()
		if s != nil {
			t.Fatalf("OneOf should return Nil")
		}
	})

}
