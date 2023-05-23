package company

import "errors"

var (
	CompanyNotFoundErr      = errors.New("компания не найдена")
	CompanyAlreadyExistsErr = errors.New("компания с таким именем уже сушествует")
)
