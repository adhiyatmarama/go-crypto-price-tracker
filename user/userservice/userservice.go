package userservice

import (
	"fmt"
	"log"

	"github.com/adhiyatmarama/go-crypto-price-tracker/database"
	"github.com/adhiyatmarama/go-crypto-price-tracker/libs/libsbcrypt"
	"github.com/adhiyatmarama/go-crypto-price-tracker/user/usermodel"
)

func CreateUser(user usermodel.User) (*usermodel.User, error) {
	passwordHash, _ := libsbcrypt.HashPassword(user.Password)

	// Add user to table
	_, err := database.DB.Exec(fmt.Sprintf("INSERT INTO Users(email, password) VALUES('%s', '%s' )", user.Email, passwordHash))
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	// Get user from table
	created, err := GetUserByEmail(user)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return &usermodel.User{
		Email: created.Email,
	}, nil
}

func GetUserByEmail(user usermodel.User) (*usermodel.User, error) {
	// Get user from table
	stmt, _ := database.DB.Prepare("select email, password from Users where email = ?")
	defer stmt.Close()
	var (
		email    string
		password string
	)
	if err := stmt.QueryRow(user.Email).Scan(&email, &password); err != nil {
		log.Print(err.Error())
		return nil, err
	}

	return &usermodel.User{
		Email:    email,
		Password: password,
	}, nil

}
