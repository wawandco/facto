package facto_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/wawandco/facto"
)

func TestBuild(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		userProduct := facto.Build(UserFactory).(user)
		if userProduct.Name != "Wawandco" {
			t.Errorf("expected '%s' but got '%s'", "Wawandco", userProduct.Name)
		}
	})

	t.Run("Dependent", func(t *testing.T) {
		event := facto.Build(EventFactory).(event)
		if event.Name != "CLICK" {
			t.Errorf("expected '%s' but got '%s'", "CLICK", event.Name)
		}

		if event.User.Name != "Wawandco" {
			t.Errorf("expected '%s' but got '%s'", "Wawandco", event.User.Name)
		}
	})

	t.Run("FakeAndBuildN", func(t *testing.T) {
		c := facto.Build(CompanyFactory).(company)
		if c.Name == "" {
			t.Errorf("should have set the Name")
		}

		if c.Address == "" {
			t.Errorf("should have set the Address")
		}

		if c.ContactEmail == "" {
			t.Errorf("should have set the ContactEmail")
		}

		if len(c.Users) != 5 {
			t.Errorf("should have set 5 users, set %v", len(c.Users))
		}

		for _, u := range c.Users {
			if u.Name != "" {
				continue
			}

			t.Errorf("should have set the Name for User %v", u)
		}

	})

	t.Run("NamedUUID", func(t *testing.T) {
		c := facto.Build(DepartmentFactory).(department)
		if c.CompanyID.String() != c.Company.ID.String() {
			t.Errorf("companyID should match")
		}
	})

}

func TestBuildFakeData(t *testing.T) {
	user := facto.Build(OtherUserFactory).(OtherUser)

	if user.FirstName == "" {
		t.Errorf("should have set the FirstName")
	}

	if user.LastName == "" {
		t.Errorf("should have set the LastName")
	}

	if user.Email == "" {
		t.Errorf("should have set the Email")
	}

	if user.Company == "" {
		t.Errorf("should have set the Company")
	}

	if user.Address == "" {
		t.Errorf("should have set the Address")
	}
}

func Test_BuildN(t *testing.T) {
	usersProduct := facto.BuildN(UserNFactory, 5).([]user)

	for i := 0; i < 5; i++ {
		if fmt.Sprintf("Wawandco %d", i) != usersProduct[i].Name {
			t.Errorf("expected '%s' but got '%s'", fmt.Sprintf("Wawandco %d", i), usersProduct[i].Name)
		}
	}
}

func Test_Build_Concurrently(t *testing.T) {
	tcases := []struct {
		factoryName string
		factory     facto.Factory
		expected    string
	}{
		{
			factoryName: "UserNumberOne",
			factory: func(f facto.Helper) facto.Product {
				u := user{
					Name: "Wawandco",
				}
				return facto.Product(u)
			},
			expected: "Wawandco",
		},

		{
			factoryName: "UserNumberTwo",
			factory: func(f facto.Helper) facto.Product {
				u := user{
					Name: "Wawandco 2",
				}
				return facto.Product(u)
			},
			expected: "Wawandco 2",
		},

		{
			factoryName: "UserNumberThree",
			factory: func(f facto.Helper) facto.Product {
				u := user{
					Name: fmt.Sprintf("Wawandco %d", f.Index),
				}
				return facto.Product(u)
			},
			expected: "Wawandco 0",
		},
	}

	var wgbuild sync.WaitGroup
	for i := range tcases {
		wgbuild.Add(1)

		gr := func(name string, factory facto.Factory, expected string, index int) {
			defer wgbuild.Done()

			userProduct, ok := facto.Build(factory).(user)
			if !ok {
				t.Fatalf("Should have got user but got %v", userProduct)
			}

			if expected != userProduct.Name {
				t.Errorf("expected '%s' but got '%s' in '%s'", expected, userProduct.Name, fmt.Sprintf("case %d", i+1))
			}
		}

		go gr(tcases[i].factoryName, tcases[i].factory, tcases[i].expected, i)
	}
	wgbuild.Wait()
}

type user struct {
	Name string
}

type event struct {
	Name string
	User user
}

type company struct {
	ID           uuid.UUID
	Name         string
	Address      string
	ContactEmail string
	Users        []user
}

type department struct {
	Name      string
	CompanyID uuid.UUID
	Company   company
}

type OtherUser struct {
	FirstName string
	LastName  string
	Email     string
	Company   string
	Address   string
}

func UserNFactory(f facto.Helper) facto.Product {
	u := user{
		Name: fmt.Sprintf("Wawandco %d", f.Index),
	}
	return facto.Product(u)
}

func UserFactory(f facto.Helper) facto.Product {
	u := user{
		Name: "Wawandco",
	}
	return facto.Product(u)
}

func OtherUserFactory(f facto.Helper) facto.Product {
	u := OtherUser{
		FirstName: f.Faker.FirstName(),
		LastName:  f.Faker.LastName(),
		Email:     f.Faker.Email(),
		Company:   f.Faker.Company(),
		Address:   f.Faker.Address(),
	}

	return facto.Product(u)
}

func EventFactory(f facto.Helper) facto.Product {
	u := event{
		Name: "CLICK",
		User: facto.Build(UserFactory).(user),
	}

	return facto.Product(u)
}

func CompanyFactory(h facto.Helper) facto.Product {
	u := company{
		ID:           h.NamedUUID("company_id"),
		Name:         h.Faker.Company(),
		Address:      h.Faker.Address(),
		ContactEmail: h.Faker.Email(),

		// Building N Users
		Users: facto.BuildN(UserFactory, 5).([]user),
	}

	return facto.Product(u)
}

func DepartmentFactory(f facto.Helper) facto.Product {
	u := department{
		Name:      "Technology",
		CompanyID: f.NamedUUID("company_id"),
		Company:   facto.Build(CompanyFactory).(company),
	}

	return facto.Product(u)
}
