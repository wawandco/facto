package facto

import (
	"reflect"
)

// Factory represents the builder of a Product.
type Factory func(f Helper) Product

// Register allows a factory to be register
// in the factories registry so it can be access
// later when building.
func Register(name string, f Factory) {
	if f == nil || name == "" {
		return
	}

	defaultRegistry.setFactory(name, f)
}

// Build requests the factory identified by factoryName
// to build a product.
func Build(name string) Product {
	f := defaultRegistry.getFactory(name)
	if f == nil {
		return nil
	}

	return f(Helper{Index: 1})
}

// BuildN requests the factory identified by factoryName
// to build n elements of a product.
func BuildN(factoryName string, n int) Product {

	factory := defaultRegistry.getFactory(factoryName)
	if factory == nil || n <= 0 {
		return nil
	}

	helper := Helper{Index: 1}
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
