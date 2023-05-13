package orderItem

import (
	"time"

	"github.com/Mickey327/rcsp-backend/internal/app/product"
)

type OrderItem struct {
	OrderID   uint64           `db:"order_id"`
	Quantity  int              `db:"quantity"`
	CreatedAt time.Time        `db:"created_at"`
	UpdatedAt time.Time        `db:"updated_at"`
	Product   *product.Product `scan:"notate"`
}

func (o *OrderItem) ToDTO() *DTO {
	orderItemDTO := &DTO{
		OrderID:  o.OrderID,
		Quantity: o.Quantity,
	}
	if o.Product != nil {
		orderItemDTO.Product = o.Product.ToDTO()
	}
	return orderItemDTO
}

func ToDTOs(orderItems []*OrderItem) []*DTO {
	var orderItemsDTOs []*DTO

	for _, orderItem := range orderItems {
		orderItemsDTOs = append(orderItemsDTOs, orderItem.ToDTO())
	}

	return orderItemsDTOs
}
