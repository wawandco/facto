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

}

func TestHelperBuild(t *testing.T) {
	h := &facto.Helper{Index: 1}

	type s struct {
		Name string
	}

	facto.Register("something", func(h *facto.Helper) facto.Product {
		return s{
			Name: "Hello",
		}
	})

	p := h.Build("something").(s)
	if p.Name != "Hello" {
		t.Errorf("Name should be Hello got %v", p.Name)
	}
}
