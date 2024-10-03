package main

import (
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	tokenString := "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJ6NkMyUXNXc1BTMnFoWVJpZUtrWWtwQ0huQ0tGWk90bnF2TEUwaHJucTVRIn0.eyJleHAiOjE3Mjc5NTkxNjIsImlhdCI6MTcyNzk1ODg2MiwianRpIjoiMjdmMWNlYWQtZDFjNS00MzczLTk3M2MtYTgyMDg4NWRiYTAxIiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3JlYWxtcy95dC1leHBlbnNlLXRyYWNrZXIiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiOWVmZWUyNWQtMTk2Zi00M2JhLTk2OTAtNWIwN2Y4NTI2YzQ0IiwidHlwIjoiQmVhcmVyIiwiYXpwIjoibXktYXBwIiwic2lkIjoiYzA0MzE0NWItYzFhYS00Yjk5LWI3Y2MtYzkyNGUyZjE2ZTUzIiwiYWNyIjoiMSIsImFsbG93ZWQtb3JpZ2lucyI6WyIvKiJdLCJyZWFsbV9hY2Nlc3MiOnsicm9sZXMiOlsiZGVmYXVsdC1yb2xlcy15dC1leHBlbnNlLXRyYWNrZXIiLCJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJlbWFpbCBwcm9maWxlIiwiZW1haWxfdmVyaWZpZWQiOmZhbHNlLCJuYW1lIjoiSm9hbyBTb3V6YSIsInByZWZlcnJlZF91c2VybmFtZSI6ImpvYW8iLCJnaXZlbl9uYW1lIjoiSm9hbyIsImZhbWlseV9uYW1lIjoiU291emEiLCJlbWFpbCI6ImpvYW9AdGVzdC5jb20ifQ.ScxMbw_uUHeQe-d0yDyIJpLKOCCLJt8vHmjDOw7gL6MhWc9BYg0475Unx9zqVqPPz5vlySPwPbg79f5wieiRx5IH2_7JzEFCTAosEP-zLO3JuKYT-FFawsGmKZYI6R8YdkNG6pG0-3sVp4_cnUAvB22hI3J3ocaj7j5Vqz93vjH5EBe4ADpDzQi6uAlm7C8eaq77FFV8mOFXtoS-Qi7qKGc0J4f_TyF2XoLFkeH7DtfrlEN6fpVgMpZsGm56A_4eGj26DIl5UnyN3bLR8ry4TVHCkma4B34fFRydIuwGt7aifvwxGasN8lO-sp41hYbhD8JYUGjMswp8ny0IPZ-70w"
	profile, err := ValidateToken(tokenString)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(profile.Name)
	fmt.Println(profile.Email)
}

type Profile struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenString string) (*Profile, error) {

	token, err := jwt.ParseWithClaims(tokenString, &Profile{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		keyString := os.Getenv("PUBLIC_KEY")

		key, err := base64.StdEncoding.DecodeString(keyString)
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

	if claims, ok := token.Claims.(*Profile); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("error getting token claims")
	}
}
