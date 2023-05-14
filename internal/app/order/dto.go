package order

import "github.com/Mickey327/rcsp-backend/internal/app/orderItem"

type DTO struct {
	ID         uint64           `json:"id"`
	Total      uint64           `json:"total"`
	Status     string           `json:"status"`
	IsArranged bool             `json:"is_arranged"`
	UserID     uint64           `json:"user_id"`
	Count      uint64           `json:"count"`
	OrderItems []*orderItem.DTO `json:"order_items"`
}

func (d *DTO) ToOrder() *Order {
	return &Order{
		ID:         d.ID,
		Total:      d.Total,
		Status:     d.Status,
		IsArranged: d.IsArranged,
		UserID:     d.UserID,
		Count:      d.Count,
		OrderItems: orderItem.ToOrderItems(d.OrderItems),
	}
}
