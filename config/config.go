package config

import (
	"log"
	"os"
)

type Environment struct {
	JwtSecret string
}

var Env = newEnvironment()

func newEnvironment() *Environment {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("missing environment variable JWT_SECRET")
	}
	return &Environment{
		JwtSecret: jwtSecret,
	}
}
