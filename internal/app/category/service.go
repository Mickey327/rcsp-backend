package category

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Create(ctx context.Context, category *Category) (uint64, error)
	Read(ctx context.Context, id uint64) (*Category, error)
	ReadAll(ctx context.Context) ([]*Category, error)
	Update(ctx context.Context, category *Category) (bool, error)
	Delete(ctx context.Context, id uint64) (bool, error)
}

type CategoryService struct {
	repository Repository
}

func NewService(repository Repository) *CategoryService {
	return &CategoryService{repository: repository}
}

func (s *CategoryService) Create(c echo.Context, categoryDTO *DTO) (uint64, error) {
	id, err := s.repository.Create(c.Request().Context(), categoryDTO.toCategory())

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *CategoryService) Read(c echo.Context, id uint64) (*DTO, error) {
	category, err := s.repository.Read(c.Request().Context(), id)

	if err != nil {
		return nil, err
	}

	return category.ToDTO(), nil
}

func (s *CategoryService) ReadAll(c echo.Context) ([]*DTO, error) {
	categories, err := s.repository.ReadAll(c.Request().Context())

	if err != nil {
		return nil, err
	}

	if len(categories) == 0 {
		return nil, CategoryNotFoundErr
	}

	return ToDTOs(categories), nil
}

func (s *CategoryService) Update(c echo.Context, categoryDTO *DTO) (bool, error) {
	isUpdated, err := s.repository.Update(c.Request().Context(), categoryDTO.toCategory())

	if err != nil {
		return false, err
	}

	return isUpdated, nil
}

func (s *CategoryService) Delete(c echo.Context, id uint64) (bool, error) {
	isDeleted, err := s.repository.Delete(c.Request().Context(), id)

	if err != nil {
		return false, err
	}

	return isDeleted, nil
}
