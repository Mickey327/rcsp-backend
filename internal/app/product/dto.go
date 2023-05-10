package product

type DTO struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Stock       uint64 `json:"stock"`
	Image       string `json:"image"`
	CategoryID  uint64 `json:"category_id"`
	CompanyID   uint64 `json:"company_id"`
}

func (d *DTO) ToProduct() *Product {
	return &Product{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Price:       d.Price,
		Stock:       d.Stock,
		Image:       d.Image,
		CategoryID:  d.CategoryID,
		CompanyID:   d.CompanyID,
	}
}
