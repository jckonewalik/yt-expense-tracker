package auth

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jckonewalik/yt-expense-tracker/services/httputils"
	"github.com/jckonewalik/yt-expense-tracker/types"
)

func validateToken(tokenString, secret string) (*types.Profile, error) {

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

func WithJWT(handleFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) != 2 || strings.TrimSpace(bearerToken[1]) == "" {
			httputils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("missing token"))
			return
		}

		publicKey := os.Getenv("PUBLIC_KEY")
		p, err := validateToken(bearerToken[1], publicKey)

		if err != nil {
			httputils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, types.UserName, p.Name)
		ctx = context.WithValue(ctx, types.UserEmail, p.Email)
		r = r.WithContext(ctx)

		handleFunc(w, r)
	}
}
