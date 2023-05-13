package orderItem

import "github.com/Mickey327/rcsp-backend/internal/app/product"

type DTO struct {
	OrderID  uint64       `json:"order_id,omitempty"`
	Product  *product.DTO `json:"product,omitempty"`
	Quantity int          `json:"quantity,omitempty"`
}

func (d *DTO) ToOrderItem() *OrderItem {
	orderItem := &OrderItem{
		OrderID:  d.OrderID,
		Quantity: d.Quantity,
	}
	if d.Product != nil {
		orderItem.Product = d.Product.ToProduct()
	}
	return orderItem
}

func ToOrderItems(orderItemDTOs []*DTO) []*OrderItem {
	var orderItems []*OrderItem

	for _, orderItemDTO := range orderItemDTOs {
		orderItems = append(orderItems, orderItemDTO.ToOrderItem())
	}

	return orderItems
}
