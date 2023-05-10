package product

import "errors"

var (
	ProductNotFoundErr      = errors.New("product not found")
	ProductAlreadyExistsErr = errors.New("product with such name already exists")
)
