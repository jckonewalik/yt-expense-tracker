package types

import "github.com/golang-jwt/jwt/v5"

type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}
