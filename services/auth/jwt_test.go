package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidateToken(t *testing.T) {

	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	expectedName := "John Doe"
	expectedEmail := "john_doe@test.com"
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"name":  expectedName,
		"email": expectedEmail,
	})

	stringToken, _ := token.SignedString(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(&key.PublicKey)

	secret := base64.StdEncoding.EncodeToString(pubKey)

	p, err := ValidateToken(stringToken, secret)

	if err != nil {
		t.Fatalf("not expecting an error at this point: %v", err)
	}

	if p.Name != expectedName {
		t.Fatalf("expecting %s, found %s", expectedName, p.Name)
	}

	if p.Email != expectedEmail {
		t.Fatalf("expecting %s, found %s", expectedEmail, p.Email)
	}

}
