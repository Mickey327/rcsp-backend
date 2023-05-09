package product

import (
	"github.com/Mickey327/rcsp-backend/internal/app/category"
	"github.com/Mickey327/rcsp-backend/internal/app/company"
)

type DTO struct {
	ID          uint64       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Price       uint64       `json:"price"`
	Stock       uint64       `json:"stock"`
	Image       string       `json:"image"`
	Category    category.DTO `json:"category"`
	Company     company.DTO  `json:"company"`
}
