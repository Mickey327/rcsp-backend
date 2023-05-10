package company

import "errors"

var (
	CompanyNotFoundErr      = errors.New("company not found")
	CompanyAlreadyExistsErr = errors.New("company with such name already exists")
)
