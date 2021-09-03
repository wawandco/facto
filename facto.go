package facto

import (
	"reflect"
	"sync"
)

// Factory represents the builder of a Product.
type Factory func(f *Helper) Product

// factoriesRegistry is the container for all
// factories registred on runtime.
var factoriesRegistry = map[string]Factory{}

// mu will allow sync access to the factoriesRegistry
// when more than one goroutines is trying to read or
// write on it.
var mu sync.Mutex

// Register allows a factory to be register
// in the factories registry so it can be access
// later when building.
func Register(name string, f Factory) {
	mu.Lock()
	defer mu.Unlock()

	if f != nil && name != "" {
		factoriesRegistry[name] = f
	}
}

// Build requests the factory identified by factoryName
// to build a product.
func Build(factoryName string) Product {
	mu.Lock()
	defer mu.Unlock()

	if factory, ok := factoriesRegistry[factoryName]; ok {
		return factory(&Helper{Index: 1})
	}

	return nil
}

// BuildN requests the factory identified by factoryName
// to build n elements of a product.
func BuildN(factoryName string, n int) Product {
	mu.Lock()
	defer mu.Unlock()

	factory, ok := factoriesRegistry[factoryName]
	if !ok || n <= 0 {
		return nil
	}

	helper := &Helper{Index: 1}
	product := factory(helper)

	products := reflect.MakeSlice(reflect.SliceOf(reflect.ValueOf(product).Type()), 0, n)
	products = reflect.Append(products, reflect.ValueOf(product))

	for i := 1; i < n; i++ {
		helper.Index = i + 1
		product = factory(helper)
		products = reflect.Append(products, reflect.ValueOf(product))
	}

	return Product(products.Interface())
}
