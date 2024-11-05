package types

type SignUpInput struct {
	Login                string `json:"login"`
	Email                string `json:"email"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}
