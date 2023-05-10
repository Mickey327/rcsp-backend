package category

import "errors"

var (
	CategoryNotFoundErr      = errors.New("category not found")
	CategoryAlreadyExistsErr = errors.New("category with such name already exists")
)
