package order

import (
	"context"
	"fmt"

	"github.com/Mickey327/rcsp-backend/internal/app/config"
	"github.com/Mickey327/rcsp-backend/internal/app/mail"
	"github.com/labstack/echo/v4"
)

type Repository interface {
	Create(ctx context.Context, userID uint64) (*Order, error)
	ReadCurrentUserArrangingOrderLazy(ctx context.Context, userID uint64) (*Order, error)
	ReadCurrentUserArrangingOrderEager(ctx context.Context, userID uint64) (*Order, error)
	ReadByIdEager(ctx context.Context, id uint64) (*Order, error)
	Update(ctx context.Context, order *Order) (bool, error)
	GetUserEmailByOrderUserID(ctx context.Context, userID uint64) (string, error)
}

type OrderService struct {
	repository Repository
}

func NewService(repository Repository) *OrderService {
	return &OrderService{
		repository: repository,
	}
}

func (s *OrderService) Create(c echo.Context, userID uint64) (*DTO, error) {
	order, err := s.repository.Create(c.Request().Context(), userID)

	if err != nil {
		return nil, err
	}

	return order.ToDTO(), nil
}

func (s *OrderService) ReadByIdEager(c echo.Context, id uint64) (*DTO, error) {
	order, err := s.repository.ReadByIdEager(c.Request().Context(), id)

	if err != nil {
		return nil, err
	}

	return order.ToDTO(), nil
}

func (s *OrderService) Update(c echo.Context, dto *DTO) (bool, error) {
	isUpdated, err := s.repository.Update(c.Request().Context(), dto.ToOrder())
	if err != nil {
		return false, err
	}
	if dto.Status == "Ожидает оплаты" && dto.IsArranged == true {
		email, err := s.repository.GetUserEmailByOrderUserID(c.Request().Context(), dto.UserID)
		if err != nil {
			return false, nil
		}
		cfg := config.GetConfig()
		message := fmt.Sprintf("Здравствуйте, спасибо что оформили у нас заказ! Для оплаты переведите деньги на карту "+
			"1234 5678 9012 3456 или по номеру телефона +7(123)456-78-90, указав в сообщении с переводом вашу почту на сайте "+
			"Статус заказа и его содержание можете отслеживать по ссылке: %s/checkout?orderID=%d", cfg.OuterClientAddress, dto.ID)
		m := mail.New(cfg.Email, email, "Заказ был успешно взят в обработку", message)
		m.SendMail()
		_, err = s.repository.Create(c.Request().Context(), dto.UserID)
		if err != nil {
			return false, err
		}
	}

	return isUpdated, nil
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
