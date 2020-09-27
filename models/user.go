package models

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/hamidteimouri/go-jwt-restful/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"
)

/* struct for token */
type Token struct {
	UserId uint
	jwt.StandardClaims
}

type User struct {
	gorm.Model
	Email    string `gorm:"size:255;not null;unique;" json:"email"`
	Password string `gorm:"size:255;not null" json:"password"`
	Token    string `json:"token";sql:"-"`
}

func (user *User) Validate() error {
	/* validating email address */
	if err := checkmail.ValidateFormat(user.Email); err != nil {
		return errors.New("invalid email address")
	}

	/* check for duplication of email address */
	temp := &User{}
	err := getDB().Table("users").Where("email = ?", user.Email).First(temp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return errors.New("connection error")
	}

	if temp.Email != "" {
		return errors.New("email already used")
	}

	if len(user.Password) < 6 {
		return errors.New("password is too short")
	}

	return nil
}

func (user *User) Create() (map[string]interface{}, error) {
	err := user.Validate()
	if err != nil {
		return nil, err
	}

	/* hashing password */
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	/* create the user record */
	err = getDB().Create(user).Error
	if err != nil {
		return utils.Message(false, "Failed to create account, connection error"), err
	}

	/* creating jwt token and set to user struct */
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	/* deleting user's password from struct */
	user.Password = ""

	/* sending back user's account detail as response, without sending its password */
	response := utils.Message(true, "account created")
	response["account"] = user
	return response, nil
}

func SignIn(email, password string) (map[string]interface{}, error) {
	user := &User{}
	err := getDB().Table("users").Where("email = ?", strings.ToLower(email)).First(user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(false, "email address not found"), err
		}
		return utils.Message(false, "connection error"), err
	}

	/* password checking */
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return utils.Message(false, "invalid login credentials. try again"), errors.New("invalid login credentials")
	}

	user.Password = ""
	tk := &Token{UserId: user.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))

	user.Token = tokenString

	/* sending back user's account detail as response, without sending its password */
	response := utils.Message(true, "logged in")
	response["user"] = user
	return response, nil
}

func DeleteUser(id uint) (map[string]interface{}, error) {
	/* check for existing user */
	user := &User{}
	err := getDB().Table("users").Where("id = ?", id).First(user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(false, "User not found"), err
		}
		return utils.Message(false, "Connection error, retry"), err
	}

	/* deleting user and send response back */
	getDB().Table("users").Where("id = ?", id).Delete(user)
	return utils.Message(true, "User Deleted"), nil
}

func UpdatePassword(id uint, newPassword string) (map[string]interface{}, error) {
	/* hashing new password and set to database */
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	getDB().Table("users").Where("id = ?", id).Update("password", hashedPassword)

	/* check whether user exist */
	user := &User{}
	err := getDB().Table("users").Where("id = ?", id).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.Message(false, "user not found"), err
		}

		return utils.Message(false, "Connection error, retry"), err
	}

	/* check whether updated hashed password same aa new input hashedPassword */
	if string(hashedPassword) != user.Password {
		return utils.Message(false, "failed to update password"), errors.New("failed to update password")
	}

	/* sending response back with message */
	return utils.Message(true, "Password Updated"), nil
}
