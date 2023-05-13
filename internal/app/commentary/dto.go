package commentary

import "time"

type DTO struct {
	UserEmail string    `json:"user_email,omitempty"`
	Message   string    `json:"message,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
