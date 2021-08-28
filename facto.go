package facto

// facto.Register("User", func(f facto.Helper) interface{} {
// 	return User{
// 		FirstName: "User",
// 		LastName: "Lastname",
// 	}
// })

// user := facto.Build("User")

// facto.Register("Event", func(f facto.Helper) interface{} {
// 	return Event{
// 		User: f.Build("User").(User),
// 		Type: "Something",
// 	}
// })

// user := facto.Build("Event")

// facto.Register("RandomUser", func(f facto.Helper) interface{} {
// 	return User{
// 		FirstName: f.Fake.FirstName(),
// 		LastName: f.Fake.LastName(),
// 	}
// })

// facto.Register("Group", func(f facto.Helper) interface{} {
// 	return Group{
// 		Name: f.Fake.Company(),
// 		Users: f.BuildN("RandomUser", 10),
// 	}
// })

// group := facto.Build("Group")
// user := facto.BuildN("RandomUser")

// facto.Register("User", UserFactory)
// func UserFactory(f facto.Helper) interface{} {
// 	return User{
// 		FirstName: "User",
// 		LastName: "Lastname",
// 	}
// }

// facto.Register("Owner", func(f facto.Helper) interface{} {
// 	return User{
// 		ID: f.UUIDWithName("owner_id"),
// 		FirstName: "User",
// 		LastName: "Lastname",
// 	}
// })

// facto.Register("Cat", func(f facto.Helper) interface{} {
// 	return Cat{
// 		Name: f.Fake.FirstName(),
// 		UserID: f.UUIDWithName("owner_id"),
// 	}
// })
