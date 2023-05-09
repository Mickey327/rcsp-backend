package order

import "time"

type Order struct {
	ID         uint64    `db:"id"`
	Total      uint64    `db:"total"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	IsFinished bool      `db:"is_finished"`
	UserID     uint64    `db:"user_id"`
}
