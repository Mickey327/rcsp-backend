package order

import "github.com/Mickey327/rcsp-backend/internal/app/orderItem"

type DTO struct {
	ID         uint64           `json:"id,omitempty"`
	Total      uint64           `json:"total,omitempty"`
	Status     string           `json:"status,omitempty"`
	IsArranged bool             `json:"is_arranged"`
	UserID     uint64           `json:"user_id,omitempty"`
	Count      uint64           `json:"count,omitempty"`
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
