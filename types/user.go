package types

type SignUpInput struct {
	Login                string `json:"login" validate:"required"`
	Email                string `validate:"required,email" json:"email"`
	FirstName            string `json:"first_name" validate:"required"`
	LastName             string `json:"last_name" validate:"required"`
	Password             string `json:"password" validate:"required"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required"`
}
