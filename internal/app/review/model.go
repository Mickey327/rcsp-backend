package review

import "time"

type Review struct {
	ID        uint64    `db:"id"`
	UserID    uint64    `db:"user_id"`
	ProductID uint64    `db:"product_id"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
