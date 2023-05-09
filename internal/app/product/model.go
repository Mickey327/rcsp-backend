package product

import "time"

type Product struct {
	ID          uint64    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       uint64    `db:"price"`
	Stock       uint64    `db:"stock"`
	Image       string    `db:"image"`
	CategoryID  uint64    `db:"category_id"`
	CompanyID   uint64    `db:"company_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
