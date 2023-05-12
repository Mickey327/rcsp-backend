package product

type DTO struct {
	ID          uint64 `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       uint64 `json:"price,omitempty"`
	Stock       uint64 `json:"stock,omitempty"`
	Image       string `json:"image,omitempty"`
	CategoryID  uint64 `json:"category_id,omitempty"`
	CompanyID   uint64 `json:"company_id,omitempty"`
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
