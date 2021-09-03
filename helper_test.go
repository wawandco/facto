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

func TestHelperBuild(t *testing.T) {
	h := &facto.Helper{Index: 1}

	type parent struct {
		Name string
	}

	type s struct {
		Name   string
		Parent parent
	}

	facto.Register("parent", func(h facto.Helper) facto.Product {
		return parent{
			Name: "Parent",
		}
	})

	facto.Register("something", func(h facto.Helper) facto.Product {
		return s{
			Name:   "Child",
			Parent: h.Build("parent").(parent),
		}
	})

	p := h.Build("something").(s)
	if p.Name != "Child" {
		t.Errorf("Name should be Hello got %v", p.Name)
	}

	if p.Parent.Name != "Parent" {
		t.Errorf("Should have build Parent %v", p.Name)
	}
}
