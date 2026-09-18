package user

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Contact  string `json:"contact" validate:"omitempty"`
	Role     string `json:"role" validate:"omitempty,oneof=SUPER_ADMIN ADMIN USER"`
}

type UpdateProfileRequest struct {
	Name    string `json:"name" validate:"omitempty"`
	Contact string `json:"contact" validate:"omitempty"`
	Image   string `json:"image" validate:"omitempty"`
}
