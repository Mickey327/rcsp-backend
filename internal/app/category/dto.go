package category

type DTO struct {
	ID   uint64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (d *DTO) ToCategory() *Category {
	return &Category{
		ID:   d.ID,
		Name: d.Name,
	}
}
