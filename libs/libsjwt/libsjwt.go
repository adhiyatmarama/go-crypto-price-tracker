package libsjwt

import (
	"time"

	config "github.com/adhiyatmarama/go-crypto-price-tracker/config"
	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(email string, expTime time.Time) (string, error) {
	claims := &config.JWTClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-crypto-price-tracker",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	tokenAlgo := jwt.NewWithClaims(config.JWT_SIGN_METHOD, claims)
	token, err := tokenAlgo.SignedString(config.JWT_KEY)

	return token, err
}

func ParseToken(tokenString string) (*config.JWTClaim, error) {
	claims := &config.JWTClaim{}
	// parsing token jwt
	_, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return config.JWT_KEY, nil
	})

	return claims, err
}
