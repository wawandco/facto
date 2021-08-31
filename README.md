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

This generates `factories/factories.go` which looks like this:
```go
//  in factories/user.go
package factories

import (
    "github.com/paganotoni/facto"
)

// Load all of the factories into facto to make them available for tests.
// Inside this function we will add the rest of the factories, 
// Facto will add newly generated factories to this function. Then 
// Its expected that in your tests you invoke this function, p.e: factories.Load()
func Load() {
    // User is how the factory will be invoked, p.e. facto.Build("User)
    facto.Register("User", UserFactory)
}

```

As well as `factories/user.go`:

```go
//  in factories/user.go
package factories

import (
    "github.com/paganotoni/facto"
)

func UserFactory(f facto.Helper) facto.Product {
    // defaults to nil, once generated you will need to add the 
    // fixture logic in here.
    return nil
}
```

It's important that before your tests you call factories.Load() to load your factories before the test. Some testing libraries provide a way to do this on a single place.

### Your first Factory

Facto exposes a small API to register and use factories. The following example shows how to create a factory named `User`:

```go
//  in factories/user.go
package factories

import (
    "github.com/paganotoni/facto"
)

func UserFactory(f facto.Helper) facto.Product {
    user := User{
        Name: f.Fake.Name(),
    }

    return facto.Product(user)
}
```

Then on your test file you can use the factory by calling:

```go
//  in test/user_test.go
package user_test

import (
    "your/package/factories"
)

func TestUser(t *testing.T) {
    // Load the factories
    factories.Load()

    user := facto.Build("User").(models.User)
    // use user for test purposes ...
}
```

### Building N   
Sometimes you need to build more than one instance of an object, in that case. Facto provides a helper to do that:

```go
//  in test/user_test.go
package user_test

func TestUser(t *testing.T) {
    users := facto.BuildN("User", 10).([]models.User)
    // use users for test purposes ...
}

```

When running `BuildN` you can use the Index variable to get the index of the object being built. This is useful when you need to build a list of objects and you want to know which one is being built.

```go
//  in factories/user.go
package factories

import (
    "fmt"
    "github.com/paganotoni/facto"
)

func UserFactory(f facto.Helper) facto.Product {
    user := User{
        Name: fmt.Sprintf("User %d", f.Index),
    }

    return facto.Product(user)
}
```

When the factory gets called individually Index will be 0.

### Dependent factories

Another case is when you need to build an object that depends on another object. You can use the `factory.Helper` parameter on your factory to get the object you need:

```go
//  in factories/event.go
package factories
import (
    "github.com/paganotoni/facto"
)

func UserFactory(f facto.Helper) facto.Product {
    user := User{
        Name: f.Fake.Name(),
    }

    return facto.Product(user)
}

func EventFactory(f facto.Helper) facto.Product {
    event := Event{
        // Here we pull the User from the helper given we know there is a 
        // factory for it, assuming it was registered in the factories.Load() as `User` we
        // can use it like this:
        User: f.Build("User").(User),
        Type: "Something",
    }

    return facto.Product(event)
})
```

Another case here is when you need to reference the ID of the previous object. You can use the `factory.Helper` parameter on your factory to build a Named UUID:

```go
//  in factories/user.go
package factories
import (
    "github.com/paganotoni/facto"
)

func UserFactory(f facto.Helper) facto.Product {
    user := User{
        // owner_id will be assigned the generated UUID 
        // and any 
        ID: f.NamedUUID("owner_id"),
        Name: f.Fake.Name(),
    }
    
    return facto.Product(user)
}

//  in factories/event.go
package factories
import (
    "github.com/paganotoni/facto"
)

func EventFactory(f facto.Helper) facto.Product {
    event := Event{
        // Here we pull the User from the helper given we know there is a 
        // factory that adds it. If there was not one it would be generated new.
        UserID: f.NamedUUID("owner_id"),
        Type: "Something",
    }

    return facto.Product(event)
})
```

### Faking data

Sometimes you need to generate data that is not real but at least looks similar to what the real data will be. To solve that need Facto provides the `Fake` method within the facto.Helper, it can be used within the factory it to generate fake data. for example:

```go
//  in factories/event.go
package factories
import (
    "github.com/paganotoni/facto"
)

func EventFactory(f facto.Helper) facto.Product {
    event := Event{
        Name: f.Fake("Name"),
        Type: "Sports",
        ContactEmail: f.Fake("Email"),
        Company: f.Fake("Company"),
        Address: f.Fake("Address"),
    }

    return facto.Product(event)
})
```

The full list of available fake data generators can be found in [here](link-to-repo).

### The CLI

The Facto CLI is a simple command line tool that allows you to generate fixtures files. It contains the following commands to facilitate the use of Facto:

 * `facto generate`: Generates fixtures for your project.
 * `facto list`: Shows the list of factories in the current codebase.
 * `facto version`: The version of the Facto CLI.


-------------
### Good to go

- Registry API ✅: Good for the first pass.
- Sequences ✅: Index field.
### Other topics

- Create vs Build
- Review terminology from Factory bot.
- Explain the "Magic" constraints and some principles.
