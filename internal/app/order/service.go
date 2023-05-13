package order

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Create(ctx context.Context, userID uint64) (uint64, error)
	ReadCurrentUserArrangingOrderLazy(ctx context.Context, userID uint64) (*Order, error)
	ReadCurrentUserArrangingOrderEager(ctx context.Context, userID uint64) (*Order, error)
	Update(ctx context.Context, order *Order) (bool, error)
}

type OrderService struct {
	repository Repository
}

func NewService(repository Repository) *OrderService {
	return &OrderService{
		repository: repository,
	}
}

func (s *OrderService) Create(c echo.Context, userID uint64) (uint64, error) {
	id, err := s.repository.Create(c.Request().Context(), userID)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *OrderService) ReadCurrentUserArrangingOrderLazy(c echo.Context, userID uint64) (*DTO, error) {
	order, err := s.repository.ReadCurrentUserArrangingOrderLazy(c.Request().Context(), userID)

	if err != nil {
		return nil, err
	}

	return order.ToDTO(), nil
}

func (s *OrderService) ReadCurrentUserArrangingOrderEager(c echo.Context, userID uint64) (*DTO, error) {
	order, err := s.repository.ReadCurrentUserArrangingOrderEager(c.Request().Context(), userID)

	if err != nil {
		return nil, err
	}

	return order.ToDTO(), nil
}
