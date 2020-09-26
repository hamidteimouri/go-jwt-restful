package models

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/hamidteimouri/go-jwt-restful/utils"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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

	/* check for duplication if email address */
	temp := &User{}
	err := *gorm.DB.Table("users").Where("email = ?", user.Email).First(temp).Error

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
	err := *gorm.DB.Create(user).Error
	if err != nil {
		return utils.Message(false, "Failed to create account, connection error"), err
	}
}
