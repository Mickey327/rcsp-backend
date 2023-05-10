package user

type DTO struct {
	ID       uint64 `json:"id,omitempty" query:"id"`
	Email    string `json:"email,omitempty" query:"email"`
	Password string `json:"password,omitempty" query:"password"`
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
