# Facto

Facto is a fixtures library with a definition syntax. It allows Go developers to define fixtures in a declarative way using Go code instead of TOML/JSON/YAML. Inspired by [factory_bot](https://github.com/thoughtbot/factory_bot).

## Installation

```bash
go get github.com/wawandco/facto
```

## Your First Factory

Facto exposes a small API to define and use factories. Here's how to create a `User` factory:

```go
// factories/user.go
package factories

import (
    "github.com/wawandco/facto"
)

func UserFactory(h facto.Helper) facto.Product {
    user := User{
        Name: h.Faker.Name(),
    }

    return facto.Product(user)
}
```

Once defined, use factories in your tests:

```go
// users/user_test.go
package user_test

import (
    "your/package/factories"
)

func TestUser(t *testing.T) {
    // Build creates an object without persisting it
    user := facto.Build(factories.UserFactory).(models.User)
    err := db.Create(&user)
    // use user for test purposes ...

    // Create builds and persists the object to the database
    p, err := facto.Create(factories.UserFactory)
    user := p.(models.User)
    // user is already stored in the database
}
```

## Building Multiple Objects

When you need multiple instances, use `BuildN` or `CreateN`:

```go
// test/user_test.go
package user_test

func TestMultipleUsers(t *testing.T) {
    // Build 10 users without persisting
    users := facto.BuildN(factories.UserFactory, 10).([]models.User)
    
    // Create 10 users and persist them
    p, err := facto.CreateN(factories.UserFactory, 10)
    users = p.([]models.User)
}
```

### Using the Index

When running `BuildN`, use `h.Index` to differentiate objects:

```go
// factories/user.go
package factories

import (
    "fmt"
    "github.com/wawandco/facto"
)

func UserFactory(h facto.Helper) facto.Product {
    user := User{
        Name: fmt.Sprintf("User %d", h.Index),
    }

    return facto.Product(user)
}
```

`h.Index` starts at 0 when a factory is called individually.

## Dependent Factories

Build objects that depend on other objects:

```go
// factories/event.go
package factories

import (
    "github.com/wawandco/facto"
)

func UserFactory(h facto.Helper) facto.Product {
    user := User{
        Name: h.Faker.FirstName(),
    }

    return facto.Product(user)
}

func EventFactory(h facto.Helper) facto.Product {
    event := Event{
        // Build a User within the Event factory
        User: facto.Build(UserFactory).(User),
        Type: "Something",
    }

    return facto.Product(event)
}
```

## Named UUIDs

Reference the same ID across multiple factories using `NamedUUID`:

```go
// factories/user.go
package factories

import (
    "github.com/wawandco/facto"
)

func UserFactory(h facto.Helper) facto.Product {
    user := User{
        // Creates or retrieves a UUID named "owner_id"
        ID:   h.NamedUUID("owner_id"),
        Name: h.Faker.FirstName(),
    }

    return facto.Product(user)
}

// factories/event.go
package factories

import (
    "github.com/wawandco/facto"
)

func EventFactory(h facto.Helper) facto.Product {
    event := Event{
        // Retrieves the same "owner_id" UUID if UserFactory was called,
        // otherwise generates a new one
        UserID: h.NamedUUID("owner_id"),
        Type:   "Something",
    }

    return facto.Product(event)
}
```

### Custom UUID Generator

By default, `NamedUUID` returns a 16-byte slice. To use a different UUID library, set `h.UUIDGenerator`:

```go
// Using google/uuid
import "github.com/google/uuid"

func UserFactory(h facto.Helper) facto.Product {
    h.UUIDGenerator = func() []byte {
        return uuid.Must(uuid.NewRandom())[:]
    }

    user := User{
        ID: h.NamedUUID("user_id"),
    }
    return facto.Product(user)
}

// Or using gofrs/uuid
import "github.com/gofrs/uuid/v5"

func CompanyFactory(h facto.Helper) facto.Product {
    h.UUIDGenerator = func() []byte {
        u, _ := uuid.NewV4()
        return u.Bytes()
    }

    company := Company{
        ID: h.NamedUUID("company_id"),
    }
    return facto.Product(company)
}
```

If `UUIDGenerator` is nil (the default), Facto uses a built-in crypto/rand-based generator.

## Fake Data Generation

Generate realistic-looking test data using `h.Faker`:

```go
// factories/event.go
package factories

import (
    "github.com/wawandco/facto"
)

func EventFactory(h facto.Helper) facto.Product {
    event := Event{
        Name:         h.Faker.FirstName(),
        Type:         "Sports",
        ContactEmail: h.Faker.Email(),
        Company:      h.Faker.Company(),
        Address:      h.Faker.Address().Address,
    }

    return facto.Product(event)
}
```

Fake data is powered by [Gofakeit](https://github.com/brianvoe/gofakeit). See the [full list of generators](https://github.com/brianvoe/gofakeit/tree/v7.2.1#functions).

## Random Selection with OneOf

Randomly select from a list of values:

```go
// factories/event.go
package factories

import (
    "github.com/wawandco/facto"
)

func EventFactory(h facto.Helper) facto.Product {
    event := Event{
        Name: h.Faker.FirstName(),
        // Randomly picks one of the provided values
        Type: h.OneOf(TypeSports, TypeMusic, TypeConcert).(EventType),
    }

    return facto.Product(event)
}
```

## CLI

Install the Facto CLI to generate factory files:

```bash
go install github.com/wawandco/facto/cmd/facto@latest
```

Generate a factory:

```bash
facto generate user
# Generates factories/user.go
```

## Ox Plugin

Facto provides a plugin for the [Ox](https://github.com/wawandco/ox) CLI:

```go
// cmd/ox/main.go
package main

import (
    // ...
    fox "github.com/wawandco/facto/ox"
)

func main() {
    cli.Use(fox.Plugin{}) // Add the facto plugin
    
    err := cli.Run(context.Background(), os.Args)
    if err != nil {
        log.Fatal(err)
    }
}
```

Now `factory` appears as a generator:

```bash
$ ox generate

Available Generators:

  Name          Plugin
  ----          ------
  factory       facto  # <-- this one
```

Generate factories with Ox:

```bash
ox generate factory [name]
```