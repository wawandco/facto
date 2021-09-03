package facto

import "testing"

func TestHelper(t *testing.T) {
	h := Helper{Index: 1}

	if h.Faker.Email() == "" {
		t.Error("Should return an email")
	}
}
