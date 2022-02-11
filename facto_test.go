package facto_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/wawandco/facto"
)

func TestBuild(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		userProduct := facto.Build(UserFactory).(user)
		if userProduct.Name != "Wawandco" {
			t.Fatalf("expected '%s' but got '%s'", "Wawandco", userProduct.Name)
		}
	})

	t.Run("Dependent", func(t *testing.T) {
		event := facto.Build(EventFactory).(event)
		if event.Name != "CLICK" {
			t.Fatalf("expected '%s' but got '%s'", "CLICK", event.Name)
		}

		if event.User.Name != "Wawandco" {
			t.Fatalf("expected '%s' but got '%s'", "Wawandco", event.User.Name)
		}
	})

	t.Run("FakeAndBuildN", func(t *testing.T) {
		c := facto.Build(CompanyFactory).(company)
		if c.Name == "" {
			t.Fatalf("should have set the Name")
		}

		if c.Address == "" {
			t.Fatalf("should have set the Address")
		}

		if c.ContactEmail == "" {
			t.Fatalf("should have set the ContactEmail")
		}

		if len(c.Users) != 5 {
			t.Fatalf("should have set 5 users, set %v", len(c.Users))
		}

		for _, u := range c.Users {
			if u.Name != "" {
				continue
			}

			t.Fatalf("should have set the Name for User %v", u)
		}

	})

	t.Run("NamedUUID", func(t *testing.T) {
		c := facto.Build(DepartmentFactory).(department)
		if c.CompanyID.String() != c.Company.ID.String() {
			t.Fatalf("companyID should match")
		}
	})

}

func TestBuildFakeData(t *testing.T) {
	user := facto.Build(OtherUserFactory).(OtherUser)

	if user.FirstName == "" {
		t.Fatalf("should have set the FirstName")
	}

	if user.LastName == "" {
		t.Fatalf("should have set the LastName")
	}

	if user.Email == "" {
		t.Fatalf("should have set the Email")
	}

	if user.Company == "" {
		t.Fatalf("should have set the Company")
	}

	if user.Address == "" {
		t.Fatalf("should have set the Address")
	}
}

func Test_BuildN(t *testing.T) {
	usersProduct := facto.BuildN(UserNFactory, 5).([]user)

	for i := 0; i < 5; i++ {
		if fmt.Sprintf("Wawandco %d", i) != usersProduct[i].Name {
			t.Fatalf("expected '%s' but got '%s'", fmt.Sprintf("Wawandco %d", i), usersProduct[i].Name)
		}
	}
}

func Test_Create(t *testing.T) {
	mc := &facto.MemoryCreator{}

	tcases := []struct {
		name    string
		creator facto.Creator
		check   func(*testing.T, facto.Product, error)
	}{
		{
			name:    "no creator",
			creator: nil,
			check: func(tt *testing.T, p facto.Product, err error) {
				if err == nil {
					tt.Fatalf("expected an error")
				}

				if err != facto.ErrNoCreatorDefined {
					tt.Fatalf("expected error to be ErrNoCreator but got %v", err)
				}
			},
		},

		{
			name:    "creator defined",
			creator: mc,
			check: func(tt *testing.T, p facto.Product, err error) {
				if err != nil {
					tt.Fatalf("expected no error but got %v", err)
				}

				if !mc.Contains(p) {
					tt.Fatalf("expected product to be in memory creator")
				}
			},
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.name, func(tt *testing.T) {
			facto.SetCreator(tcase.creator)

			p, err := facto.Create(UserFactory)
			tcase.check(tt, p, err)
		})
	}
}

func Test_CreateN(t *testing.T) {
	mc := &facto.MemoryCreator{}

	tcases := []struct {
		name    string
		creator facto.Creator
		check   func(*testing.T, facto.Product, error)
	}{
		{
			name:    "no creator",
			creator: nil,
			check: func(tt *testing.T, p facto.Product, err error) {
				if err == nil {
					tt.Fatalf("expected an error")
				}

				if err != facto.ErrNoCreatorDefined {
					tt.Fatalf("expected error to be ErrNoCreator but got %v", err)
				}
			},
		},

		{
			name:    "creator defined",
			creator: mc,
			check: func(tt *testing.T, p facto.Product, err error) {
				if err != nil {
					tt.Fatalf("expected no error but got %v", err)
				}

				if !mc.Contains(p) {
					tt.Fatalf("expected product to be in memory creator")
				}
			},
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.name, func(tt *testing.T) {
			facto.SetCreator(tcase.creator)

			p, err := facto.CreateN(UserFactory, 10)
			tcase.check(tt, p, err)
		})
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
				t.Fatalf("expected '%s' but got '%s' in '%s'", expected, userProduct.Name, fmt.Sprintf("case %d", i+1))
			}
		}

		go gr(tcases[i].factoryName, tcases[i].factory, tcases[i].expected, i)
	}
	wgbuild.Wait()
}
