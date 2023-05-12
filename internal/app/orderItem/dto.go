package orderItem

type DTO struct {
	OrderID       uint64 `json:"order_id,omitempty"`
	ProductID     uint64 `json:"product_id,omitempty"`
	Quantity      uint64 `json:"quantity,omitempty"`
	QuantityDelta int    `json:"quantity_delta,omitempty"`
}

func (d *DTO) ToOrderItem() *OrderItem {
	return &OrderItem{
		OrderID:   d.OrderID,
		ProductID: d.ProductID,
		Quantity:  d.Quantity,
	}
}

func ToOrderItems(orderItemDTOs []*DTO) []*OrderItem {
	var orderItems []*OrderItem

	for _, orderItemDTO := range orderItemDTOs {
		orderItems = append(orderItems, orderItemDTO.ToOrderItem())
	}

	return orderItems
}
