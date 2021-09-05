package ox

import (
	"context"

	"github.com/wawandco/facto"
)

// Plugin for Ox, its mainly an adapter to Ox generators
// that can be invoked from the Ox command line, it reuses
// the facto.Generate function which the facto CLI uses as well.
// this plugin should not have any facto logic as that pertains
// to Facto, and not to the plugin.
type Plugin struct{}

// Name of the plugin
func (p Plugin) Name() string {
	return "facto"
}

// InvocationName of the plugin, this is how it will be invoked
// as part of the generators.
func (p Plugin) InvocationName() string {
	return "factory"
}

// Generate:
// - the factory file within app/factories.
// - factories/factories.go with the Load() method if it does not exist.
// And
// - Add facto.Register("Name", NameFactory) to the load method in factories/factories.go
func (p Plugin) Generate(ctx context.Context, root string, args []string) error {
	// Pass arguments to the generate function
	return facto.Generate(root, args)
}
