package cart

import (
	"context"
	"log"

	"github.com/Mickey327/rcsp-backend/internal/app/order"
	"github.com/Mickey327/rcsp-backend/internal/app/orderItem"
	"github.com/labstack/echo/v4"
)

type OrderRepository interface {
	GetUserIDByNotArrangedOrderID(ctx context.Context, orderID uint64) (uint64, error)
	ReadCurrentUserArrangingOrderLazy(ctx context.Context, userID uint64) (*order.Order, error)
	ReadCurrentUserArrangingOrderEager(ctx context.Context, userID uint64) (*order.Order, error)
}

type OrderItemRepository interface {
	Create(ctx context.Context, orderItem *orderItem.OrderItem) (bool, error)
	ReadByOrderAndProductID(ctx context.Context, productID, orderID uint64) (*orderItem.OrderItem, error)
	Update(ctx context.Context, orderItem *orderItem.OrderItem) (bool, error)
	Delete(ctx context.Context, orderItem *orderItem.OrderItem) (bool, error)
}

type CartService struct {
	orderRepository     OrderRepository
	orderItemRepository OrderItemRepository
}

func NewService(orderRepository OrderRepository, orderItemRepository OrderItemRepository) *CartService {
	return &CartService{
		orderRepository:     orderRepository,
		orderItemRepository: orderItemRepository,
	}
}

//func (s *CartService) Create(c echo.Context, dto *orderItem.DTO) (bool, error) {
//	_, err := s.orderRepository.GetUserIDByNotArrangedOrderID(c.Request().Context(), dto.OrderID)
//	if err != nil {
//		return false, WrongCartErr
//	}
//
//	isCreated, err := s.orderItemRepository.Create(c.Request().Context(), dto.ToOrderItem())
//	if err != nil {
//		return false, err
//	}
//
//	return isCreated, nil
//}

func (s *CartService) UpdateCart(c echo.Context, dto *orderItem.DTO) (*order.DTO, error) {
	userID, err := s.orderRepository.GetUserIDByNotArrangedOrderID(c.Request().Context(), dto.OrderID)
	if err != nil {
		return nil, WrongCartErr
	}

	item, err := s.orderItemRepository.ReadByOrderAndProductID(c.Request().Context(), dto.OrderID, dto.Product.ID)

	if item != nil {
		log.Println("ITEM QUANTITY:", item.Quantity)
		if dto.Quantity+item.Quantity > 0 {
			log.Println("SUM QUANTITY:", dto.Quantity+item.Quantity)
			_, err = s.orderItemRepository.Update(c.Request().Context(), dto.ToOrderItem())
		} else {
			err = NotPositiveQuantityErr
		}
	} else {
		if dto.Quantity > 0 {
			_, err = s.orderItemRepository.Create(c.Request().Context(), dto.ToOrderItem())
		} else {
			err = NotPositiveQuantityErr
		}
	}

	if err != nil {
		return nil, err
	}

	o, err := s.orderRepository.ReadCurrentUserArrangingOrderEager(c.Request().Context(), userID)
	if err != nil {
		return nil, err
	}

	return o.ToDTO(), nil
}

//func (s *CartService) ReadByOrderAndProductID(c echo.Context, productID, orderID uint64) (*orderItem.DTO, error) {
//	item, err := s.orderItemRepository.ReadByOrderAndProductID(c.Request().Context(), productID, orderID)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return item.ToDTO(), nil
//}
//
//func (s *CartService) Update(c echo.Context, dto *orderItem.DTO) (bool, error) {
//	_, err := s.orderRepository.GetUserIDByNotArrangedOrderID(c.Request().Context(), dto.OrderID)
//	if err != nil {
//		return false, WrongCartErr
//	}
//
//	item, err := s.orderItemRepository.ReadByOrderAndProductID(c.Request().Context(),
//		dto.Product.ID, dto.OrderID)
//
//	if item.Quantity+dto.Quantity < 1 {
//		return false, NotPositiveQuantityErr
//	}
//
//	isUpdated, err := s.orderItemRepository.Update(c.Request().Context(), dto.ToOrderItem())
//	if err != nil {
//		return false, err
//	}
//
//	return isUpdated, nil
//}

func (s *CartService) RemoveFromCart(c echo.Context, dto *orderItem.DTO) (*order.DTO, error) {
	userID, err := s.orderRepository.GetUserIDByNotArrangedOrderID(c.Request().Context(), dto.OrderID)
	if err != nil {
		return nil, WrongCartErr
	}

	isDeleted, err := s.orderItemRepository.Delete(c.Request().Context(), dto.ToOrderItem())
	if err != nil {
		return nil, err
	}

	if !isDeleted {
		return nil, orderItem.OrderItemNotFound
	}

	o, err := s.orderRepository.ReadCurrentUserArrangingOrderEager(c.Request().Context(), userID)
	if err != nil {
		return nil, err
	}

	return o.ToDTO(), nil
}
