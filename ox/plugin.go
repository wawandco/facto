package ox

import (
	"context"
	"errors"

	"github.com/wawandco/facto"
)

var (
	// ErrIncompleteArgs is returned when the arguments are not enough to generate
	// the factory.
	ErrIncompleteArgs = errors.New("incomplete arguments")
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

// Generate the factory file within factories, e.g. factories/user.go
func (p Plugin) Generate(ctx context.Context, root string, args []string) error {
	// Pass arguments to the generate function
	if len(args) < 3 {
		return ErrIncompleteArgs
	}

	// pass "generate" and the rest of the parameters as the Generate function
	// expects that format.
	args = append([]string{"generate"}, args[2:]...)
	return facto.Generate(root, args)
}
