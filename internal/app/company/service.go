package company

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Create(ctx context.Context, company *Company) (uint64, error)
	Read(ctx context.Context, id uint64) (*Company, error)
	ReadAll(ctx context.Context) ([]*Company, error)
	Update(ctx context.Context, company *Company) (bool, error)
	Delete(ctx context.Context, id uint64) (bool, error)
}

type CompanyService struct {
	repository Repository
}

func NewService(repository Repository) *CompanyService {
	return &CompanyService{
		repository: repository,
	}
}

func (s *CompanyService) Create(c echo.Context, companyDTO *DTO) (uint64, error) {
	id, err := s.repository.Create(c.Request().Context(), companyDTO.ToCompany())

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *CompanyService) Read(c echo.Context, id uint64) (*DTO, error) {
	company, err := s.repository.Read(c.Request().Context(), id)

	if err != nil {
		return nil, err
	}

	return company.ToDTO(), nil
}

func (s *CompanyService) ReadAll(c echo.Context) ([]*DTO, error) {
	companies, err := s.repository.ReadAll(c.Request().Context())

	if err != nil {
		return nil, err
	}

	if len(companies) == 0 {
		return nil, CompanyNotFoundErr
	}

	return ToDTOs(companies), nil
}

func (s *CompanyService) Update(c echo.Context, companyDTO *DTO) (bool, error) {
	isUpdated, err := s.repository.Update(c.Request().Context(), companyDTO.ToCompany())

	if err != nil {
		return false, err
	}

	return isUpdated, nil
}

func (s *CompanyService) Delete(c echo.Context, id uint64) (bool, error) {
	isDeleted, err := s.repository.Delete(c.Request().Context(), id)

	if err != nil {
		return false, err
	}

	return isDeleted, nil
}
