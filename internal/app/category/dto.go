package category

type DTO struct {
	ID   uint64 `json:"id,omitempty"`
	Name string `json:"name"`
}

func (d *DTO) toCategory() *Category {
	return &Category{
		ID:   d.ID,
		Name: d.Name,
	}
}
