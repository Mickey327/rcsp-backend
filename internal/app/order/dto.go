package order

import order_item "github.com/Mickey327/rcsp-backend/internal/app/order-item"

type DTO struct {
	ID         uint64           `json:"id"`
	Total      uint64           `json:"total"`
	Status     string           `json:"status"`
	IsFinished bool             `json:"is_finished"`
	OrderItems []order_item.DTO `json:"order_items"`
}
