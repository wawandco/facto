# Facto

Facto is a fixtures library with a definition syntax. It aims to allow Go developers to define fixtures in a declarative way and using the Go language as the means to describe the fixtures instead of TOML/JSON/YAML. It is inspired by [factory_bot](https://github.com/thoughtbot/factory_bot).

### Your first Factory

Facto exposes a small API to register and use factories. The following example shows how to create a factory named `User`:

```go
//  in factories/user.go
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

One these are added we can use our factories in our tests, e.g:

```go
//  in users/user_test.go
package user_test

import (
    "your/package/factories"
)

func TestUser(t *testing.T) {
    user := facto.Build(factories.UserFactory).(models.User)
    err := db.Create(&user)
    // use user for test purposes ...

    // Another alternative is to use create, which will attempt at creating the object directly in your database.
    // In which case the user variable has already been stored in the database.
    p, err := facto.Create("User")
    user := p.(models.User)
    // use user for test purposes ...

}
```

### Building/Creating N
Sometimes you need to build more than one instance of an object, in that case. Facto provides 2 functions to do that: `BuildN` and `CreateN`.

```go
//  in test/user_test.go
package user_test

func TestUser(t *testing.T) {
    users := facto.BuildN(factories.User, 10).([]models.User)
    // use users for test purposes ...

    // you can also create N Users
    p, err := facto.CreateN(factories.User, 10)
    // make sure you check the error
    users := p.([]models.User)
}

```

When running `BuildN` you can use the Index variable to get the index of the object being built. This is useful when you need to build a list of objects and you want to know which one is being built.

```go
//  in factories/user.go
package factories

import (
    "fmt"
    "github.com/wawandco/facto"
)

func UserFactory(h facto.Helper) facto.Product {
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
        // Here we pull the User from the helper given we know there is a 
        // factory for it. can use it like this:
        User: facto.Build(UserFactory).(User),
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
    "github.com/wawandco/facto"
)

func UserFactory(h facto.Helper) facto.Product {
    user := User{
        // owner_id will be assigned the generated UUID
        // and any
        ID: h.NamedUUID("owner_id"),
        Name: h.Faker.FirstName(),
    }

    return facto.Product(user)
}

//  in factories/event.go
package factories
import (
    "github.com/wawandco/facto"
)

func EventFactory(h facto.Helper) facto.Product {
    event := Event{
        // Here we pull the User from the helper given we know there is a 
        // factory that adds it. If there was not one it would be generated new.
        UserID: h.NamedUUID("owner_id"),
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
    "github.com/wawandco/facto"
)

func EventFactory(h facto.Helper) facto.Product {
    event := Event{
        Name: h.Faker.FirstName(),
        Type: "Sports",
        ContactEmail: h.Email(),
        Company: h.Faker.Company(),
        Address: h.Faker.Address(),
    }

    return facto.Product(event)
})
```

The fake data generation comes by courtesy of [`Gofakeit`](https://github.com/brianvoe/gofakeit) The full list of available fake data generators can be found in [here](https://github.com/brianvoe/gofakeit/tree/v6.7.1#functions).

### One of

Another thing you could do with Facto is randomize the selection from a list of passed elements. For example:
```go
//  in factories/event.go
package factories
import (
    "github.com/wawandco/facto"
)

func EventFactory(h facto.Helper) facto.Product {
    event := Event{
        Name: h.Faker.FirstName(),
        // You can pass here a list of elements to randomly select
        // from and the facto helper will pick one of these.
        Type: f.OneOf(TypeSports, TypeMusic, TypeConcert).(EventType),
        ...
    }

    return facto.Product(event)
})
```

### The CLI

The Facto CLI is a simple command line tool that allows you to generate fixtures files. To install it you can use the following command:

```
go get github.com/wawandco/facto/cmd/facto@latest
```

The CLI allows you to generate factory files based on a given name. For example:

```sh
facto generate user
# Generates factories/user.go
```
### Facto & Ox
Facto provides a plugin for the Ox CLI plugin system. To use it add the following to your `cmd/ox/main.go` file:

```go
import (
    ...
    fox "github.com/wawandco/facto/cmd/ox" // Add the facto ox package, we called it fox here 😉
)

func main() {
    // You may have other plugins here.
    //...
    cli.Use(fox.Plugin{}) // Add the facto Ox plugin.

    err := cli.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
```

Once added you can use facto with Ox within the `generate` command, you can see that invoking the `generate` command will list `factory` as a generator.

```sh
$ ox generate
[info] Using cmd/ox/main.go 

Available Generators:

  Name          Plugin
  ----          ------
  ox            ox/generate-cli-main
  template      buffalo/generate-template
  model         buffalo/generate-model
  action        buffalo/generate-action
  resource      buffalo/generate-resource
  task          grift/generate-task
  migration     liquo/generate-migration
  factory       facto ## (this is the one you are looking for)
```

Which you can then use to generate your factories with:

```
$ ox generate factory [name]
```

-------------
### Pending
- Create API: Create a new object in the database.
- Review terminology from Factory bot.
- Explain the "Magic" constraints and some principles.