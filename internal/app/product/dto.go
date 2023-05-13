package product

import (
	"github.com/Mickey327/rcsp-backend/internal/app/category"
	"github.com/Mickey327/rcsp-backend/internal/app/company"
)

type DTO struct {
	ID          uint64        `json:"id,omitempty"`
	Name        string        `json:"name,omitempty"`
	Description string        `json:"description,omitempty"`
	Price       uint64        `json:"price,omitempty"`
	Stock       uint64        `json:"stock,omitempty"`
	Image       string        `json:"image,omitempty"`
	Category    *category.DTO `json:"category,omitempty"`
	Company     *company.DTO  `json:"company,omitempty"`
}

func (d *DTO) ToProduct() *Product {
	product := &Product{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Stock:       d.Stock,
		Image:       d.Image,
	}
	if d.Company != nil {
		product.Company = d.Company.ToCompany()
	}
	if d.Category != nil {
		product.Category = d.Category.ToCategory()
	}

	return product
}
