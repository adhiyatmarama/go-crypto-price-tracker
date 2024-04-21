package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("secret-code-for-go-crypto-price-tracker")

type JWTClaim struct {
	Email string
	jwt.RegisteredClaims
}
