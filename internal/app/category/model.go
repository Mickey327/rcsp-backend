package category

import "time"

type Category struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *Category) ToDTO() *DTO {
	return &DTO{
		ID:   c.ID,
		Name: c.Name,
	}
}

func ToDTOs(categories []*Category) []*DTO {
	var categoryDTOs []*DTO

	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, category.ToDTO())
	}

	return categoryDTOs
}
