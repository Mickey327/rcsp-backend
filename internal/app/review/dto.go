package review

import "time"

type DTO struct {
	UserEmail string    `json:"user_email"`
	Message   string    `json:"message"`
	UpdatedAt time.Time `json:"updated_at"`
}
