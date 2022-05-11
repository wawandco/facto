package facto

import (
	"fmt"
	"reflect"
)

// Factory represents the builder of a Product.
type Factory func(f Helper) Product

// Build requests the factory identified by factoryName
// to build a product.
func Build(f Factory) Product {
	return f(NewHelper())
}

// BuildN requests the factory identified by factoryName
// to build n elements of a product.
func BuildN(f Factory, n int) Product {
	h := NewHelper()
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

// Create a record using the configured creator. It calls the Build function
// and then passes the result product to the creator.
func Create(f Factory) (Product, error) {
	if creator == nil {
		return nil, ErrNoCreatorDefined
	}

	p := Build(f)
	if err := creator.Create(p); err != nil {
		return p, fmt.Errorf("could not create products: %w", err)
	}

	return p, nil
}

// Create n records of a given factory. It calls BuildN and then
// passes the result to the creator.
func CreateN(f Factory, n int) (Product, error) {
	if creator == nil {
		return nil, ErrNoCreatorDefined
	}

	p := BuildN(f, n)
	if err := creator.Create(p); err != nil {
		return p, fmt.Errorf("could not create products: %w", err)
	}

	return p, nil
}
