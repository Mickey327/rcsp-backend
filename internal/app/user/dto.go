package user

type DTO struct {
	ID       uint64 `json:"id,omitempty" query:"id"`
	Email    string `json:"email,omitempty" query:"email" validate:"required,email"`
	Password string `json:"password,omitempty" query:"password" validate:"required,min=5"`
	Role     string `json:"role,omitempty" query:"role"`
}

func (d *DTO) ToUser() *User {
	return &User{
		ID:       d.ID,
		Email:    d.Email,
		Password: d.Password,
		Role:     d.Role,
	}
}
