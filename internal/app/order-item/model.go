package order_item

import "time"

type OrderItem struct {
	ID        uint64    `db:"id"`
	Quantity  uint64    `db:"quantity"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	OrderID   uint64    `db:"order_id"`
	ProductID uint64    `db:"product_id"`
}
