package orderItem

import "time"

type OrderItem struct {
	OrderID   uint64    `db:"order_id"`
	ProductID uint64    `db:"product_id"`
	Quantity  uint64    `db:"quantity"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (o *OrderItem) ToDTO() *DTO {
	return &DTO{
		OrderID:   o.OrderID,
		ProductID: o.ProductID,
		Quantity:  o.Quantity,
	}
}

func ToDTOs(orderItems []*OrderItem) []*DTO {
	var orderItemsDTOs []*DTO

	for _, orderItem := range orderItems {
		orderItemsDTOs = append(orderItemsDTOs, orderItem.ToDTO())
	}

	return orderItemsDTOs
}
