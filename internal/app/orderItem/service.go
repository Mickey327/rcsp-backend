package orderItem

import (
	"context"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	ChangeQuantityProductToOrderByID(ctx context.Context, orderItem *OrderItem, quantityDelta int) (bool, error)
}

type OrderItemService struct {
	repository Repository
}

func NewService(repository Repository) *OrderItemService {
	return &OrderItemService{
		repository: repository,
	}
}

func (s *OrderItemService) ChangeOrderItemQuantity(c echo.Context, orderItemDTO *DTO) (bool, error) {
	isChanged, err := s.repository.ChangeQuantityProductToOrderByID(c.Request().Context(), orderItemDTO.ToOrderItem(), orderItemDTO.QuantityDelta)
	if err != nil {
		return false, err
	}

	return isChanged, nil
}
