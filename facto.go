package facto

import (
	"reflect"
)

// Factory represents the builder of a Product.
type Factory func(f Helper) Product

// Build requests the factory identified by factoryName
// to build a product.
func Build(f Factory) Product {
	return f(Helper{Index: 0})
}

// BuildN requests the factory identified by factoryName
// to build n elements of a product.
func BuildN(f Factory, n int) Product {
	h := Helper{Index: 0}
	product := f(h)

	products := reflect.MakeSlice(reflect.SliceOf(reflect.ValueOf(product).Type()), 0, n)
	products = reflect.Append(products, reflect.ValueOf(product))

	for i := 1; i < n; i++ {
		h.Index++
		product = f(h)
		products = reflect.Append(products, reflect.ValueOf(product))
	}

	return Product(products.Interface())
}
