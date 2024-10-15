package auth

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jckonewalik/yt-expense-tracker/types"
)

func ValidateToken(tokenString, secret string) (*types.Profile, error) {

	token, err := jwt.ParseWithClaims(tokenString, &types.Profile{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		key, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			return nil, err
		}

		pub, err := x509.ParsePKIXPublicKey(key)
		if err != nil {
			return nil, fmt.Errorf("error parsing public key. %v", err)
		}

		return pub, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*types.Profile); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("error getting token claims")
	}
}
