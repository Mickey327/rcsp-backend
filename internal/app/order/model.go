package order

import (
	"time"

	"github.com/Mickey327/rcsp-backend/internal/app/orderItem"
)

type Order struct {
	ID         uint64    `db:"id"`
	Total      uint64    `db:"total"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	IsArranged bool      `db:"is_arranged"`
	UserID     uint64    `db:"user_id"`
	OrderItems []*orderItem.OrderItem
}

func (o *Order) ToDTO() *DTO {
	return &DTO{
		ID:         o.ID,
		Total:      o.Total,
		Status:     o.Status,
		IsArranged: o.IsArranged,
		UserID:     o.UserID,
		OrderItems: orderItem.ToDTOs(o.OrderItems),
	}
}

func ToDTOs(orders []*Order) []*DTO {
	var orderDTOs []*DTO

	for _, order := range orders {
		orderDTOs = append(orderDTOs, order.ToDTO())
	}

	return orderDTOs
}
