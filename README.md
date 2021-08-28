# Facto

Facto is a fixtures library with a definition syntax. It aims to allow Go developers to define fixtures in a declarative way and using the Go language as the means to describe the fixtures instead of TOML/JSON/YAML. It is inspired on factory_bot.

### Getting started

To get started install the Facto CLI, it will allow you to perform basic operations on your fixtures.

```sh
go install github.com/paganotoni/facto/cmd/facto@latest
```

Once initialized you can generate your first factory with:

```sh
facto generate user
```

This generates the `factories/user.go` and (`factories/factories.go` in case it does not exist).

### Your first Factory

Facto exposes a small API to register and use factories. The following example shows how to create a factory named `User`:

```go
//  in factories/user.go
package factories

import (
    "github.com/paganotoni/facto"
)

func init() {
    // User is how the factory will be invoked.
    facto.Register("User", UserFactory)
}

func UserFactory(f facto.Helper) interface{} {
    return &User{
        Name: f.Fake.Name(),
    }
}
```

Then on your test file you can use the factory by calling:

```go
//  in test/user_test.go
package user_test

func TestUser(t *testing.T) {
    user := facto.Build("User").(models.User)
    // use user for test purposes ...
}
```

### Building N objects
Sometimes you need to build more than one instance of an object, in that case. Facto provides a helper to do that:

```go
//  in test/user_test.go
package user_test

func TestUser(t *testing.T) {
    users := facto.BuildN("User").([]models.User)
    // use users for test purposes ...
}

```

### Dependent factories

Another case is when you need to build an object that depends on another object. You can use the `factory.Helper` parameter on your factory to get the object you need:

```go
//  in factories/event.go
package factories
import (
    "github.com/paganotoni/facto"
)

func init() {
    facto.Register("User", UserFactory)
    facto.Register("Event", UserFactory)
}

func UserFactory(f facto.Helper) interface{} {
    return &User{
        Name: f.Fake.Name(),
    }
}

func EventFactory(f facto.Helper) interface{} {
	return Event{
        // Here we pull the User from the helper given we know there is a 
        // factory for it.
		User: f.Build("User").(User),
		Type: "Something",
	}
})
```

Another case here is when you need to reference the ID of the previous object. You can use the `factory.Helper` parameter on your factory to build a Named UUID:

```go
//  in factories/user.go
package factories
import (
    "github.com/paganotoni/facto"
)

func init() {
    facto.Register("User", UserFactory)
}

func UserFactory(f facto.Helper) interface{} {
    return &User{
        // owner_id will be assigned the generated UUID 
        // and any 
        ID: f.NamedUUID("owner_id"),
        Name: f.Fake.Name(),
    }
}

//  in factories/event.go
package factories
import (
    "github.com/paganotoni/facto"
)

func init() {
    facto.Register("Event", EventFactory)
}

func EventFactory(f facto.Helper) interface{} {
	return Event{
        // Here we pull the User from the helper given we know there is a 
        // factory that adds it. If there was not one it would be generated new.
		UserID: f.NamedUUID("owner_id"),
		Type: "Something",
	}
})
```

### Faking data

Sometimes you need to generate data that is not realistic, in that case you can use the `factory.Fake` method on your factory to generate fake data, for example:

```go
//  in factories/event.go
package factories
import (
    "github.com/paganotoni/facto"
)

func init() {
    facto.Register("Event", EventFactory)
}

func EventFactory(f facto.Helper) interface{} {
	return Event{
        Name: f.Fake("Name"),
		Type: "Sports",

        ContactEmail: f.Fake("Email"),
        Company: f.Fake("Company"),
        Address: f.Fake("Address"),
	}
})
```

The full list of available fake data generators can be found in [here](link-to-repo).

### The CLI

The Facto CLI is a simple command line tool that allows you to generate fixtures and initialize your project for fixtures generation. It contains the following commands to facilitate the use of Facto:

 * `facto generate`: Generates fixtures for your project.
 * `facto list --help`: Shows the list of factories in the current codebase.
 * `facto version`: The version of the Facto CLI.



