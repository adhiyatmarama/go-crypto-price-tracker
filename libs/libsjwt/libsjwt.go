package libsjwt

import (
	"time"

	config "github.com/adhiyatmarama/go-crypto-price-tracker/config"
	"github.com/golang-jwt/jwt/v4"
)

var SIGN_METHOD = jwt.SigningMethodHS256

func CreateToken(email string, expTime time.Time) (string, error) {
	claims := &config.JWTClaim{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-crypto-price-tracker",
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	}

	tokenAlgo := jwt.NewWithClaims(SIGN_METHOD, claims)
	token, err := tokenAlgo.SignedString(config.JWT_KEY)

	return token, err
}
