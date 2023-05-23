package product

import "errors"

var (
	ProductNotFoundErr      = errors.New("товар не найден")
	ProductAlreadyExistsErr = errors.New("товар с таким именем уже существует")
)
