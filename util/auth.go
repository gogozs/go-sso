package util

import (
	"go-ops/models"
	"go-ops/models/users"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(user users.User) error {
	if user.Password != "" {
		var err error
		user.Password, err = GeneratePassword(user.Password)
		if err != nil {
			return err
		}
	}
	u := models.DB.Create(&user) // gorm会自动将id插入struct中
	if u.Error != nil {
		return u.Error
	}
	return nil
}

func GeneratePassword(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ComparePassword(userPwd, hashPwd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(userPwd))
	return err
}

func CheckUser(username, password string) bool {
	user, err := users.GetUser(username)
	if err != nil {
		return false
	} else {
		if err = ComparePassword(password, user.Password); err != nil {
			return false
		}
	}
	return true
}
