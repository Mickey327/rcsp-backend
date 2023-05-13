package product

import (
	"time"

	"github.com/Mickey327/rcsp-backend/internal/app/category"
	"github.com/Mickey327/rcsp-backend/internal/app/company"
)

type Product struct {
	ID          uint64             `db:"id"`
	Name        string             `db:"name"`
	Description string             `db:"description"`
	Price       uint64             `db:"price"`
	Stock       uint64             `db:"stock"`
	Image       string             `db:"image"`
	CreatedAt   time.Time          `db:"created_at"`
	UpdatedAt   time.Time          `db:"updated_at"`
	Category    *category.Category `scan:"notate"`
	Company     *company.Company   `scan:"notate"`
}

func (p *Product) ToDTO() *DTO {
	productDTO := &DTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		Image:       p.Image,
	}
	if p.Company != nil {
		productDTO.Company = p.Company.ToDTO()
	}
	if p.Category != nil {
		productDTO.Category = p.Category.ToDTO()
	}

	return productDTO
}

func ToDTOs(products []*Product) []*DTO {
	var productDTOs []*DTO

	for _, product := range products {
		productDTOs = append(productDTOs, product.ToDTO())
	}

	return productDTOs
}
