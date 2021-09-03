package facto_test

import (
	"testing"

	"github.com/paganotoni/facto"
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
