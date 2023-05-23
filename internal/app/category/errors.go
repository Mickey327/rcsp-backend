package category

import "errors"

var (
	CategoryNotFoundErr      = errors.New("категория не найдена")
	CategoryAlreadyExistsErr = errors.New("категория с таким именем уже существует")
)
