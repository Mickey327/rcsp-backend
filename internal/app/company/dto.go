package company

type DTO struct {
	ID   uint64 `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (d *DTO) toCompany() *Company {
	return &Company{
		ID:   d.ID,
		Name: d.Name,
	}
}
