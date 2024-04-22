package config

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("secret-code-for-go-crypto-price-tracker")
var JWT_SIGN_METHOD = jwt.SigningMethodHS256

type JWTClaim struct {
	Email string
	jwt.RegisteredClaims
}
