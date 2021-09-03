package facto_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/paganotoni/facto"
)

type User struct {
	Name string
}

func TestBuild(t *testing.T) {
	type event struct {
		Name string
		User User
	}

	type company struct {
		Name         string
		Address      string
		ContactEmail string
		Users        []User
	}

	facto.Register("User", func(f facto.Helper) facto.Product {
		u := User{
			Name: "Wawandco",
		}
		return facto.Product(u)
	})

	facto.Register("Event", func(f facto.Helper) facto.Product {
		u := event{
			Name: "CLICK",
			User: facto.Build("User").(User),
		}

		return facto.Product(u)
	})

	facto.Register("company", func(h facto.Helper) facto.Product {
		u := company{
			Name:         h.Faker.Company(),
			Address:      h.Faker.Address(),
			ContactEmail: h.Faker.Email(),

			// Building N Users
			Users: facto.BuildN("User", 5).([]User),
		}

		return facto.Product(u)
	})

	facto.Register("Company", func(f facto.Helper) facto.Product {
		u := event{
			Name: "CLICK",
			User: facto.Build("User").(User),
		}

		return facto.Product(u)
	})

	t.Run("Simple", func(t *testing.T) {
		userProduct := facto.Build("User").(User)
		if userProduct.Name != "Wawandco" {
			t.Errorf("expected '%s' but got '%s'", "Wawandco", userProduct.Name)
		}
	})

	t.Run("Dependent", func(t *testing.T) {
		event := facto.Build("Event").(event)
		if event.Name != "CLICK" {
			t.Errorf("expected '%s' but got '%s'", "CLICK", event.Name)
		}

		if event.User.Name != "Wawandco" {
			t.Errorf("expected '%s' but got '%s'", "Wawandco", event.User.Name)
		}
	})

	t.Run("FakeAndBuildN", func(t *testing.T) {
		c := facto.Build("company").(company)
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

}

func TestBuildFakeData(t *testing.T) {
	type OtherUser struct {
		FirstName string
		LastName  string
		Email     string
		Company   string
		Address   string
	}

	facto.Register("User", func(f facto.Helper) facto.Product {
		u := OtherUser{
			FirstName: f.Faker.FirstName(),
			LastName:  f.Faker.LastName(),
			Email:     f.Faker.Email(),
			Company:   f.Faker.Company(),
			Address:   f.Faker.Address(),
		}

		return facto.Product(u)
	})

	user := facto.Build("User").(OtherUser)

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
	facto.Register("Users", func(f facto.Helper) facto.Product {
		u := User{
			Name: fmt.Sprintf("Wawandco %d", f.Index),
		}
		return facto.Product(u)
	})

	usersProduct := facto.BuildN("Users", 5).([]User)

	for i := 0; i < 5; i++ {
		if fmt.Sprintf("Wawandco %d", i+1) != usersProduct[i].Name {
			t.Errorf("expected '%s' but got '%s'", fmt.Sprintf("Wawandco %d", i+1), usersProduct[i].Name)
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
				u := User{
					Name: "Wawandco",
				}
				return facto.Product(u)
			},
			expected: "Wawandco",
		},

		{
			factoryName: "UserNumberTwo",
			factory: func(f facto.Helper) facto.Product {
				u := User{
					Name: "Wawandco 2",
				}
				return facto.Product(u)
			},
			expected: "Wawandco 2",
		},

		{
			factoryName: "UserNumberThree",
			factory: func(f facto.Helper) facto.Product {
				u := User{
					Name: fmt.Sprintf("Wawandco %d", f.Index),
				}
				return facto.Product(u)
			},
			expected: "Wawandco 1",
		},
	}

	var wgreg sync.WaitGroup
	for i := range tcases {
		wgreg.Add(1)
		gr := func(fname string, factory facto.Factory) {
			defer wgreg.Done()

			facto.Register(fname, factory)
		}

		go gr(tcases[i].factoryName, tcases[i].factory)
	}
	wgreg.Wait()

	var wgbuild sync.WaitGroup
	for i := range tcases {
		wgbuild.Add(1)

		gr := func(name, expected string, index int) {
			defer wgbuild.Done()

			userProduct, ok := facto.Build(name).(User)
			if !ok {
				t.Fatalf("Should have got user but got %v", userProduct)
			}

			if expected != userProduct.Name {
				t.Errorf("expected '%s' but got '%s' in '%s'", expected, userProduct.Name, fmt.Sprintf("case %d", i+1))
			}
		}

		go gr(tcases[i].factoryName, tcases[i].expected, i)
	}
	wgbuild.Wait()
}
