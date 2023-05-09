package order_item

import "github.com/Mickey327/rcsp-backend/internal/app/product"

type DTO struct {
	ID       uint64      `json:"id"`
	Product  product.DTO `json:"product"`
	Quantity uint64      `json:"quantity"`
}
