package ox

import (
	"context"

	"github.com/spf13/pflag"
)

// Plugin for Ox, its mainly a generator but it could be used in here.
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
	return nil
}

// ParseFlags parses the flags used by the plugin
func (p Plugin) ParseFlags([]string) {

}

// Flags that the plugin uses
func (p Plugin) Flags() *pflag.FlagSet {
	return nil
}
