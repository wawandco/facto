package facto_test

import (
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/wawandco/facto"
)

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
		Address:   f.Faker.Address().Address,
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
		Name:         h.Faker.Name(),
		Address:      h.Faker.Address().Address,
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
