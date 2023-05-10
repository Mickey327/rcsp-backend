package company

import "time"

type Company struct {
	ID        uint64    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c *Company) ToDTO() *DTO {
	return &DTO{
		ID:   c.ID,
		Name: c.Name,
	}
}

func ToDTOs(companies []*Company) []*DTO {
	var companyDTOs []*DTO

	for _, company := range companies {
		companyDTOs = append(companyDTOs, company.ToDTO())
	}

	return companyDTOs
}
