package usercontroller

import (
	"database/sql"
	"log"
	"time"

	"github.com/adhiyatmarama/go-crypto-price-tracker/middlewares"
	"github.com/adhiyatmarama/go-crypto-price-tracker/user/usermodel"
	"github.com/adhiyatmarama/go-crypto-price-tracker/user/userservice"

	"github.com/adhiyatmarama/go-crypto-price-tracker/libs/libsbcrypt"
	"github.com/adhiyatmarama/go-crypto-price-tracker/libs/libsjwt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

func GetRoutes() *fiber.App {
	userRoute := fiber.New()

	userRoute.Post("/signup", SignUp)
	userRoute.Post("/signin", SignIn)
	userRoute.Get("/signout", middlewares.JWTMiddleware, SignOut)

	return userRoute
}

func validatSignupBody(user usermodel.User) (bool, string) {
	if user.Email == "" || user.Password == "" || user.ConfirmPassword == "" {
		return false, "Email, password, and confirm password must be not empty"
	}

	if user.Password != user.ConfirmPassword {
		return false, "Confirm password should be same as password"
	}

	return true, ""
}

func SignUp(c *fiber.Ctx) error {
	var user usermodel.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
			"error":   err.Error(),
		})
	}

	// validate body
	isValid, validationMessage := validatSignupBody(user)
	if !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": validationMessage,
		})
	}

	createdUser, err := userservice.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error when create user to DB",
			"error":   err.Error(),
		})
	}

	// create jwt token
	expTime := time.Now().Add(time.Minute * 1)
	token, err := libsjwt.CreateToken(createdUser.Email, expTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// create cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expTime,
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": fiber.Map{
			"email": createdUser.Email,
		},
		"message": "Successfully registered and signed in",
	})
}

func SignIn(c *fiber.Ctx) error {
	var userLogin usermodel.User

	if err := c.BodyParser(&userLogin); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "bad request",
			"error":   err.Error(),
		})
	}

	// Get user
	user, err := userservice.GetUserByEmail(userLogin)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "email or password is invalid",
			})
		default:
			log.Print(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Error when get user",
			})
		}
	}

	// check password
	if !libsbcrypt.CheckPasswordHash(userLogin.Password, user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "email or password is invalid",
		})
	}

	// create jwt token
	expTime := time.Now().Add(time.Minute * 1)
	token, err := libsjwt.CreateToken(user.Email, expTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// create cookie
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  expTime,
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"email": user.Email,
		},
		"message": "Successfully signed in",
	})
}

func SignOut(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Minute),
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully signed out",
	})
}
