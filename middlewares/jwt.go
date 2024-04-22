package middlewares

import (
	"github.com/adhiyatmarama/go-crypto-price-tracker/libs/libsjwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(c *fiber.Ctx) error {
	tokenString := c.Cookies("token")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized missing token",
		})
	}

	parsed, err := libsjwt.ParseToken(tokenString)

	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		switch v.Errors {
		case jwt.ValidationErrorSignatureInvalid:
			// token invalid
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized, Token Invalid",
			})
		case jwt.ValidationErrorExpired:
			// token expired
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized, Token expired",
			})
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorizd",
			})
		}
	}

	c.Locals("user_email", parsed.Email)

	return c.Next()
}
