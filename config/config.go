package config

import (
	"log"
	"os"
)

type Environment struct {
	JwtSecret           string
	AuthApiClientSecret string
}

var Env = NewEnvironment()

func NewEnvironment() Environment {

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("missing environment variable JWT_SECRET")
	}

	authApiClientSecret := os.Getenv("AUTH_API_CLIENT_SECRET")
	if authApiClientSecret == "" {
		log.Fatal("missing environtment variable AUTH_API_CLIENT_SECRET")
	}

	return Environment{
		JwtSecret:           jwtSecret,
		AuthApiClientSecret: authApiClientSecret,
	}
}
